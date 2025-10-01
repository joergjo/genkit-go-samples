# RAG Server Sample

## About
In 2024 the Go team published the article [Building LLM-powered applications in Go](https://go.dev/blog/llmpowered) on their blog. This article included a sample RAG server written in Go using the alpha version of Genkit Go. I've updated the sample to the GA version of Genkit. 

## Running the Sample
Follow the instructions in the [blog post](https://go.dev/blog/llmpowered).

If your editor or IDE supports `.http` [files](https://www.jetbrains.com/help/idea/exploring-http-syntax.html), you can use 
- [tests/add-documents.http](./tests/add-documents.http) to upload sample data to the document store and
- [tests/query](tests/query.http) to execute a sample query
