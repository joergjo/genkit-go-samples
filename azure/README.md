# Azure OpenAI Sample

## About
This sample shows hwo to build a custom [Azure OpenAI](https://learn.microsoft.com/en-us/azure/ai-foundry/openai/overview) plugin for Genkit Go.

Why do even need a custom plugin? Genkit Go includes both an [OpenAI plugin](https://genkit.dev/docs/integrations/openai/?lang=go) and an [OpenAI-Compatible Plugin](https://genkit.dev/docs/integrations/openai-compatible/?lang=go). While Azure OpenAI's REST API is almost identical to the OpenAI REST API, there is one fundamental difference: OpenAI is model centric, whereas Azure OpenAI is endpoint centric using *deployments*. Hence, in Azure OpenAI the deployment abstracts the model and version being used. Since the deployment is an integral part of the endpoint's URL, Genkit Go must be adapted accordingly. 

Note that the sample still requires specifying the model name as well. This is because Genkit Go maintains an internal catalog of OpenAI models and their capabilities, which this implementation reuses. 


The sample plugin also demonstrates Entra based access to Azure OpenAI instead of using API keys. Note that this a proof of concept, since Genkit right now lacks any concept for renewing credentials like OAuth2 access tokens. 

>Make sure the user principal accessing the API has been assigned the required roles: https://learn.microsoft.com/en-us/azure/ai-foundry/openai/how-to/managed-identity#assign-role.    

## Running the Sample
```bash
cd azure
export AZ_OAI_ENDPOINT=<your-azure-openai-endpoint>
export AZ_OAI_DEPLOYMENT_NAME=<your-azure-oepnai-deployment>
export AZ_OAI_MODEL_NAME=<your-azure-openai-model>
export AZ_OAI_API_KEY=<your-azure-openai-api-key>
go run .
```