package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const indexName = "research_notes"

type Note struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

type Client struct {
	es *elasticsearch.Client
}

func NewClient(address string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create ES client: %w", err)
	}
	return &Client{es: es}, nil
}

// EnsureIndex creates the index if it doesn't exist
func (c *Client) EnsureIndex(ctx context.Context) error {
	res, err := c.es.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		return nil
	}

	mapping := `{
		"mappings": {
			"properties": {
				"title":      { "type": "text" },
				"content":    { "type": "text" },
				"tags":       { "type": "keyword" },
				"created_at": { "type": "date" }
			}
		}
	}`
	req := esapi.IndicesCreateRequest{
		Index: indexName,
		Body:  strings.NewReader(mapping),
	}
	_, err = req.Do(ctx, c.es)
	return err
}

// IndexNote saves a note into Elasticsearch
func (c *Client) IndexNote(ctx context.Context, note Note) error {
	note.CreatedAt = time.Now()
	body, err := json.Marshal(note)
	if err != nil {
		return err
	}

	res, err := c.es.Index(
		indexName,
		bytes.NewReader(body),
		c.es.Index.WithDocumentID(note.ID),
		c.es.Index.WithContext(ctx),
	)
	if err != nil {
		return fmt.Errorf("index error: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("ES error: %s", res.String())
	}
	return nil
}

// SearchNotes performs a full-text search across title and content
func (c *Client) SearchNotes(ctx context.Context, query string) ([]Note, error) {
	var buf bytes.Buffer
	q := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^2", "content", "tags"},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return nil, err
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(indexName),
		c.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("search error: %w", err)
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source Note `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	notes := make([]Note, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		notes = append(notes, h.Source)
	}
	return notes, nil
}

// ListAllNotes returns all notes sorted by creation time
func (c *Client) ListAllNotes(ctx context.Context) ([]Note, error) {
	var buf bytes.Buffer
	q := map[string]interface{}{
		"query": map[string]interface{}{"match_all": map[string]interface{}{}},
		"sort":  []map[string]interface{}{{"created_at": "desc"}},
	}
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return nil, err
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(ctx),
		c.es.Search.WithIndex(indexName),
		c.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source Note `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	notes := make([]Note, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		notes = append(notes, h.Source)
	}
	return notes, nil
}
