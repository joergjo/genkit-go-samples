package main

import (
	"context"
	"fmt"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/core/x/session"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/googlegenai"
)

type ChatState struct {
	History []*ai.Message `json:"history"`
}

type sessionKey struct{}

func main() {
	ctx := context.Background()
	g := genkit.Init(ctx,
		genkit.WithPlugins(&googlegenai.GoogleAI{}),
		genkit.WithDefaultModel("googleai/gemini-2.5-flash"))

	// Create a store to persist sessions across requests
	store := session.NewInMemoryStore[ChatState]()

	chat := genkit.DefineFlow(g, "chat", func(ctx context.Context, input string) (string, error) {
		sessionID, ok := ctx.Value(sessionKey{}).(string)
		if !ok {
			// We ignore missing session ID and create a new session
			sessionID = ""
		}

		// Load existing session or create new one
		sess, err := session.Load(ctx, store, sessionID)
		if err != nil {
			sess, err = session.New(ctx,
				session.WithID[ChatState](sessionID),
				session.WithStore(store),
				session.WithInitialState(ChatState{}),
			)
			if err != nil {
				return "", err
			}
		}
		state := sess.State()

		// Attach session to context for use in tools and prompts
		ctx = session.NewContext(ctx, sess)

		// Generate with the session-aware context
		resp, err := genkit.Generate(ctx, g,
			ai.WithPrompt(input),
			ai.WithMessages(state.History...))

		state.History = append(resp.History(), resp.Message)
		if err := sess.UpdateState(ctx, state); err != nil {
			panic(err)
		}
		return resp.Text(), err
	})

	// The session ID could come from an HTTP request header, cookie, etc.
	// In this example, we just use a fixed value to demonstrate session persistence across
	// multiple calls to the "chat" flow.
	ctx = context.WithValue(ctx, sessionKey{}, "local-session")
	prompts := []string{
		"Tell me a pirate joke!",
		"Now tell that joke in the voice of a pirate.",
	}

	for _, prompt := range prompts {
		response, err := chat.Run(ctx, prompt)
		if err != nil {
			panic(err)
		}
		fmt.Println(response)
	}
}
