package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func MoveHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		request := struct {
			File        string `json:"file"`
			Destination string `json:"destination"`
		}{}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			log.Printf("Invalid request: %v", err)
			return
		}

		filePath := filepath.Join(request.Destination, request.File)

		if _, err := os.Stat(filePath); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File already exists"})
			log.Printf("File %s already exists", filePath)
			return
		}

		filePath = filepath.Join(request.Destination, filepath.Base(filePath))

		if err := os.Rename(request.File, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to move file"})
			log.Printf("Failed to move file %s to %s: %v", request.File, filePath, err)
			return
		}

		log.Printf("User %s moving file %s to %s", username, request.File, request.Destination)
		c.JSON(http.StatusOK, gin.H{"message": "File moved successfully"})

	}
}
