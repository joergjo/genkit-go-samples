package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
)

func main() {
	ctx := context.Background()

	baseURL := os.Getenv("AZ_OPENAI_BASE_URL")
	apiKey := os.Getenv("AZ_OPENAI_API_KEY")
	if baseURL == "" || apiKey == "" {
		log.Fatal("export AZ_OPENAI_BASE_URL and AZ_OPENAI_API_KEY to run this sample")
	}

	fmt.Println("Using Entra ID authentication for Azure OpenAI")
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("could not create credential: %v\n", err)
	}

	azOpenAI := &AzureOpenAI{
		BaseURL:         baseURL,
		TokenCredential: cred,
	}

	g := genkit.Init(ctx, genkit.WithPlugins(azOpenAI))
	model := azOpenAI.Model(g, "gpt-5-mini")

	text, err := generate(ctx, g, model, "Invent a menu for a pirate-themed restaurant")
	if err != nil {
		log.Fatalf("could not generate model response: %v\n", err)
	}
	fmt.Println(text)

	fmt.Println("")
	fmt.Println("---------------------------------------------------")
	fmt.Println("")

	fmt.Println("Using API key for Azure OpenAI")
	// We already know that the API key is not empty.
	azOpenAI = &AzureOpenAI{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
	g = genkit.Init(ctx, genkit.WithPlugins(azOpenAI))

	text, err = generate(ctx, g, model, "Invent a menu for a pirate-themed restaurant")
	if err != nil {
		log.Fatalf("could not generate model response: %v\n", err)
	}
	log.Println(text)

}

func generate(ctx context.Context, g *genkit.Genkit, model ai.Model, prompt string) (string, error) {
	// Simple chat completion
	resp, err := genkit.Generate(ctx, g,
		ai.WithPrompt(prompt),
		ai.WithModel(model))
	if err != nil {
		return "", err
	}
	return resp.Text(), nil
}
