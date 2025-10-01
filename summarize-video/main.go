package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
	"google.golang.org/genai"
)

func main() {
	// Get command line arguments
	var url, prompt string
	flag.StringVar(&url, "url", "", "URL of video to summarize")
	flag.StringVar(&prompt, "prompt", "Please summarize the following video:", "Prompt to use for summarization")
	flag.Parse()

	if url == "" {
		log.Fatalf("Please provide a URL of a video to summarize using the -url flag.")
	}

	fmt.Printf("=== Summarizing video %q ===\n", url)

	// Initialize genkit with Google GenAI plugin
	ctx := context.Background()
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"))

	// Generate summary
	resp, err := genkit.Generate(ctx, g,
		ai.WithMessages(ai.NewUserMessage(
			ai.NewMediaPart("video/mp4", url),
			ai.NewTextPart(prompt))),
		ai.WithConfig(&genai.GenerateContentConfig{
			MaxOutputTokens: 65536,
		}))
	if err != nil {
		log.Fatalf("Failed to generate summary: %v", err)
	}
	fmt.Printf("=== Summary ===\n%s", resp.Text())
}
