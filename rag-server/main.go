package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/weaviate"
)

const systemPrompt = `
I will ask you a question and will provide some additional context information.
Assume this context information is factual and correct, as part of internal
documentation.
If the question relates to the context, answer it using the context.
If the question does not relate to the context, answer it as normal.

For example, let's say the context has nothing in it about tropical flowers;
then if I ask you about tropical flowers, just answer what you know about them
without referring to the context.

For example, if the context does mention minerology and I ask you about that,
provide information from the context along with general knowledge.
`

func main() {
	generativeModelName := cmp.Or(os.Getenv("GENKIT_MODEL"), "googleai/gemini-2.5-flash")
	embeddingModelName := cmp.Or(os.Getenv("EMBEDDING_MODEL"), "text-embedding-004")

	ctx := context.Background()

	googleAI := &googlegenai.GoogleAI{}

	g := genkit.Init(ctx, genkit.WithPlugins(googleAI, &weaviate.Weaviate{
		Addr:   "localhost:9035",
		Scheme: "http",
		APIKey: "", // No auth for local Weaviate
	}))

	embedder, err := googleAI.DefineEmbedder(g, embeddingModelName, &ai.EmbedderOptions{
		Dimensions: 768,
	})
	if err != nil {
		log.Fatalf("unable to set up embedder %q: %v", embeddingModelName, err)
	}

	indexer, retriever, err := weaviate.DefineRetriever(ctx, g, weaviate.ClassConfig{
		Class:    "Document",
		Embedder: embedder,
	}, nil)

	if err != nil {
		log.Fatalf("unable to set up Weaviate retriever: %v", err)
	}

	model := genkit.LookupModel(g, generativeModelName)
	if model == nil {
		log.Fatalf("unable to find model %q: %v", generativeModelName, err)
	}

	server := &ragServer{
		ctx:       ctx,
		g:         g,
		indexer:   indexer,
		retriever: retriever,
		model:     model,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /add/", server.addDocumentsHandler)
	mux.HandleFunc("POST /query/", server.queryHandler)

	port := cmp.Or(os.Getenv("SERVERPORT"), "9020")
	address := net.JoinHostPort("localhost", port)
	log.Println("listening on", address)
	log.Fatal(http.ListenAndServe(address, mux))
}

type ragServer struct {
	ctx       context.Context
	g         *genkit.Genkit
	indexer   *weaviate.Docstore
	retriever ai.Retriever
	model     ai.Model
}

func (rs *ragServer) addDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type document struct {
		Text string
	}
	type addRequest struct {
		Documents []document
	}
	ar := &addRequest{}
	if err := readRequestJSON(req, ar); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert request documents into Weaviate documents used for embedding.
	var wvDocs []*ai.Document
	for _, doc := range ar.Documents {
		wvDocs = append(wvDocs, ai.DocumentFromText(doc.Text, nil))
	}

	// Index the requested documents.
	if err := weaviate.Index(rs.ctx, wvDocs, rs.indexer); err != nil {
		http.Error(w, fmt.Errorf("indexing: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (rs *ragServer) queryHandler(w http.ResponseWriter, req *http.Request) {
	// Parse HTTP request from JSON.
	type queryRequest struct {
		Content string
	}
	qr := &queryRequest{}
	err := readRequestJSON(req, qr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the most similar documents using the retriever.
	resp, err := genkit.Retrieve(rs.ctx, rs.g, ai.WithRetriever(rs.retriever),
		ai.WithTextDocs(qr.Content), ai.WithConfig(&weaviate.RetrieverOptions{Count: 3}))
	if err != nil {
		http.Error(w, fmt.Errorf("retrieval: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("retrieved %d documents", len(resp.Documents))

	// Create a RAG query for the LLM with the most relevant documents as
	// context.
	genResp, err := genkit.Generate(rs.ctx, rs.g, ai.WithModel(rs.model),
		ai.WithSystem(systemPrompt), ai.WithDocs(resp.Documents...), ai.WithPrompt(qr.Content))
	if err != nil {
		log.Printf("calling generative model: %v", err.Error())
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	if len(genResp.Message.Content) != 1 {
		log.Printf("got %v candidates, expected 1", len(genResp.Message.Content))
		http.Error(w, "generative model error", http.StatusInternalServerError)
		return
	}

	renderJSON(w, genResp.Text())
}
