package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/api"
	"github.com/firebase/genkit/go/genkit"
	_ "github.com/lib/pq"
	pgv "github.com/pgvector/pgvector-go"
)

const (
	provider     = "pgvector"
	embedderName = "text-embedding-3-small"
)

var (
	connString = flag.String("dbconn", "postgres://postgres:password@localhost:5432/mydb?sslmode=disable", "database connection string")
	index      = flag.Bool("index", false, "index the existing data")
)

func main() {
	baseURL := os.Getenv("AZ_OPENAI_BASE_URL")
	apiKey := os.Getenv("AZ_OPENAI_API_KEY")
	if baseURL == "" || apiKey == "" {
		log.Fatal("export AZ_OPENAI_BASE_URL and AZ_OPENAI_API_KEY to run this sample")
	}

	flag.Parse()
	ctx := context.Background()
	azOpenAI := &AzureOpenAI{
		BaseURL: baseURL,
		APIKey:  apiKey,
		// The sample assumes the use of the Azure OpenAI v1 API version.
		// If you want to use 2024-10-21 instead, make sure to deploy the
		// text-embedding-3-small model with exactly that deployment name and
		// uncomment the following line.
		// Deployment: embedderName,
	}
	g := genkit.Init(ctx, genkit.WithPlugins(azOpenAI))
	if err := run(g, azOpenAI); err != nil {
		log.Fatal(err)
	}
}

func run(g *genkit.Genkit, aoai *AzureOpenAI) error {
	if *connString == "" {
		return errors.New("need -dbconn")
	}
	ctx := context.Background()
	embedder := aoai.Embedder(g, embedderName)
	if embedder == nil {
		return fmt.Errorf("embedder %s is not known to the googlegenai plugin", embedderName)
	}

	db, err := sql.Open("postgres", *connString)
	if err != nil {
		return err
	}
	defer db.Close()

	if *index {
		if err := indexExistingRows(ctx, g, db, embedder); err != nil {
			return err
		}
	}

	retOpts := &ai.RetrieverOptions{
		ConfigSchema: nil,
		Label:        "pgVector",
		Supports: &ai.RetrieverSupports{
			Media: false,
		},
	}
	retriever := defineRetriever(g, db, embedder, retOpts)

	type input struct {
		Question string
		Show     string
	}

	genkit.DefineFlow(g, "askQuestion", func(ctx context.Context, in input) (string, error) {
		res, err := genkit.Retrieve(ctx, g,
			ai.WithRetriever(retriever),
			ai.WithConfig(in.Show),
			ai.WithTextDocs(in.Question))
		if err != nil {
			return "", err
		}
		for _, doc := range res.Documents {
			fmt.Printf("%+v %q\n", doc.Metadata, doc.Content[0].Text)
		}
		// Use documents in RAG prompts.
		return "", nil
	})

	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	fmt.Println("Press Ctrl-C to stop")
	<-sigCtx.Done()
	return nil
}

func defineRetriever(g *genkit.Genkit, db *sql.DB, embedder ai.Embedder, retOpts *ai.RetrieverOptions) ai.Retriever {
	f := func(ctx context.Context, req *ai.RetrieverRequest) (*ai.RetrieverResponse, error) {
		eres, err := genkit.Embed(ctx, g,
			ai.WithEmbedder(embedder),
			ai.WithDocs(req.Query))
		if err != nil {
			return nil, err
		}
		rows, err := db.QueryContext(ctx, `
			SELECT episode_id, season_number, chunk as content
			FROM embeddings
			WHERE show_id = $1
		  	ORDER BY embedding <#> $2
		  	LIMIT 2`,
			req.Options, pgv.NewVector(eres.Embeddings[0].Embedding))
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		res := &ai.RetrieverResponse{}
		for rows.Next() {
			var eid, sn int
			var content string
			if err := rows.Scan(&eid, &sn, &content); err != nil {
				return nil, err
			}
			meta := map[string]any{
				"episode_id":    eid,
				"season_number": sn,
			}
			doc := &ai.Document{
				Content:  []*ai.Part{ai.NewTextPart(content)},
				Metadata: meta,
			}
			res.Documents = append(res.Documents, doc)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		return res, nil
	}
	return genkit.DefineRetriever(g, api.NewName(provider, "shows"), retOpts, f)
}

// Helper function to get started with indexing
func Index(ctx context.Context, g *genkit.Genkit, db *sql.DB, embedder ai.Embedder, docs []*ai.Document) error {
	// The indexer assumes that each Document has a single part, to be embedded, and metadata fields
	// for the table primary key: show_id, season_number, episode_id.
	const query = `
			UPDATE embeddings
			SET embedding = $4
			WHERE show_id = $1 AND season_number = $2 AND episode_id = $3
		`
	res, err := genkit.Embed(ctx, g,
		ai.WithEmbedder(embedder),
		ai.WithDocs(docs...))
	if err != nil {
		return err
	}
	// You may want to use your database's batch functionality to insert the embeddings
	// more efficiently.
	for i, emb := range res.Embeddings {
		doc := docs[i]
		args := make([]any, 4)
		for j, k := range []string{"show_id", "season_number", "episode_id"} {
			if a, ok := doc.Metadata[k]; ok {
				args[j] = a
			} else {
				return fmt.Errorf("doc[%d]: missing metadata key %q", i, k)
			}
		}
		args[3] = pgv.NewVector(emb.Embedding)
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}
	return nil

}

func indexExistingRows(ctx context.Context, g *genkit.Genkit, db *sql.DB, embedder ai.Embedder) error {
	rows, err := db.QueryContext(ctx, `SELECT show_id, season_number, episode_id, chunk FROM embeddings`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var docs []*ai.Document
	for rows.Next() {
		var sid, chunk string
		var sn, eid int
		if err := rows.Scan(&sid, &sn, &eid, &chunk); err != nil {
			return err
		}
		docs = append(docs, &ai.Document{
			Content: []*ai.Part{ai.NewTextPart(chunk)},
			Metadata: map[string]any{
				"show_id":       sid,
				"season_number": sn,
				"episode_id":    eid,
			},
		})
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return Index(ctx, g, db, embedder, docs)
}
