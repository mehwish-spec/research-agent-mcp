package tools

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	htmlTagRe    = regexp.MustCompile(`<[^>]+>`)
	multiSpaceRe = regexp.MustCompile(`\s{2,}`)
)

// FetchPage retrieves a URL and returns cleaned plain text
func FetchPage(url string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ResearchBot/1.0)")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024)) // cap at 512KB
	if err != nil {
		return "", err
	}

	text := stripHTML(string(body))
	if len(text) > 4000 {
		text = text[:4000] + "\n\n[content truncated]"
	}
	return text, nil
}

func stripHTML(html string) string {
	// remove script and style blocks entirely
	scriptRe := regexp.MustCompile(`(?is)<(script|style)[^>]*>.*?</(script|style)>`)
	html = scriptRe.ReplaceAllString(html, "")
	html = htmlTagRe.ReplaceAllString(html, " ")
	html = strings.ReplaceAll(html, "&nbsp;", " ")
	html = strings.ReplaceAll(html, "&amp;", "&")
	html = strings.ReplaceAll(html, "&lt;", "<")
	html = strings.ReplaceAll(html, "&gt;", ">")
	html = multiSpaceRe.ReplaceAllString(html, " ")
	return strings.TrimSpace(html)
}
