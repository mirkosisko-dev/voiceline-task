package service

import (
	"context"
	"time"
)

type AudioService struct {
	openAIService *OpenAIService
	sinkService   *SinkService
	failOnSink    bool
}

func NewAudioService(openAIService *OpenAIService, sinkService *SinkService, failOnSink bool) *AudioService {
	return &AudioService{
		openAIService: openAIService,
		sinkService:   sinkService,
		failOnSink:    failOnSink,
	}
}

type TranscriptResult struct {
	Transcript string
	Answer     string
	Prompt     string
}

func (s *AudioService) ProcessAudio(ctx context.Context, filePath string, prompt string) (*TranscriptResult, error) {
	transcript, err := s.openAIService.Transcribe(ctx, filePath)
	if err != nil {
		return nil, err
	}

	answer, err := s.openAIService.Chat(ctx, transcript, prompt)
	if err != nil {
		return nil, err
	}

	if s.sinkService != nil {
		subErr := s.sinkService.Submit(ctx, SinkPayload{
			Transcript: transcript,
			Answer:     answer,
			Prompt:     prompt,
			CreatedAt:  time.Now().UTC(),
		})

		if subErr != nil && s.failOnSink {
			return nil, subErr
		}
	}

	return &TranscriptResult{
		Transcript: transcript,
		Answer:     answer,
		Prompt:     prompt,
	}, nil
}
