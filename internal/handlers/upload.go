package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/gin-gonic/gin"
)

// UploadHandler обрабатывает загрузку файлов
func UploadHandler(rootDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем имя пользователя из контекста
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		path := c.Query("path")
		// Читаем путь назначения из параметров запроса
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		files := form.File["files"]

		log.Printf("path: %s \n Files: %v \n", path, files)

		for _, file := range files {
			// Создаем путь для сохраненея
			filePath := filepath.Join(path, file.Filename)

			// Сохраняем файл
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
			database.SaveUploadedFile(filePath, username.(string))

			log.Printf("User %s uploaded file %s:", username, filePath)
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d файлов загружено.", len(files))})
	}
}
