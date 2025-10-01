package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/firebase/genkit/go/core/api"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

const (
	apiVersion = "2024-10-21"
)

type AzureOpenAI struct {
	*oai.OpenAI
	APIKey      string
	AccessToken string
	Endpoint    string
	Deployment  string
}

func (a *AzureOpenAI) Init(ctx context.Context) []api.Action {
	if a.APIKey == "" && a.AccessToken == "" || a.APIKey != "" && a.AccessToken != "" {
		panic("azopenai plugin initialization failed: either APIKey or AccessToken is required")
	}
	if a.Deployment == "" && a.Endpoint == "" {
		panic("azopenai plugin initialization failed: Endpoint and Deployment are required")
	}

	if a.OpenAI == nil {
		// Overwrite base URL, set "api-version" query parameter, and remove JSON attribute "model"
		a.OpenAI = &oai.OpenAI{
			APIKey: a.APIKey,
			Opts: []option.RequestOption{
				option.WithBaseURL(baseURL(a.Endpoint, a.Deployment)),
				option.WithQuery("api-version", apiVersion),
				option.WithJSONDel("model"),
			},
		}

		switch a.APIKey {
		case "":
			// Satisfy OpenAI's requirement for a non-empty string
			a.OpenAI.APIKey = "notused"
			// Set "Authorization" header with bearer token
			a.OpenAI.Opts = append(a.OpenAI.Opts, option.WithHeader("Authorization", "Bearer "+a.AccessToken))
		default:
			// Satisfy OpenAI's requirement for an API key
			a.OpenAI.APIKey = a.APIKey
			// Use the "api-key" header instead of "Authorization"
			a.OpenAI.Opts = append(a.OpenAI.Opts,
				option.WithHeader("api-key", a.APIKey),
				option.WithHeaderDel("Authorization"))
		}
	}

	// Enable HTTP request/response logging if AZ_OAI_DEBUG_HTTP environment variable is set to "1" or "true"
	debug := os.Getenv("AZ_OAI_DEBUG_HTTP")
	if cmp.Or(debug == "1", strings.EqualFold(debug, "true")) {
		a.OpenAI.Opts = append(a.OpenAI.Opts, option.WithDebugLog(log.Default()))
	}

	return a.OpenAI.Init(ctx)
}

func baseURL(endpoint, deployment string) string {
	u, err := url.JoinPath(endpoint, "openai", "deployments", deployment)
	if err != nil {
		panic(fmt.Sprintf("unexpected error generating base URL: %v", err))
	}
	return u
}
