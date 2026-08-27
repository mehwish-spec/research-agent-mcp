package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ddgResult struct {
	AbstractText   string `json:"AbstractText"`
	AbstractURL    string `json:"AbstractURL"`
	AbstractSource string `json:"AbstractSource"`
	RelatedTopics  []struct {
		Text      string `json:"Text"`
		FirstURL  string `json:"FirstURL"`
	} `json:"RelatedTopics"`
}

// WebSearch queries DuckDuckGo Instant Answer API and returns results
func WebSearch(query string) (string, error) {
	apiURL := "https://api.duckduckgo.com/?q=" + url.QueryEscape(query) + "&format=json&no_redirect=1&no_html=1"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	var result ddgResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse search response: %w", err)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Search results for: %q\n\n", query))

	if result.AbstractText != "" {
		sb.WriteString(fmt.Sprintf("Summary: %s\nSource: %s\n\n", result.AbstractText, result.AbstractURL))
	}

	count := 0
	for _, topic := range result.RelatedTopics {
		if topic.FirstURL == "" || topic.Text == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("- %s\n  URL: %s\n", topic.Text, topic.FirstURL))
		count++
		if count >= 5 {
			break
		}
	}

	if count == 0 && result.AbstractText == "" {
		return fmt.Sprintf("No instant results found for %q. Try fetch_page with a direct URL.", query), nil
	}

	return sb.String(), nil
}
