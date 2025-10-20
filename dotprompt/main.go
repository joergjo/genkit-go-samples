package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

func main() {
	var city, likes, dislikes string
	var days, attractions int
	flag.StringVar(&city, "city", "", "City to get the trip advice for")
	flag.StringVar(&likes, "likes", "", "Things the user likes")
	flag.StringVar(&dislikes, "dislikes", "", "Things the user dislikes")
	flag.IntVar(&days, "days", 0, "Numnber of days to stay")
	flag.IntVar(&attractions, "attractions", 0, "Number of attractions to visit")
	flag.Parse()

	if city == "" {
		log.Fatal("city is required")
	}
	promptName := "tripadvisor_single"
	if likes != "" && dislikes != "" && days > 0 && attractions > 0 {
		promptName = "tripadvisor_multi"
	}

	ctx := context.Background()
	g := genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{}))

	prompt := genkit.LookupPrompt(g, promptName)
	if prompt == nil {
		log.Fatalf("prompt %q not found", promptName)
	}

	input := map[string]any{
		"city": city,
	}
	if promptName == "tripadvisor_multi" {
		input["likes"] = likes
		input["dislikes"] = dislikes
		input["n_days"] = days
		input["n_attractions"] = attractions
	}
	resp, err := prompt.Execute(ctx, ai.WithInput(input))
	if err != nil {
		log.Fatalf("failed to execute prompt: %v", err)
	}
	fmt.Println(resp.Text())

	notifyCtx, stop := signal.NotifyContext(ctx)
	defer stop()
	<-notifyCtx.Done()
}
