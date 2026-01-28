package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"github.com/firebase/genkit/go/plugins/localvec"
	"github.com/ledongthuc/pdf"
	"github.com/tmc/langchaingo/textsplitter"
)

func readPDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		return "", err
	}

	reader, err := r.GetPlainText()
	if err != nil {
		return "", err
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func main() {
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(2000),
		textsplitter.WithChunkOverlap(20))

	ctx := context.Background()

	g := genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{}))

	docStore, pdfRetriever, err := localvec.DefineRetriever(
		g, "travelQA", localvec.Config{Embedder: googlegenai.GoogleAIEmbedder(g, "gemini-embedding-001")}, nil)
	if err != nil {
		log.Fatalf("unable to create docstore/retriever: %v", err)
	}

	genkit.DefineFlow(g, "indexBrochure",
		func(ctx context.Context, path string) (any, error) {
			pdfText, err := genkit.Run(ctx, "extract", func() (string, error) {
				return readPDF(path)
			})
			if err != nil {
				return nil, err
			}

			docs, err := genkit.Run(ctx, "chunk", func() ([]*ai.Document, error) {
				chunks, err := splitter.SplitText(pdfText)
				if err != nil {
					return nil, err
				}

				var docs []*ai.Document
				for _, chunk := range chunks {
					docs = append(docs, ai.DocumentFromText(chunk, nil))
				}
				return docs, nil
			})
			if err != nil {
				return nil, err
			}
			localvec.Index(ctx, docs, docStore)

			return map[string]any{
				"success":          true,
				"documentsIndexed": len(docs),
			}, nil
		})
	if err := localvec.Init(); err != nil {
		log.Fatalf("unable to index documents: %v", err)
	}

	genkit.DefineFlow(g, "travelQA", func(ctx context.Context, question string) (string, error) {
		// Retrieve text relevant to the user's question.
		resp, err := genkit.Retrieve(ctx, g, ai.WithRetriever(pdfRetriever), ai.WithTextDocs(question))
		if err != nil {
			return "", err
		}
		if len(resp.Documents) == 0 {
			fmt.Println("No documents found by retriever.")
		}

		// Call Generate, including the menu information in your prompt
		return genkit.GenerateText(ctx, g,
			ai.WithModelName("googleai/gemini-3-flash-preview"),
			ai.WithDocs(resp.Documents...),
			ai.WithSystem(`
You are an AI assistant that helps with travel-related inquiries, offering tips, advice, and recommendations 
as a knowledgeable travel agent.
Use only the context provided to answer the question. If you don't know, do not
make up an answer. Do not add or change details of the travel destinations you have been provided.`),
			ai.WithPrompt(question))
	})

	// Prevent main() to terminate
	fmt.Println("RAG application started. Press Ctrl-C to stop.")
	fmt.Println("Use the Genkit CLI to run flows, e.g. genkit flow:run indexBrochure '\"brochures/Dubai Brochure.pdf\"'")
	notifyCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	<-notifyCtx.Done()
}
