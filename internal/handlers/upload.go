package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadHandler обрабатывает загрузку файлов
func UploadHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Читаем путь назначения из параметров запроса
		targetPath := c.PostForm("path") // Например, "docs/"
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		files := form.File["files"]

		log.Printf("targetPath: %s \n Files: %v \n", targetPath, files)

		for _, file := range files {
			// Создаем путь для сохраненея
			filePath := filepath.Join(uploadDir, username.(string), targetPath, file.Filename)

			// Сохраняем файл
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d файлов загружено.", len(files))})
	}

	/*
		// Получаем загружаемый файл
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
			return
		}
		defer file.Close()

		// Создаём путь назначения
		userDir := filepath.Join(uploadDir, username.(string), targetPath)
		os.MkdirAll(userDir, 0755) // Создаём папку, если её нет

		filePath := filepath.Join(userDir, header.Filename)

		// Сохраняем файл
		outFile, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
			return
		}
		defer outFile.Close()

		_, err = outFile.ReadFrom(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write file"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded", "path": targetPath})
	} */
}
