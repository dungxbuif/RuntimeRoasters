package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ElasticsearchClient struct {
	baseURL string
	index   string
	http    *http.Client
}

type TraceDocument struct {
	EntityID   string      `json:"entity_id"`
	EntityType string      `json:"entity_type"`
	Status     string      `json:"status"`
	Order      interface{} `json:"order,omitempty"`
	Payment    interface{} `json:"payment,omitempty"`
	Shipment   interface{} `json:"shipment,omitempty"`
	Timeline   interface{} `json:"timeline"`
	Events     interface{} `json:"events"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func NewElasticsearchClient(baseURL string, index string) *ElasticsearchClient {
	return &ElasticsearchClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		index:   index,
		http:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *ElasticsearchClient) EnsureIndex(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, c.indexURL(), nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("elasticsearch index check failed: status %d", resp.StatusCode)
	}

	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"entity_id": map[string]string{
					"type": "keyword",
				},
				"entity_type": map[string]string{
					"type": "keyword",
				},
				"status": map[string]string{
					"type": "keyword",
				},
				"updated_at": map[string]string{
					"type": "date",
				},
				"order": map[string]string{
					"type": "object",
				},
				"payment": map[string]string{
					"type": "object",
				},
				"shipment": map[string]string{
					"type": "object",
				},
				"timeline": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"topic":       map[string]string{"type": "keyword"},
						"status":      map[string]string{"type": "keyword"},
						"title":       map[string]string{"type": "text"},
						"description": map[string]string{"type": "text"},
						"timestamp":   map[string]string{"type": "date"},
						"payload":     map[string]interface{}{"type": "object", "enabled": false},
					},
				},
				"events": map[string]interface{}{
					"type": "nested",
					"properties": map[string]interface{}{
						"id": map[string]string{
							"type": "keyword",
						},
						"message_id": map[string]string{
							"type": "keyword",
						},
						"topic": map[string]string{
							"type": "keyword",
						},
						"batch_id": map[string]string{
							"type": "keyword",
						},
						"order_id": map[string]string{
							"type": "keyword",
						},
						"shipment_id": map[string]string{
							"type": "keyword",
						},
						"payload": map[string]string{
							"type": "text",
						},
						"occurred_at": map[string]string{
							"type": "date",
						},
						"created_at": map[string]string{
							"type": "date",
						},
					},
				},
			},
		},
	}
	body, err := json.Marshal(mapping)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPut, c.indexURL(), body, nil)
}

func (c *ElasticsearchClient) UpsertTraceDocument(ctx context.Context, document TraceDocument) error {
	body, err := json.Marshal(document)
	if err != nil {
		return err
	}
	return c.do(ctx, http.MethodPut, c.indexURL()+"/_doc/"+document.EntityID+"?refresh=true", body, nil)
}

func (c *ElasticsearchClient) GetTraceDocument(ctx context.Context, entityID string) (map[string]interface{}, bool, error) {
	var raw map[string]interface{}
	err := c.do(ctx, http.MethodGet, c.indexURL()+"/_doc/"+entityID, nil, &raw)
	if err != nil {
		if strings.Contains(err.Error(), "status 404") {
			return nil, false, nil
		}
		return nil, false, err
	}
	source, ok := raw["_source"].(map[string]interface{})
	if !ok {
		return nil, false, nil
	}
	return source, true, nil
}

func (c *ElasticsearchClient) indexURL() string {
	return c.baseURL + "/" + c.index
}

func (c *ElasticsearchClient) do(ctx context.Context, method string, url string, body []byte, out interface{}) error {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("elasticsearch request failed: method %s url %s status %d body %s", method, url, resp.StatusCode, string(responseBody))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}
