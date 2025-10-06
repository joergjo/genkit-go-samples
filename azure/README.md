# Azure OpenAI Sample

## About
This sample shows hwo to build a custom [Azure OpenAI](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/overview) plugin for Genkit Go.

Why do even need a custom plugin? This plugin is just a thin adapter around the OpenAI plugin that wires up the Azure OpenAI endpoint. In order for this to work, you have to use the new Azure OpenAI endpoint documented [here](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/api-version-lifecycle?tabs=go), i.e. `https://YOUR-RESOURCE-NAME.openai.azure.com/openai/v1/`. 

> If you are using Azure OpenAI with its older endpoints based on deployments, have a look at [this version](https://github.com/joergjo/genkit-go-samples/blob/6ea363b6cb7564d0bb5fa8f46b9c183881d33b03/azure/azopenai.go) of the sample, which shows how to adapt the OpenAI plugin for these endpoints.

The sample uses GPT-5-mini, so make sure to deploy this model before running the sample.

The sample plugin also demonstrates Entra based access to Azure OpenAI instead of using API keys. Note that this only meant to be a proof of concept, since Genkit currently lacks any means for renewing short-lived credentials like OAuth2 access tokens. 

>Make sure the user principal accessing the API has been assigned the required roles: https://learn.microsoft.com/en-us/azure/ai-foundry/openai/how-to/managed-identity#assign-role.    

## Running the Sample
```bash
cd azure
export AZ_OPENAI_BASE_URL=<your-azure-openai-endpoint>
export AZ_OPENAI_API_KEY=<your-azure-openai-api-key>
# optional - if you want to log API requests sent to your endpoint.
export AZ_OPENAI_DEBUG_HTTP=true
go run .
```