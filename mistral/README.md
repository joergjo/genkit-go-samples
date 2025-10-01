# Mistral Sample

## About 
This sample demonstrates using Genkit Go with [Mistral AI](https://mistral.ai/) LLMs and embedding models. Genkit Go does not ship a Mistral AI plugin, but there is an [open source plugin by Thomas Marquis](github.com/thomas-marquis/genkit-mistral) which I'm using here. Thomas has also published an extensive article on [Genkit Go and RAG on Medium](https://medium.com/@thomas.marquis314/when-go-meets-ai-building-a-rag-application-with-genkit-3f0a2734eca7).

## Running the Sample
```bash
cd mistral
export MISTRAL_API_KEY=<your-api-key>
go run .
```