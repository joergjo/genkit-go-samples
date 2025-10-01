# Summarize Video Sample

## About
This is my Go version of the JavaScriot tutorial in Genkit's documentation: https://genkit.dev/docs/tutorials/summarize-youtube-videos/. At the time of writing, there was no Go version available, but it was very simple to adapt the JavaScript version.

## Running the Sample
```bash
cd summarize-video
export GEMINI_API_KEY=<your-api-key>

# Use the default prompt
go run . -url "https://www.youtube.com/watch?v=kj80m-umOxs&t=2s"

# Use your own prompt
go run . -url "https://www.youtube.com/watch?v=YUgXJkNqH9Q" -prompt "Please provide a concise summary of the video segments that pertain to Genkit"
```