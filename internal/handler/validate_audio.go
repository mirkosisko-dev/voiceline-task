package handler

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"
)

var allowedExt = map[string]bool{
	".wav":  true,
	".mp3":  true,
	".m4a":  true,
	".webm": true,
	".ogg":  true,
}

func validateAudioUpload(filename string, contentType string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedExt[ext] {
		return fmt.Errorf("unsupported file extension: %s", ext)
	}

	if contentType != "" {
		mt, _, _ := mime.ParseMediaType(contentType)
		if mt != "" && !strings.HasPrefix(mt, "audio/") && mt != "application/octet-stream" {
			return fmt.Errorf("unsupported content-type: %s", mt)
		}
	}

	return nil
}
