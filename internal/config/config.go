package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string
	}
	OpenAI struct {
		APIKey          string
		TranscribeModel string
		ChatModel       string
	}
	FileMetadata struct {
		MaxFileSize int64
	}
	Sink struct {
		WebhookURL     string
		TimeoutSeconds int64
		FailRequest    bool
	}
}

func Load() *Config {
	return &Config{
		Server: struct {
			Port string
		}{Port: getEnv("PORT", "8080")},
		OpenAI: struct {
			APIKey          string
			TranscribeModel string
			ChatModel       string
		}{
			APIKey:          getEnv("OPENAI_API_KEY", ""),
			TranscribeModel: getEnv("OPENAI_TRANSCRIBE_MODEL", "gpt-4o-transcribe"),
			ChatModel:       getEnv("OPENAI_CHAT_MODEL", "gpt-4o-mini"),
		},
		FileMetadata: struct {
			MaxFileSize int64
		}{
			MaxFileSize: getEnvInt("MAX_FILE_SIZE", 15728640),
		},
		Sink: struct {
			WebhookURL     string
			TimeoutSeconds int64
			FailRequest    bool
		}{
			WebhookURL:     getEnv("SINK_WEBHOOK_URL", ""),
			TimeoutSeconds: getEnvInt("SINK_WEBHOOK_TIMEOUT_SECONDS", 10),
			FailRequest:    getEnvBool("SINK_FAIL_REQUEST", false),
		},
	}
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return defaultValue
		}

		return i
	}

	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return defaultValue
		}
		return b
	}
	return defaultValue
}
