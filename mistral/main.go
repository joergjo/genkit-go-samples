package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/thomas-marquis/genkit-mistral/mistral"
	mistralclient "github.com/thomas-marquis/mistral-client/mistral"
)

func main() {
	apiKey := os.Getenv("MISTRAL_API_KEY")
	if apiKey == "" {
		log.Fatal("MISTRAL_API_KEY environment variable not set")
	}

	ctx := context.Background()

	g := genkit.Init(ctx, genkit.WithPlugins(mistral.NewPlugin(apiKey, mistral.WithClientOptions(mistralclient.WithClientTimeout(time.Duration(60*time.Second))))),
		genkit.WithDefaultModel("mistral/mistral-small-latest"))

	// Simple chat completion
	fmt.Println("=== Chat Completion ===")
	resp, err := genkit.Generate(ctx, g, ai.WithPrompt("Invent a menu for a pirate-themed restaurant"))
	if err != nil {
		log.Fatalf("could not generate model response: %v", err)
	}
	fmt.Println(resp.Text())

	// Simple embedding generation
	fmt.Println("=== Embedding Generation ===")
	res, err := genkit.Embed(ctx, g,
		ai.WithEmbedderName("mistral/mistral-embed"),
		ai.WithDocs(ai.DocumentFromText("GenKit is a Go library for working with generative AI models.", nil)))
	if err != nil {
		log.Fatalf("could not generate embeddings: %v", err)
	}
	fmt.Printf("Embedding size: %d\n", len(res.Embeddings[0].Embedding))
}
