package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// DownloadFilesHandler создаёт ZIP-архив с выбранными файлами
func DownloadFilesHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var request struct {
			Files          []string `json:"files"`
			CurrentDirName string   `json:"currentDirName"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if len(request.Files) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No files selected"})
			return
		}

		// Создаём временный файл для ZIP
		zipFilePath := filepath.Join(os.TempDir(), fmt.Sprintf("files_%s.zip", username))
		outFile, err := os.Create(zipFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create zip file"})
			return
		}
		defer outFile.Close()

		zipWriter := zip.NewWriter(outFile)

		// Добавляем файлы в ZIP
		for _, file := range request.Files {
			// Проверяем, что путь безопасен
			//filePath := filepath.Join(uploadDir, username.(string), file)
			filePath := file
			// Убеждаемся, что файл существует
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("File %s not found", file)})
				log.Printf("File %s not found", file)
				return
			}

			fileToZip, err := os.Open(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Cannot open file %s", file)})
				return
			}
			defer fileToZip.Close()

			// Создаём запись в ZIP с базовым именем файла
			w, err := zipWriter.Create(filepath.Base(file))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Cannot add file %s to zip", file)})
				return
			}

			_, err = io.Copy(w, fileToZip)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error copying file %s", file)})
				return
			}
		}

		// Закрываем ZIP перед отправкой
		zipWriter.Close()

		// Отправляем файл клиенту
		currentDirName := filepath.Base(request.CurrentDirName)
		c.FileAttachment(zipFilePath, currentDirName+"_"+username.(string)+".zip")
		log.Printf("✅ Архив %s отправлены пользователем %s", currentDirName+"_"+username.(string)+".zip", username.(string))
	}
}

// DownloadOnClickHandler

func DownloadOnClickHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := c.Param("filename")
		// Получаем имя пользователя из контекста
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		filePath := filepath.Join(uploadDir, username.(string), fileName)
		log.Printf("Downloading file: %s", filePath)
		c.File(filePath)
	}
}
