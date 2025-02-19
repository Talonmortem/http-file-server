package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DeleteFilesHandler удаляет файлы
func DeleteFilesHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем имя пользователя из контекста
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
			log.Printf("Invalid request: %v", err)
			return
		}

		for _, file := range request.Files {
			log.Printf("User %s deleting file/directory %s", username, file)
			//check if file is directory
			fileInfo, err := os.Stat(filepath.Join(file))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
				log.Printf("Failed to get file %s info: %v", file, err)
				return
			}
			if fileInfo.IsDir() {
				//delete directory
				if err := os.RemoveAll(filepath.Join(file)); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete directory"})
					log.Printf("Failed to delete directory %s: %v", file, err)
					return
				}
				log.Printf("User %s deleted directory %s", username, file)
				continue
			}

			filePath := filepath.Join(file)
			if err := os.Remove(filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete " + file})
				log.Printf("Failed to delete file %s: %v", file, err)
				return
			}
			log.Printf("User %s deleted file %s", username, file)
		}

		c.JSON(http.StatusOK, gin.H{"message": "Files deleted"})
	}
}
