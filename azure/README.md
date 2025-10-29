# Azure OpenAI Sample

## About
This sample shows how to build a custom [Azure OpenAI](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/overview) plugin for Genkit Go.

The sample plugin supports both Azure OpenAI [`v1`](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/api-version-lifecycle?tabs=go) and [`2024-10-21`](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/reference). Make sure to specify `AZ_OPENAI_BASE_URL` correctly——the `v1` base URL must include the path `openai/v1`. 

The sample uses GPT-5-mini, so make sure to deploy this model before running the sample. When using Azure OpenAI `2024-10-21`, you must set `AZ_OPENAI_DEPLOYMENT` to your model deployment's name.

The sample also demonstrates Entra-based access to Azure OpenAI instead of using API keys. The underlying OpenAI SDK client will automatically refresh access tokens because of the [middleware](https://github.com/openai/openai-go/blob/c5fd07f55034e2f14d3c3566d24973b903ad5761/azure/azure.go#L98) injected by the `azure` package.  

> Make sure the user principal accessing the API has been assigned the required roles: https://learn.microsoft.com/en-us/azure/ai-foundry/openai/how-to/managed-identity#assign-role.    

## Running the Sample
```bash
cd azure
export AZ_OPENAI_BASE_URL=<your-azure-openai-endpoint>
export AZ_OPENAI_API_KEY=<your-azure-openai-api-key>

# When using Azure OpenAI 2024-10-21
# export AZ_OPENAI_DEPLOYMENT=<your-model-deployment-name>

# optional - if you want to log API requests sent to your endpoint.
export AZ_OPENAI_DEBUG_HTTP=true
go run .
```