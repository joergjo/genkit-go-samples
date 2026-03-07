# Sessions Sample

## About
This sample demonstrates Genkit Go's session management using the experimental [`core/x/session`](https://pkg.go.dev/github.com/firebase/genkit/go/core/x/session) package. It uses an in-memory session store to persist chat history across multiple calls to a `"chat"` flow. The second prompt references the first response, showing that session state is preserved between calls.

## Running the Sample
```bash
cd sessions
export GEMINI_API_KEY=<your-api-key>
go run .
```
