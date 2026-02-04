# Voiceline - AI-Powered Conversational Voice Assistant

### Backend Service Details

#### **Core Features**

1. **Audio Upload Endpoint** (`POST /audio`)

   - Accepts audio files via multipart/form-data
   - Stores files temporarily (local temp dir)
   - Processes with Whisper API for transcription
   - Sends transcript to LLM for response generation
   - Sends webhook payload after successful processing (if configured)

2. **Health Check Endpoint** (`GET /health`)

   - Monitors service availability
   - Returns API version and status

3. **Response Format**

   ```json
   {
     "transcript": "What time is it?",
     "answer": "It's currently 3:45 PM",
     "prompt": "Provide a helpful response to this message."
   }
   ```

4. **Webhook Integration**

   After processing, the service sends a POST request to the configured webhook URL with:

   ```json
   {
     "transcript": "What time is it?",
     "answer": "It's currently 3:45 PM",
     "prompt": "Provide a helpful response to this message.",
     "created_at": "2026-02-03T11:30:00Z"
   }
   ```

## Setup Instructions

### Prerequisites

- Go 1.24+ installed
- OpenAI API key

### Quick Start

1. **Clone the repository**

   ```bash
   git clone https://github.com/mirkosisko-dev/voiceline.git
   cd voiceline
   ```

2. **Create environment file**

   ```bash
   cp .env.example .env
   ```

3. **Configure your OpenAI API key**

   ```bash
   # Open .env and set your API key
   OPENAI_API_KEY=your-api-key-here
   OPENAI_TRANSCRIBE_MODEL=gpt-4o-transcribe
   OPENAI_CHAT_MODEL=gpt-4o-mini

   PORT=8080

   MAX_FILE_SIZE=15728640 # 15MB

   # Optional webhook configuration
   SINK_WEBHOOK_URL=https://webhook.site/your-unique-url
   SINK_WEBHOOK_TIMEOUT_SECONDS=10
   SINK_FAIL_REQUEST=false
   ```

4. **Run the server**

   ```bash
   go run cmd/server/main.go
   ```

5. **Test the API**

   ```bash
   # Health check
   curl http://localhost:8080/health

   # Upload audio file
   curl -X POST http://localhost:8080/audio \
     -F "file=@path/to/your/audio.wav"
   ```

## API Documentation

### Endpoints

#### GET /health

Check service health and status.

**Response:**

```json
{
  "status": "healthy",
  "service": "voiceline-api"
}
```

#### POST /audio

Upload audio file for processing.

**Request:**

- Content-Type: multipart/form-data
- Form field: `file` (audio file), `prompt` (instructions for the LLM)

**Success Response:**

```json
{
  "transcript": "What time is it?",
  "answer": "It's currently 3:45 PM"
}
```

## Project Structure

```
voiceline/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handler/
│   │   ├── audio.go             # Audio upload handler
│   │   └── health.go            # Health check handler
│   ├── service/
│   │   ├── openai.go            # OpenAI API integration
│   │   ├── audio.go             # Audio processing orchestration
│   │   └── sink.go              # Webhook notification service
│   └── config/
│       └── config.go            # Configuration loading
├── pkg/
│   └── response/
│       └── response.go          # API response structures
├── .env.example                 # Environment template
├── go.mod                       # Go module dependencies
└── go.sum                       # Dependency checksums
```

## Configuration

### Environment Variables

| Variable                          | Description                     | Default             |
| --------------------------------- | ------------------------------- | ------------------- |
| `OPENAI_API_KEY`                  | OpenAI API key (required)       | -                   |
| `OPENAI_TRANSCRIBE_MODEL`         | Whisper model for transcription | `gpt-4o-transcribe` |
| `OPENAI_CHAT_MODEL`               | GPT model for conversation      | `gpt-4o-mini`       |
| `PORT`                            | Server port                     | `8080`              |
| `MAX_FILE_SIZE`                   | Max file size allowed           | `15728640`          |
| `SINK_WEBHOOK_URL`                | Webhook URL for notifications   | -                   |
| `SINK_WEBHOOK_TIMEOUT_SECONDS=10` | Webhook timeout in seconds      | `10`                |
| `SINK_FAIL_REQUEST`               | Fail request on webhook failure | `false`             |

## Technical Architecture

### Complete VoiceLine Tech Stack (Theoretical)

#### **Mobile Application (React Native/Expo)**

- **Frontend Framework**: React Native with Expo
- **UI Components**: NativeWind for modern, performant UI
- **State Management**: Zustand for global state + React Query for server state
- **Audio Capture**: Expo Audio for high-quality voice recording
- **Push Notifications**: Expo Push Notifications SDK
- **Offline Support**: AsyncStorage

#### **Backend Service (Go + Gin)**

- **Language**: Go 1.24+
- **Framework**: Gin for HTTP routing and middleware
- **API Documentation**: Swagger/OpenAPI

#### **Infrastructure & DevOps**

- **Containerization**: Docker + Docker Compose
- **Orchestration**: Kubernetes
- **Database**: PostgreSQL (user data, conversation history, transcripts)
- **Message Queue**: Redis for async job processing
- **CDN**: Cloudflare for static assets and DDoS protection
- **Monitoring**: Prometheus + Grafana
- **CI/CD**: GitHub Actions with automated testing and deployment

#### **AI Services**

- **Audio Transcription**: OpenAI Whisper API (state-of-the-art speech-to-text)
- **LLM Provider**: OpenAI GPT-4o-mini or Anthropic Claude 3.5 Sonnet
- **Rate Limiting**: Per-user token management to prevent abuse
