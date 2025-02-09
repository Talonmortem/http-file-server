package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// CreateDirHandler - обработчик создания директорий
func CreateDirHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var request struct {
			DirName string `json:"dir_name"`
			Path    string `json:"path"`
		}

		request.Path = filepath.Join(uploadDir, username.(string))

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Формируем полный путь
		userDir := filepath.Join(uploadDir, username.(string))
		targetPath := filepath.Join(userDir, request.Path, request.DirName)

		if len(targetPath) < len(userDir) || targetPath[:len(userDir)] != userDir {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid directory path"})
			return
		}

		// Создаем директорию
		if err := os.Mkdir(targetPath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Directory created"})
	}
}
