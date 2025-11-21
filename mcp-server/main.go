package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/beevik/ntp"
	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/mcp"
)

type NTPRequest struct {
	Host string `json:"host" jsonschema:"description=The NTP server hostname or IP address"`
}

func main() {
	ctx := context.Background()
	g := genkit.Init(ctx)
	genkit.DefineTool(g, "getTime", "Get the current time from an NTP server",
		func(ctx *ai.ToolContext, input NTPRequest) (string, error) {
			resp, err := ntp.Query(input.Host)
			if err != nil {
				return "", err
			}
			return resp.Time.Format(time.RFC3339), nil
		})
	server := mcp.NewMCPServer(g, mcp.MCPServerOptions{
		Name:    "ntp-mcp-go",
		Version: "0.0.1",
	})

	rc := 0
	slog.Info("starting MCP server on stdout")
	if err := server.ServeStdio(); err != nil {
		slog.Error("MCP server exited with error", "error", err)
		rc = 1
	}
	os.Exit(rc)
}
