package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type OpenAIService struct {
	apiKey          string
	transcribeModel string
	chatModel       string
	httpClient      *http.Client
}

func NewOpenAIService(apiKey, transcribeModel, chatModel string) *OpenAIService {
	return &OpenAIService{
		apiKey:          apiKey,
		transcribeModel: transcribeModel,
		chatModel:       chatModel,
		httpClient: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}
func (s *OpenAIService) Transcribe(ctx context.Context, filePath string) (string, error) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	_ = writer.WriteField("model", s.transcribeModel)
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/audio/transcriptions", &body)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("transcribe http %d: %s", resp.StatusCode, string(respBytes))
	}

	var parsed struct {
		Text string `json:"text"`
	}

	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", err
	}

	text := strings.TrimSpace(parsed.Text)

	if text == "" {
		return "", fmt.Errorf("empty transcript")
	}

	return text, nil
}
func (s *OpenAIService) Chat(ctx context.Context, transcript string, prompt string) (string, error) {
	input := []map[string]any{}

	if strings.TrimSpace(prompt) != "" {
		input = append(input, map[string]any{
			"role": "developer",
			"content": []map[string]any{
				{"type": "input_text", "text": prompt},
			},
		})
	}

	input = append(input, map[string]any{
		"role": "user",
		"content": []map[string]any{
			{"type": "input_text", "text": transcript},
		},
	})

	payload := map[string]any{
		"model": s.chatModel,
		"input": input}

	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chat request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/responses", bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("responses http error %d: %s", resp.StatusCode, string(respBytes))
	}

	var parsed struct {
		Output []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	var sb strings.Builder

	for _, out := range parsed.Output {
		for _, c := range out.Content {
			if c.Type == "output_text" && c.Text != "" {
				sb.WriteString(c.Text)
			}
		}
	}

	answer := strings.TrimSpace(sb.String())
	if answer == "" {
		return "", fmt.Errorf("no output_text found in response")
	}

	return answer, nil
}
