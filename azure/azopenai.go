package main

import (
	"cmp"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/firebase/genkit/go/core/api"
	oai "github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/azure"
	"github.com/openai/openai-go/option"
)

const apiVersion = "2024-10-21"

type AzureOpenAI struct {
	*oai.OpenAI
	APIKey          string
	TokenCredential azcore.TokenCredential
	BaseURL         string
	Deployment      string
}

func (a *AzureOpenAI) Init(ctx context.Context) []api.Action {
	if a.APIKey == "" && a.TokenCredential == nil || a.APIKey != "" && a.TokenCredential != nil {
		panic("Azure OpenAI plugin initialization failed: either APIKey or TokenCredential is required")
	}
	if a.BaseURL == "" {
		panic("Azure OpenAI plugin initialization failed: Endpoint is required")
	}

	if a.OpenAI == nil {
		switch a.Deployment {
		case "":
			a.init()
		default:
			a.initWithDeployment()
		}
	}

	// Enable HTTP request/response logging if AZ_OPENAI_DEBUG_HTTP environment variable is set to "1" or "true"
	debug := os.Getenv("AZ_OPENAI_DEBUG_HTTP")
	if cmp.Or(debug == "1", strings.EqualFold(debug, "true")) {
		a.OpenAI.Opts = append(a.OpenAI.Opts, option.WithDebugLog(log.Default()))
	}

	return a.OpenAI.Init(ctx)
}

func (a *AzureOpenAI) init() {
	// Overwrite base URL and provide API key
	a.OpenAI = &oai.OpenAI{
		APIKey: a.APIKey,
		Opts: []option.RequestOption{
			option.WithBaseURL(a.BaseURL),
		},
	}

	// If no API key is provided, use TokenCredential (Entra) for authorization
	if a.APIKey == "" {
		// Satisfy the OpenAI plugin's requirement for a non-empty string
		a.OpenAI.APIKey = "notused"
		// Inject bearer token middleware
		a.OpenAI.Opts = append(a.OpenAI.Opts, azure.WithTokenCredential(a.TokenCredential))
	}
}

func (a *AzureOpenAI) initWithDeployment() {
	// Build the effective base URL with deployment path
	// Note: This should never fail unless BaseURL was a non-empty string that is not a valid URL
	u, err := url.JoinPath(a.BaseURL, "openai", "deployments", a.Deployment)
	if err != nil {
		panic(fmt.Sprintf("unexpected error generating base URL: %v", err))
	}

	// Overwrite base URL, set "api-version" query parameter, and remove JSON attribute "model"
	a.OpenAI = &oai.OpenAI{
		APIKey: a.APIKey,
		Opts: []option.RequestOption{
			option.WithBaseURL(u),
			option.WithQuery("api-version", apiVersion),
			option.WithJSONDel("model"),
		},
	}

	switch a.APIKey {
	case "":
		// Satisfy OpenAI's requirement for a non-empty string
		a.OpenAI.APIKey = "notused"
		// Inject bearer token middleware
		a.OpenAI.Opts = append(a.OpenAI.Opts, azure.WithTokenCredential(a.TokenCredential))

	default:
		// Use the "api-key" header instead of "Authorization"
		a.OpenAI.Opts = append(a.OpenAI.Opts,
			option.WithHeader("api-key", a.APIKey),
			option.WithHeaderDel("Authorization"))
	}
}
