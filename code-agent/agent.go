package main

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

type Agent struct {
	flow           *core.Flow[string, string, struct{}]
	g              *genkit.Genkit
	getUserMessage func(context.Context) (string, bool)
	history        []*ai.Message
	tools          []ai.ToolRef
}

func NewAgent(ctx context.Context, getUserMessage func(context.Context) (string, bool)) *Agent {
	a := &Agent{getUserMessage: getUserMessage}
	g := genkit.Init(ctx, genkit.WithPlugins(&googlegenai.GoogleAI{}), genkit.WithDefaultModel("googleai/gemini-2.5-flash"))
	readFile := genkit.DefineTool(g, ReadFileDefinition.Name, ReadFileDefinition.Description, ReadFile)
	listFiles := genkit.DefineTool(g, ListFilesDescription.Name, ListFilesDescription.Description, ListFiles)
	editFile := genkit.DefineTool(g, EditFileDescription.Name, EditFileDescription.Description, EditFile)
	a.tools = []ai.ToolRef{readFile, listFiles, editFile}
	a.g = g

	a.flow = genkit.DefineFlow(g, "run_inference", func(ctx context.Context, input string) (string, error) {
		resp, err := genkit.Generate(ctx, a.g, ai.WithPrompt(input),
			ai.WithMessages(a.history...), ai.WithTools(a.tools...))
		if err != nil {
			return "", err
		}
		a.history = append(resp.History(), resp.Message)
		return resp.Text(), nil
	})

	return a
}

func (a *Agent) Run(ctx context.Context) error {
	fmt.Println("Chat with your code. Use CTRL-C to quit.")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInput, ok := a.getUserMessage(ctx)
		if !ok {
			break
		}

		content, err := a.flow.Run(ctx, userInput)
		if err != nil {
			return err
		}
		fmt.Printf("\u001b[93mAgent\u001b[0m: %s\n", content)
	}

	return nil
}
