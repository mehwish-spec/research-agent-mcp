# Research Agent MCP Server

A Go-based Model Context Protocol (MCP) server designed to extend AI assistants with web research, persistent research notes, full-text search, and automated Markdown report generation.

## Overview

Research Agent MCP provides a bridge between AI assistants and external research infrastructure.

An AI assistant can use the MCP server to:

- Search the web for relevant information
- Retrieve and clean content from web pages
- Save useful findings as structured research notes
- Search previously stored research
- Organize collected information by topic and tags
- Generate a structured Markdown report from stored research

The project uses Go for the MCP server, Elasticsearch for persistent research storage and search, and Docker Compose for running the complete development environment.

## Why This Project?

AI assistants can reason over information provided to them, but research workflows often require multiple external steps:

1. Finding relevant sources
2. Reading web pages
3. Extracting useful information
4. Remembering findings across multiple research steps
5. Searching previously collected information
6. Converting the collected findings into a structured report

This project combines these capabilities into an MCP-compatible server so that an AI client can interact with the research system through a defined set of tools.

## Architecture

The system consists of three primary components:

### MCP Server

The MCP server is implemented in Go and communicates with compatible AI clients through the Model Context Protocol.

It exposes research-related tools and handles communication using JSON-RPC over standard input and output.

### Elasticsearch

Elasticsearch provides persistent storage for research notes.

Each stored note can contain information such as:

- Title
- Content
- Tags
- Research metadata

The stored information can later be searched using full-text queries.

### Kibana

Kibana provides a visual interface for inspecting and exploring the data stored in Elasticsearch.

### Docker Compose

Docker Compose is used to run the MCP server, Elasticsearch, and Kibana together as a development environment.

## Features

### Web Search

The `web_search` tool searches the web and returns relevant search results.

```text
web_search(query)
```

It is useful for discovering sources related to a research topic before retrieving their full content.

### Web Page Fetching

The `fetch_page` tool retrieves and cleans the text content of a web page.

```text
fetch_page(url)
```

This allows the AI client to work with readable page content rather than relying only on search-result metadata.

### Research Note Storage

The `store_note` tool saves research findings into Elasticsearch.

```text
store_note(title, content, tags)
```

Research notes can be organized using tags, making it easier to build a searchable research knowledge base.

### Research Search

The `search_notes` tool performs full-text searches over previously stored research.

```text
search_notes(query)
```

This allows information collected during earlier research steps to be retrieved and reused.

### Report Generation

The `generate_report` tool creates a structured Markdown report from the available research notes.

```text
generate_report(topic)
```

Generated reports are written to the project's data directory.

## Research Workflow

A typical research workflow can be represented as:

```text
AI Assistant
     |
     v
MCP Server
     |
     +------------------+
     |                  |
     v                  v
Web Search          Stored Research
     |                  |
     v                  v
Fetch Pages        Elasticsearch
     |                  |
     +--------+---------+
              |
              v
       Research Notes
              |
              v
       Report Generation
              |
              v
       Markdown Report
```

For example, an AI assistant can receive a research request and perform the following sequence:

```text
1. Search for relevant sources
2. Retrieve selected web pages
3. Extract useful information
4. Store important findings
5. Search the accumulated notes
6. Generate a structured Markdown report
```

## Technology Stack

### Backend

- Go
- Model Context Protocol
- JSON-RPC
- HTTP-based web fetching

### Storage and Search

- Elasticsearch
- Full-text search
- Structured research notes

### Observability and Data Exploration

- Kibana

### Infrastructure

- Docker
- Docker Compose

## Project Structure

```text
research-agent-mcp/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── mcp/
│   │   └── server.go
│   │
│   ├── storage/
│   │   └── elasticsearch.go
│   │
│   └── tools/
│       ├── fetch.go
│       ├── notes.go
│       ├── report.go
│       ├── runcode.go
│       └── search.go
│
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Directory Responsibilities

### `cmd/server`

Contains the application's entry point.

`main.go` initializes the MCP server and starts the application.

### `internal/mcp`

Contains the MCP server implementation and protocol-level handling.

### `internal/tools`

Contains the individual capabilities exposed through MCP.

The tool implementations include:

- Web page retrieval
- Web search
- Research note management
- Research search
- Report generation
- Code execution functionality

### `internal/storage`

Contains Elasticsearch-related functionality, including connection and research data operations.

### Docker Configuration

`Dockerfile` defines the container image for the Go application.

`docker-compose.yml` defines the services required to run the complete environment.

## Running the Project

### Prerequisites

Install the following before running the project:

- Go
- Docker
- Docker Compose

Verify the installations:

```bash
go version
docker --version
docker compose version
```

## Run With Docker Compose

Build and start the services:

```bash
docker compose up --build
```

The environment provides the following services:

| Service | Address |
|---|---|
| Elasticsearch | `http://localhost:9200` |
| Kibana | `http://localhost:5601` |
| MCP Server | Standard input/output |

To stop the services:

```bash
docker compose down
```

To rebuild the containers after making code changes:

```bash
docker compose up --build
```

## Run the Go Server Locally

If Elasticsearch is already running locally on port `9200`, the MCP server can also be started directly with:

```bash
go run ./cmd/server
```

The server communicates through standard input and output as required by the MCP architecture.

## Connecting to an MCP Client

The server can be configured as an MCP server in a compatible AI client.

For example, a Docker-based configuration can launch the MCP process inside the running container:

```json
{
  "mcpServers": {
    "research-agent": {
      "command": "docker",
      "args": [
        "exec",
        "-i",
        "research-mcp",
        "./mcp-server"
      ],
      "env": {
        "ELASTICSEARCH_URL": "http://elasticsearch:9200"
      }
    }
  }
}
```

The exact configuration path and integration method depend on the MCP client being used.

## Example Research Session

A research request could look like:

```text
Research the advantages and disadvantages of
microservices compared with monolithic architecture.
Save the important findings and generate a report.
```

The research agent can then use the available tools to:

```text
web_search
     |
     v
Find relevant sources
     |
     v
fetch_page
     |
     v
Extract source content
     |
     v
store_note
     |
     v
Save findings in Elasticsearch
     |
     v
search_notes
     |
     v
Retrieve relevant research
     |
     v
generate_report
     |
     v
Generate Markdown report
```

## Generated Reports

Reports are generated in Markdown format.

A generated report can contain:

- Research topic
- Collected findings
- Source information
- Organized research notes
- Structured conclusions

The generated output can then be used for further analysis, documentation, or review.

## Elasticsearch Data Layer

Elasticsearch acts as the persistent research layer of the system.

Instead of keeping research findings only inside a single AI conversation, the MCP server can store notes and make them searchable.

This provides several useful capabilities:

- Persistent research storage
- Full-text search
- Tag-based organization
- Retrieval of previously collected information
- Reuse of research across different tasks

## Development

Install the Go dependencies:

```bash
go mod download
```

Run the application:

```bash
go run ./cmd/server
```

Build the application:

```bash
go build ./...
```

Run the project's Go tests, when available:

```bash
go test ./...
```

## Configuration

The Elasticsearch connection can be configured using the environment variable:

```text
ELASTICSEARCH_URL
```

For a Docker Compose environment, the Elasticsearch service can be referenced using:

```text
http://elasticsearch:9200
```

For a local Elasticsearch installation, the URL can typically be:

```text
http://localhost:9200
```

## Design Goals

The project is designed around several principles:

### Tool-Based AI Interaction

Instead of embedding every research capability directly into an AI application, functionality is exposed through MCP tools.

This allows compatible AI clients to discover and use research capabilities through a standardized interface.

### Persistent Research

Research findings are stored externally in Elasticsearch rather than existing only within a temporary conversation.

### Modular Components

Search, fetching, storage, retrieval, and report generation are implemented as separate tools and components.

### Containerized Development

Docker Compose makes it possible to run the main services together without manually configuring every dependency.

## Potential Extensions

The architecture can be extended with additional capabilities such as:

- Additional search providers
- More advanced document extraction
- Source credibility scoring
- Research deduplication
- Citation generation
- PDF report generation
- Additional storage backends
- Authentication and access control
- Background research jobs
- Research history and versioning
- More advanced Elasticsearch indexing
- Automated testing and CI/CD

## Limitations

This project is intended primarily as a research-agent infrastructure and development project.

Search quality, extracted page content, and generated reports depend on the external sources and the behavior of the configured tools.

For production use, additional considerations would be required around:

- Authentication
- Rate limiting
- Error handling
- Security
- Resource limits
- Search provider reliability
- Web content validation
- Elasticsearch security configuration
- Production deployment

## License

Add the project's license information here if a license has been selected.

## Project Status

This project is intended as an extensible foundation for building AI-powered research workflows with the Model Context Protocol.

The architecture separates the AI-facing MCP interface from web retrieval, persistent research storage, search, and report generation, making the system suitable for further experimentation and development.