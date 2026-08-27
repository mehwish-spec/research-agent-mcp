package tools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"research-agent-mcp/internal/storage"
)

func randomID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// StoreNote indexes a note into Elasticsearch
func StoreNote(ctx context.Context, client *storage.Client, title, content string, tags []string) (string, error) {
	note := storage.Note{
		ID:      randomID(),
		Title:   title,
		Content: content,
		Tags:    tags,
	}
	if err := client.IndexNote(ctx, note); err != nil {
		return "", fmt.Errorf("failed to store note: %w", err)
	}
	return fmt.Sprintf("Note stored successfully. ID: %s, Title: %q, Tags: [%s]",
		note.ID, title, strings.Join(tags, ", ")), nil
}

// SearchNotes performs full-text search in Elasticsearch
func SearchNotes(ctx context.Context, client *storage.Client, query string) (string, error) {
	notes, err := client.SearchNotes(ctx, query)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}
	if len(notes) == 0 {
		return "No notes found matching: " + query, nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d note(s) for query %q:\n\n", len(notes), query))
	for i, n := range notes {
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n   Tags: %s\n   %s\n\n",
			i+1, n.ID, n.Title, strings.Join(n.Tags, ", "), truncate(n.Content, 200)))
	}
	return sb.String(), nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
