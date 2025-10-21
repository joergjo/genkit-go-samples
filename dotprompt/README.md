# Dotprompt Sample

## About
This sample demonstrates externalizing prompt templates using Genkit Go's [Dotprompt](https://github.com/google/dotprompt) support.  The templates are stored as `.prompt` text files in the [`prompts`](./prompts/) directory and use the Handlebars template language, allowing Genkit Go to dynamically inject variables at runtime.

This sample is inspired by a C# sample from the book [_Building AI Applications with Microsoft Semantic Kernel_](https://www.packtpub.com/en-us/product/building-ai-applications-with-microsoft-semantic-kernel-9781835469590). The original C# source code can be found [here](https://github.com/PacktPublishing/Building-AI-Applications-with-Microsoft-Semantic-Kernel/blob/main/dotnet/ch2/ex03/Program.cs).     

## Running the Sample
```bash
cd dotprompt
export GEMINI_API_KEY=<your-api-key>

# Create a simple trip advice - pick any city you like
go run . -city "London"

# Create a detailed trip advice - set all flags shown below
go run . -city "London" -days 3 -attractions 5 -likes "Sports, Food, History" -dislikes "Art"
```