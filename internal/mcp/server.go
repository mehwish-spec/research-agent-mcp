package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"research-agent-mcp/internal/storage"
	"research-agent-mcp/internal/tools"
)

// JSON-RPC 2.0 types
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// Server is the MCP server that reads from stdin and writes to stdout
type Server struct {
	esClient *storage.Client
}

func NewServer(esClient *storage.Client) *Server {
	return &Server{esClient: esClient}
}

func (s *Server) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.writeError(nil, -32700, "parse error")
			continue
		}

		s.handle(req)
	}
}

func (s *Server) handle(req Request) {
	switch req.Method {
	case "initialize":
		s.writeResult(req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
			"serverInfo":      map[string]interface{}{"name": "research-agent-mcp", "version": "1.0.0"},
		})

	case "tools/list":
		s.writeResult(req.ID, map[string]interface{}{
			"tools": toolDefinitions(),
		})

	case "tools/call":
		var params ToolCallParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			s.writeError(req.ID, -32602, "invalid params")
			return
		}
		result, err := s.executeTool(params.Name, params.Arguments)
		if err != nil {
			s.writeResult(req.ID, map[string]interface{}{
				"content": []map[string]interface{}{{"type": "text", "text": "Error: " + err.Error()}},
				"isError": true,
			})
			return
		}
		s.writeResult(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": result}},
		})

	default:
		s.writeError(req.ID, -32601, "method not found: "+req.Method)
	}
}

func (s *Server) executeTool(name string, args map[string]interface{}) (string, error) {
	ctx := context.Background()

	switch name {
	case "web_search":
		query, _ := args["query"].(string)
		return tools.WebSearch(query)

	case "fetch_page":
		url, _ := args["url"].(string)
		return tools.FetchPage(url)

	case "store_note":
		title, _ := args["title"].(string)
		content, _ := args["content"].(string)
		var tags []string
		if raw, ok := args["tags"].([]interface{}); ok {
			for _, t := range raw {
				if s, ok := t.(string); ok {
					tags = append(tags, s)
				}
			}
		}
		return tools.StoreNote(ctx, s.esClient, title, content, tags)

	case "search_notes":
		query, _ := args["query"].(string)
		return tools.SearchNotes(ctx, s.esClient, query)

	case "generate_report":
		topic, _ := args["topic"].(string)
		return tools.GenerateReport(ctx, s.esClient, topic)

	case "run_code":
		language, _ := args["language"].(string)
		code, _ := args["code"].(string)
		return tools.RunCode(language, code)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func (s *Server) writeResult(id interface{}, result interface{}) {
	s.write(Response{JSONRPC: "2.0", ID: id, Result: result})
}

func (s *Server) writeError(id interface{}, code int, msg string) {
	s.write(Response{JSONRPC: "2.0", ID: id, Error: &RPCError{Code: code, Message: msg}})
}

func (s *Server) write(resp Response) {
	b, _ := json.Marshal(resp)
	fmt.Fprintf(os.Stdout, "%s\n", b)
}

func toolDefinitions() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "web_search",
			"description": "Search the web using DuckDuckGo and return relevant results with URLs",
			"inputSchema": schema(map[string]interface{}{
				"query": prop("string", "The search query"),
			}, []string{"query"}),
		},
		{
			"name":        "fetch_page",
			"description": "Fetch the content of a web page and return cleaned plain text",
			"inputSchema": schema(map[string]interface{}{
				"url": prop("string", "The URL to fetch"),
			}, []string{"url"}),
		},
		{
			"name":        "store_note",
			"description": "Store a research note in Elasticsearch for later retrieval",
			"inputSchema": schema(map[string]interface{}{
				"title":   prop("string", "Short title for the note"),
				"content": prop("string", "The full content of the note"),
				"tags":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Tags to categorize this note"},
			}, []string{"title", "content"}),
		},
		{
			"name":        "search_notes",
			"description": "Full-text search across all stored research notes in Elasticsearch",
			"inputSchema": schema(map[string]interface{}{
				"query": prop("string", "Search query to find relevant notes"),
			}, []string{"query"}),
		},
		{
			"name":        "generate_report",
			"description": "Compile all stored notes into a structured Markdown report file",
			"inputSchema": schema(map[string]interface{}{
				"topic": prop("string", "The topic/title for the report"),
			}, []string{"topic"}),
		},
		{
			"name":        "run_code",
			"description": "Execute code in an isolated Docker container and return the output. Supported languages: python, javascript, go, bash",
			"inputSchema": schema(map[string]interface{}{
				"language": prop("string", "Programming language: python, javascript, go, or bash"),
				"code":     prop("string", "The code to execute"),
			}, []string{"language", "code"}),
		},
	}
}

func schema(props map[string]interface{}, required []string) map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": props,
		"required":   required,
	}
}

func prop(t, desc string) map[string]interface{} {
	return map[string]interface{}{"type": t, "description": desc}
}
