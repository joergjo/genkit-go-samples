# Genkit Go Samples

## About
This repository contains code samples for the [Genkit Go framework](https://genkit.dev/docs/get-started/?lang=go). All samples use the GA version of Genkit Go. The samples work with free tiers of Gemini, Mistral, etc. For Azure OpenAI, you can get a subscription for free [here](https://azure.microsoft.com/en-us/free/). 

Some of the samples require the use of the [Genkit CLI](https://genkit.dev/docs/devtools/?lang=go#command-line-interface-cli-1).

>Genkit supports macOS, Windows, and Linux. These samples have been built and tested on macOS 26 Tahoe.

## Table of Contents
[aoai-pgvector](./aoai-pgvector/): This sample shows how to use the [Azure OpenAI sample plugin](./azure/) for embedding creation and vector search using PostgreSQL and the [pgvector](https://github.com/pgvector/pgvector) extension.

[azure](./azure/): This sample demonstrates building a custom [Azure OpenAI](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/overview) plugin based on Genkit Go's [OpenAI plugin](https://genkit.dev/docs/integrations/openai/?lang=go).

[code-agent](./code-agent/): A basic CLI coding agent using Genkit Go [flows](https://genkit.dev/docs/flows/?lang=go) and [tools](https://genkit.dev/docs/tool-calling/?lang=go). 

[dotprompt](./dotprompt/): This sample demonstrates externalizing prompt templates using Genkit Go's [Dotprompt](https://github.com/google/dotprompt) support.  

[mcp-server](./mcp-server/): This sample demonstrates building a simple MPC server with Genkit Go and integrating it with [Claude Desktop](https://claude.ai/download). 

[mistral](./mistral/): This sample demonstrates using the [genkit-mistral](https://pkg.go.dev/github.com/thomas-marquis/genkit-mistral) plugin.

[rag](./rag/): A RAG (retrieval augmented generation) sample that demonstrates using Genkit Go's Dev Local Vector Store based on Genkit Go's RAG documentation.

[rag-server](./rag-server/): This sample is an updated version of the [demo application](https://github.com/golang/example/tree/master/ragserver/ragserver-genkit) published by the Go team for Genkit Go's alpha version.

[tools](./tools/): This sample demonstrates Genkit Go's [tool calling](https://genkit.dev/docs/tool-calling/?lang=go) capabilities based on a fictitious home automation scenario. 

[summarize-video](./summarize-video/): A Go version of the [JavaScript tutorial](https://genkit.dev/docs/tutorials/summarize-youtube-videos/) published by the Genkit team.

## Other Samples
I've also published a Go SDK for Microsoft's Foundry Local. An example for using Genkit Go with Foundry Local is in that repo's [example folder](https://github.com/joergjo/go-foundry-local/tree/main/examples/genkit-go). 