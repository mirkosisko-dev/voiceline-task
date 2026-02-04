package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	OpenAI struct {
		APIKey          string `mapstructure:"api_key"`
		TranscribeModel string `mapstructure:"transcribe_model"`
		ChatModel       string `mapstructure:"chat_model"`
	} `mapstructure:"openai"`
	FileMetadata struct {
		MaxFileSize int64 `mapstructure:"max_file_size"`
	} `mapstructure:"filemetadata"`
	Sink struct {
		WebhookURL     string `mapstructure:"webhook_url"`
		TimeoutSeconds int64  `mapstructure:"timeout_seconds"`
		FailRequest    bool   `mapstructure:"fail_request"`
	} `mapstructure:"sink"`
}

func Load() *Config {
	return &Config{
		Server: struct {
			Port string `mapstructure:"port"`
		}{Port: getEnv("PORT", "8080")},
		OpenAI: struct {
			APIKey          string `mapstructure:"api_key"`
			TranscribeModel string `mapstructure:"transcribe_model"`
			ChatModel       string `mapstructure:"chat_model"`
		}{
			APIKey:          getEnv("OPENAI_API_KEY", ""),
			TranscribeModel: getEnv("OPENAI_TRANSCRIBE_MODEL", "gpt-4o-transcribe"),
			ChatModel:       getEnv("OPENAI_CHAT_MODEL", "gpt-4o-mini"),
		},
		FileMetadata: struct {
			MaxFileSize int64 `mapstructure:"max_file_size"`
		}{
			MaxFileSize: getEnvInt("MAX_FILE_SIZE", 15728640),
		},
		Sink: struct {
			WebhookURL     string `mapstructure:"webhook_url"`
			TimeoutSeconds int64  `mapstructure:"timeout_seconds"`
			FailRequest    bool   `mapstructure:"fail_request"`
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
