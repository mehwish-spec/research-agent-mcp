package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"research-agent-mcp/internal/storage"
)

// GenerateReport compiles all stored notes into a Markdown report
func GenerateReport(ctx context.Context, client *storage.Client, topic string) (string, error) {
	notes, err := client.ListAllNotes(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve notes: %w", err)
	}
	if len(notes) == 0 {
		return "No notes stored yet. Use store_note to add research before generating a report.", nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Research Report: %s\n\n", topic))
	sb.WriteString(fmt.Sprintf("*Generated on %s*\n\n", time.Now().Format("January 2, 2006 at 15:04")))
	sb.WriteString(fmt.Sprintf("**Total notes compiled:** %d\n\n---\n\n", len(notes)))

	for i, note := range notes {
		sb.WriteString(fmt.Sprintf("## %d. %s\n\n", i+1, note.Title))
		if len(note.Tags) > 0 {
			sb.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(note.Tags, ", ")))
		}
		sb.WriteString(note.Content + "\n\n")
		sb.WriteString("---\n\n")
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./data"
	}

	filename := fmt.Sprintf("report_%s.md", time.Now().Format("20060102_150405"))
	outPath := filepath.Join(dataDir, filename)

	if err := os.WriteFile(outPath, []byte(sb.String()), 0644); err != nil {
		return "", fmt.Errorf("failed to write report: %w", err)
	}

	return fmt.Sprintf("Report generated: %s (%d notes, %d bytes)", outPath, len(notes), len(sb.String())), nil
}
