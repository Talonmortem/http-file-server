package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ListFilesHandler возвращает файлы с поддержкой поддиректорий
// Обновленная функция ListFilesHandler
func ListFilesHandler(uploadDir, publicDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")

		userDir := filepath.Join(uploadDir, username)
		// Создаем директорию, если она не существует
		if err := os.MkdirAll(userDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user directory"})
			return
		}

		// Чтение файлов с путями
		myFiles := listFilesRecursive(userDir, "")
		publicFiles := listFilesRecursive(publicDir, "")

		c.JSON(http.StatusOK, gin.H{
			"my_files":     myFiles,
			"public_files": publicFiles,
		})
	}
}

// Обновленная listFilesRecursive для относительных путей
func listFilesRecursive(root, prefix string) []string {
	files := []string{}
	entries, err := os.ReadDir(root)
	if err != nil {
		return files
	}

	for _, entry := range entries {
		relPath := filepath.Join(prefix, entry.Name())
		if entry.IsDir() {
			subFiles := listFilesRecursive(filepath.Join(root, entry.Name()), relPath)
			files = append(files, subFiles...)
		} else {
			files = append(files, relPath)
		}
	}
	return files
}
