package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type SinkService struct {
	webhookURL string
	client     *http.Client
}

func NewSinkService(webhookURL string, timeout time.Duration) *SinkService {
	return &SinkService{
		webhookURL: strings.TrimSpace(webhookURL),
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

type SinkPayload struct {
	Transcript string    `json:"transcript"`
	Answer     string    `json:"answer"`
	Prompt     string    `json:"prompt"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *SinkService) Submit(ctx context.Context, payload SinkPayload) error {
	if s == nil || s.webhookURL == "" {
		return nil // disabled
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("sink marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhookURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("sink create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("sink send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("sink http error: %d", resp.StatusCode)
	}

	return nil
}
