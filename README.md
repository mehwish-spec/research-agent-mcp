# Research Agent MCP Server

A production-oriented **Model Context Protocol (MCP) server** built in Go that gives AI agents access to research, information retrieval, persistent notes, automated report generation, and isolated code execution.

The server communicates through **JSON-RPC 2.0 over stdio** and provides a collection of tools that can be consumed by MCP-compatible AI applications.

Features

* Web Search — Retrieve information from the web for research workflows.
* Persistent Notes — Store and retrieve research notes using Elasticsearch.
* Report Generation — Convert collected research into structured Markdown reports.
* Sandboxed Code Execution — Execute code inside isolated Docker containers.
* Resource Controls — Apply container-level CPU and memory constraints during execution.
* MCP Communication — JSON-RPC 2.0 communication over standard input/output.
* Containerized Infrastructure — Docker-based deployment for reproducible environments.
* Elasticsearch + Kibana — Persistent storage and visualization for research data.

 Architecture

                    ┌──────────────────────┐
                    │     AI / MCP Host    │
                    └──────────┬───────────┘
                               │
                         JSON-RPC / stdio
                               │
                    ┌──────────▼───────────┐
                    │    MCP Go Server     │
                    └──────────┬───────────┘
                               │
          ┌────────────────┼────────────────────┐
          │                    │                    │
          ▼                    ▼                    ▼
     Web Search          Research Notes       Report Generator
                              │
                              ▼
                       ┌─────────────┐
                       │Elasticsearch│
                       └──────┬──────┘
                              │
                              ▼
                          Kibana UI

                    ┌────────────────────┐
                    │  Code Execution    │
                    └─────────┬──────────┘
                              │
                              ▼
                       Docker Container
                       ┌───────────────┐
                         Isolated Code │
                       │   Execution   │
                       └───────────────┘


 Technology Stack

| Technology         | Purpose                                |
| ------------------ | -------------------------------------- |
| Go                 | MCP server and backend implementation  |
| MCP                | Tool integration with AI agents        |
| JSON-RPC 2.0       | Communication protocol                 |
| Elasticsearch      | Persistent research/note storage       |
| Kibana             | Data visualization and inspection      |
| Docker             | Isolated code execution and deployment |
| Docker Compose     | Local multi-service orchestration      |
| Markdown           | Generated research reports             |

Project Structure


research-agent-mcp/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│  ├── mcp/
│  │   └── server.go
│  │
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


Getting Started

Prerequisites

Make sure the following are installed:

* Go
* Docker
* Docker Compose
* Git

Clone the Repository

git clone https://github.com/mehwish-spec/research-agent-mcp.git
cd research-agent-mcp

Start the Services

docker compose up --build

This starts the required containerized services used by the project.

Run the Server

For local development:

go run ./cmd/server

MCP Tools

The server exposes tools that allow an MCP-compatible AI client to perform research-related operations.

Search

Search for relevant information and return results that can be used during a research workflow.

Fetch

Retrieve content from a specified resource for further processing.

Notes

Persist research notes in Elasticsearch and retrieve them when needed.

Report

Generate structured Markdown reports from gathered research information.

Run Code

Execute code inside a temporary Docker-based sandbox with resource restrictions.

Sandboxed Code Execution

The code execution component is designed around process isolation.

For each execution request, the system can:

1. Create an ephemeral Docker container.
2. Apply CPU and memory restrictions.
3. Execute the requested code inside the container.
4. Capture the execution output.
5. Return the result.
6. Remove the temporary execution environment.

This architecture keeps arbitrary code execution separated from the main application process.

Elasticsearch & Kibana

Elasticsearch provides persistent storage for research-related information such as notes and indexed content.

Kibana can be used to inspect and visualize the stored data during development and debugging.

The services can be started together through Docker Compose.

Use Cases

This server can be used as the backend for AI-powered workflows such as:

* Automated research assistants
* Technical research
* Information gathering
* Research note management
* Automated report generation
* AI-assisted data analysis
* Controlled code execution
* Multi-step agent workflows

Research Workflow

A typical workflow can look like:

User Question
      │
      ▼
   AI Agent
      │
      ▼
  MCP Server
      │
      ├── Search Web
      │
      ├── Fetch Information
      │
      ├── Save Notes
      │
      ├── Execute Analysis
      │
      └── Generate Report
              │
              ▼
       Markdown Report

Engineering Highlights

* Implemented an MCP-compatible server in Go.
* Used JSON-RPC 2.0 for structured tool communication.
* Integrated Elasticsearch for persistent storage.
* Containerized infrastructure with Docker.
* Designed isolated execution environments for running code.
* Added resource constraints for containerized workloads.
* Organized functionality into modular MCP tools.
* Supported automated Markdown report generation.

Future Improvements

Potential extensions include:

* Authentication and authorization
* More granular container security policies
* Additional research tools
* Streaming tool responses
* Improved observability and logging
* Persistent execution history
* Automated test coverage
* Production deployment configuration

License

See the repository license for usage and distribution terms.
