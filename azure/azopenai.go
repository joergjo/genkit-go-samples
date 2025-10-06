package main

import (
	"cmp"
	"context"
	"log"
	"os"
	"strings"

	"github.com/firebase/genkit/go/core/api"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

type AzureOpenAI struct {
	*oai.OpenAI
	APIKey      string
	AccessToken string
	Endpoint    string
}

func (a *AzureOpenAI) Init(ctx context.Context) []api.Action {
	if a.APIKey == "" && a.AccessToken == "" || a.APIKey != "" && a.AccessToken != "" {
		panic("azopenai plugin initialization failed: either APIKey or AccessToken is required")
	}
	if a.Endpoint == "" {
		panic("azopenai plugin initialization failed: Endpoint is required")
	}

	if a.OpenAI == nil {
		// Overwrite base URL and provide API key
		a.OpenAI = &oai.OpenAI{
			APIKey: a.APIKey,
			Opts: []option.RequestOption{
				option.WithBaseURL(a.Endpoint),
			},
		}

		switch a.APIKey {
		case "":
			// Satisfy OpenAI's requirement for a non-empty string
			a.OpenAI.APIKey = "notused"
			// Set "Authorization" header with bearer token
			a.OpenAI.Opts = append(a.OpenAI.Opts, option.WithHeader("Authorization", "Bearer "+a.AccessToken))
		default:
			a.OpenAI.APIKey = a.APIKey
		}
	}

	// Enable HTTP request/response logging if AZ_OPENAI_DEBUG_HTTP environment variable is set to "1" or "true"
	debug := os.Getenv("AZ_OPENAI_DEBUG_HTTP")
	if cmp.Or(debug == "1", strings.EqualFold(debug, "true")) {
		a.OpenAI.Opts = append(a.OpenAI.Opts, option.WithDebugLog(log.Default()))
	}

	return a.OpenAI.Init(ctx)
}
