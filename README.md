# Research Agent MCP Server

A **Model Context Protocol (MCP) server** written in Go that gives AI agents (Claude, Cursor, etc.) the ability to search the web, store research notes in Elasticsearch, and generate structured Markdown reports.

## Problem Statement

AI agents are powerful reasoners but have no persistent memory and cannot browse the web by default. This MCP server bridges that gap — it exposes a set of tools that agents can call to fetch real-time information, store findings in a searchable Elasticsearch index, and compile everything into a report.

## Tools Exposed

| Tool | Description |
|------|-------------|
| `web_search(query)` | Search the web via DuckDuckGo, returns titles + URLs |
| `fetch_page(url)` | Fetch and clean the text content of any web page |
| `store_note(title, content, tags)` | Index a research note into Elasticsearch |
| `search_notes(query)` | Full-text search across all stored notes |
| `generate_report(topic)` | Compile all notes into a Markdown report in `/data` |

## Stack

- **Go** — MCP protocol implementation (JSON-RPC 2.0 over stdio), concurrent HTTP fetching
- **Elasticsearch** — full-text search and note storage
- **Kibana** — visual dashboard to browse indexed notes
- **Docker + Docker Compose** — multi-container orchestration, single-command deploy

## Running with Docker Compose

```bash
docker compose up --build
```

Services:
- Elasticsearch → `http://localhost:9200`
- Kibana → `http://localhost:5601`
- MCP Server → connected via stdio (see Claude Desktop setup below)

## Connecting to Claude Desktop

Add this to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "research-agent": {
      "command": "docker",
      "args": ["exec", "-i", "research-mcp", "./mcp-server"],
      "env": {
        "ELASTICSEARCH_URL": "http://elasticsearch:9200"
      }
    }
  }
}
```

Or run locally without Docker (requires Elasticsearch running on port 9200):

```bash
go run ./cmd/server
```

## Example Usage

Ask Claude:
> "Research the pros and cons of microservices vs monolith architecture, store your findings, and generate a report."

Claude will:
1. Call `web_search` to find relevant articles
2. Call `fetch_page` on each URL
3. Call `store_note` to save key findings to Elasticsearch
4. Call `generate_report` to produce `data/report_<timestamp>.md`

## Project Structure

```
research-agent-mcp/
├── cmd/server/main.go              # entry point
├── internal/
│   ├── mcp/server.go               # MCP protocol (JSON-RPC over stdio)
│   ├── tools/                      # one file per tool
│   └── storage/elasticsearch.go    # ES client, index/search logic
├── data/                           # generated reports land here
├── Dockerfile                      # multi-stage build
└── docker-compose.yml              # ES + Kibana + MCP server
```
