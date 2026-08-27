package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"research-agent-mcp/internal/mcp"
	"research-agent-mcp/internal/storage"
)

func main() {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200"
	}

	esClient, err := storage.NewClient(esURL)
	if err != nil {
		log.Fatalf("failed to connect to Elasticsearch at %s: %v", esURL, err)
	}

	if err := esClient.EnsureIndex(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not ensure index: %v\n", err)
	}

	log.Printf("Research Agent MCP Server started (ES: %s)", esURL)
	server := mcp.NewServer(esClient)
	server.Run()
}
