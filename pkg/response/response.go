package response

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
type SuccessResponse struct {
	Transcript string `json:"transcript"`
	Answer     string `json:"answer"`
	Prompt     string `json:"prompt,omitempty"`
}

func Error(c *gin.Context, statusCode int, error, details string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   error,
		Details: details,
	})
}
func Success(c *gin.Context, transcript, answer string, prompt string) {
	c.JSON(200, SuccessResponse{
		Transcript: transcript,
		Answer:     answer,
		Prompt:     prompt,
	})
}
