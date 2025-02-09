package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DeleteFilesHandler удаляет файлы
func DeleteFilesHandler(uploadDir, publicDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var request struct {
			Files []string `json:"files"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		for _, file := range request.Files {
			filePath := filepath.Join(uploadDir, username.(string), file)
			if err := os.Remove(filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete " + file})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Files deleted"})
	}
}
