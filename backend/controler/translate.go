package controler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type TranslateRequest struct {
	Text string `json:"text" binding:"required"`
}
type TranslateItem struct {
	SourceText string `json:"source_text"`
	TargetText string `json:"target_text"`
}

func (AppControler) Translate(c *gin.Context) {
	var req TranslateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text cannot be empty"})
		return
	}
	
}
