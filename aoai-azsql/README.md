# Azure OpenAI Azure SQL Database Vector Sample

## About
This sample shows how to use the [Azure OpenAI sample plugin](./azure/) for embedding creation and vector search using Azure SQL's [native vector type and functions](https://learn.microsoft.com/en-us/sql/t-sql/data-types/vector-data-type?view=azuresqldb-current&tabs=csharp). It is based on the original [`pgvector` sample](https://github.com/firebase/genkit/tree/genkit%401.22.0/go/samples/pgvector) for Genkit Go that uses Google's `embedding-001` model. This sample uses Azure OpenAI's `text-embedding-3-small` instead.

Deploy `text-embedding-3-small` to your Azure OpenAI resource or Azure AI Foundry project and an Azure SQL database before running the sample. The sample uses the Azure OpenAI `v1` API, so make sure to specify the correct base URL (i.e., ending with `/openai/v1`). 

## Setting up Azure SQL
Execute the included [`vector.sql`](./vector.sql) script on your Azure database using a SQL client of your choice (`sqlcmd`, Visual Studio Code's [MSSQL extension](https://learn.microsoft.com/en-us/sql/tools/visual-studio-code-extensions/mssql/mssql-extension-visual-studio-code?view=sql-server-ver17), etc.). You can use the included [deployment script](./deploy.sh) to deploy the required Azure resources and sample data.

>The script assumes you have `bash` installed on either macOS, WSL2, or Linux, as well as OpenSSL and [sqlcmd](https://github.com/microsoft/go-sqlcmd). 

## Running the Sample
Open two terminal windows or tabs in your preferred terminal application.

### Run App
In window/tab #1

```bash
cd aoai-azsql

export AZ_OPENAI_BASE_URL=<your-azure-openai-endpoint>
export AZ_OPENAI_API_KEY=<your-azure-openai-api-key>
export GENKIT_ENV='dev'

# The init flag triggers the embedding generation
go run . -dbconn "sqlserver://<username>:<password>@<servername>.database.windows.net?database=<database-name>" -index
```

### Run Vector Search Flow 
In window/tab #2

```bash
genkit flow:run askQuestion '{"Show": "La Vie", "Question": "Who gets divorced?"}'
genkit flow:run askQuestion '{"Show": "Best Friends", "Question": "Who does Alice love?"}' 
```

The vector search result is written as log output to stderr (i.e., in terminal/tab #1).

![Vector search output](media/output.jpg)