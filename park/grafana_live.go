package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// GrafanaLiveClient handles streaming data to Grafana Live
type GrafanaLiveClient struct {
	baseURL  string
	apiKey   string
	streamID string
	client   *http.Client
	enabled  bool
}

// NewGrafanaLiveClient creates a new Grafana Live client
func NewGrafanaLiveClient(grafanaURL, apiKey string) *GrafanaLiveClient {
	enabled := grafanaURL != "" && apiKey != ""

	if !enabled {
		slog.Info("Grafana Live streaming disabled - missing URL or API key")
	} else {
		slog.Info("Grafana Live streaming enabled", "url", grafanaURL, "api_key_length", len(apiKey))
	}

	return &GrafanaLiveClient{
		baseURL:  grafanaURL,
		apiKey:   apiKey,
		streamID: "kubepark",
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		enabled: enabled,
	}
}

// PushMetric sends a metric to Grafana Live
func (g *GrafanaLiveClient) PushMetric(metricName string, value float64, tags map[string]string) error {
	if !g.enabled {
		return nil // Silently skip if not enabled
	}

	// Build tags string
	tagsStr := "source=kubepark"
	for k, v := range tags {
		tagsStr += fmt.Sprintf(",%s=%s", k, v)
	}

	// Format data in InfluxDB line protocol
	now := time.Now().UnixNano()
	data := fmt.Sprintf("%s,%s value=%.2f %d\n", metricName, tagsStr, value, now)

	// Create HTTP request to Grafana Live Push API
	url := fmt.Sprintf("%s/api/live/push/%s", g.baseURL, g.streamID)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "text/plain")

	// Send request
	resp, err := g.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		// Read response body for more details
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("grafana live push failed with status: %d, response: %s", resp.StatusCode, string(body))
	}

	slog.Debug("Successfully pushed metric to Grafana Live",
		"metric", metricName,
		"value", value,
		"tags", tags,
		"status", resp.StatusCode)

	return nil
}
