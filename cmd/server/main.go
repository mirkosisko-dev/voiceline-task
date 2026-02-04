package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mirkosisko-dev/voiceline/internal/config"
	"github.com/mirkosisko-dev/voiceline/internal/handler"
	"github.com/mirkosisko-dev/voiceline/internal/middleware"
	"github.com/mirkosisko-dev/voiceline/internal/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}
	cfg := config.Load()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	maxFileSize := cfg.FileMetadata.MaxFileSize
	if maxFileSize <= 0 {
		maxFileSize = 10 << 20
	}

	r.Use(middleware.MaxBodyBytes(maxFileSize))

	r.MaxMultipartMemory = maxFileSize

	openaiSvc := service.NewOpenAIService(cfg.OpenAI.APIKey, cfg.OpenAI.TranscribeModel, cfg.OpenAI.ChatModel)

	sinkTimeout := time.Duration(cfg.Sink.TimeoutSeconds) * time.Second
	sinkSvc := service.NewSinkService(cfg.Sink.WebhookURL, sinkTimeout)

	audioSvc := service.NewAudioService(openaiSvc, sinkSvc, cfg.Sink.FailRequest)

	fmt.Printf("running on: %s\n", cfg.Server.Port)

	audioHandler := handler.NewAudioHandler(audioSvc, cfg)
	healthHandler := handler.NewHealthHandler()

	r.POST("/audio", audioHandler.HandleAudioUpload)
	r.GET("/health", healthHandler.HealthCheck)

	r.Run(":" + cfg.Server.Port)
}
