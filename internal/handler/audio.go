package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mirkosisko-dev/voiceline/internal/config"
	"github.com/mirkosisko-dev/voiceline/internal/service"
	"github.com/mirkosisko-dev/voiceline/pkg/response"
)

type AudioHandler struct {
	audioService *service.AudioService
	cfg          *config.Config
}

func NewAudioHandler(audioService *service.AudioService, cfg *config.Config) *AudioHandler {
	return &AudioHandler{
		audioService: audioService,
		cfg:          cfg,
	}
}
func (h *AudioHandler) HandleAudioUpload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "missing form field 'file'", "")
		return
	}

	ct := fileHeader.Header.Get("Content-Type")
	if err := validateAudioUpload(fileHeader.Filename, ct); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid audio file", err.Error())
		return
	}

	if h.cfg.FileMetadata.MaxFileSize > 0 {
		if fileHeader.Size > h.cfg.FileMetadata.MaxFileSize {
			maxMB := h.cfg.FileMetadata.MaxFileSize / (1024 * 1024)
			response.Error(c, http.StatusRequestEntityTooLarge, "file size exceeds limit", fmt.Sprintf("max file size is %dMB (you uploaded %d MB)", maxMB, fileHeader.Size/(1024*1024)))
			return
		}
	}

	prompt := c.PostForm("prompt")
	if prompt == "" {
		prompt = "Provide a helpful response to this message."
	}

	tmpDir, err := os.MkdirTemp("", "mvp-audio-*")
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create temp dir", "")
		return
	}

	defer os.RemoveAll(tmpDir)

	inPath := filepath.Join(tmpDir, filepath.Base(fileHeader.Filename))

	if err := c.SaveUploadedFile(fileHeader, inPath); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to save upload", "")
		return
	}

	transcript, err := h.audioService.ProcessAudio(c.Request.Context(), inPath, prompt)
	if err != nil {
		response.Error(c, http.StatusBadGateway, "processing failed", err.Error())
		return
	}

	response.Success(c, transcript.Transcript, transcript.Answer, transcript.Prompt)
}
