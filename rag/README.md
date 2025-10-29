# RAG Sample

## About 
The Genkit documentation contains a chapter on [retrieval augmented generation (RAG)](https://genkit.dev/docs/rag/?lang=go) with a sample application, but at the time of writing, this sample is incomplete and even [lacks the PDF file](https://github.com/firebase/genkit/issues/1405) the application is meant to use. I've filled in the missing parts and adopted the sample to use the "Margies Travels" sample data used by various [Microsoft Learn](https://learn.microsoft.com/en-us/) courses. 

## Running the Sample
Open two terminal windows or tabs in your preferred terminal application.

### Run App
In window/tab #1

```bash
cd rag
export GEMINI_API_KEY=<your-api-key>
export GENKIT_ENV=dev
go run .
```

### Create Embeddings
In window/tab #2

```bash
cd rag
genkit flow:run indexBrochure '"brochures/Dubai Brochure.pdf"'
genkit flow:run indexBrochure '"brochures/Las Vegas Brochure.pdf"'
genkit flow:run indexBrochure '"brochures/London Brochure.pdf"'
genkit flow:run indexBrochure '"brochures/Margies Travel Company Info.pdf"'
genkit flow:run indexBrochure '"brochures/New York Brochure.pdf"'
genkit flow:run indexBrochure '"brochures/San Francisco Brochure.pdf"'
```

### Query Documents
In window/tab #2

```bash
genkit flow:run travelQA '"Where can I stay in NY?"'
genkit flow:run travelQA '"When is a great time to visit London?"'
```