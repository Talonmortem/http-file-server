package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/gin-gonic/gin"
)

func CreateFolderHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, _ := c.Get("username")
		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Read dir from params
		path := c.Query("path")

		dirName := filepath.Join(path, "New Folder")
		if _, err := os.Stat(dirName); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Directory already exists"})
			log.Printf("Directory %s already exists", dirName)
			return
		}

		// Создаем директорию
		if err := os.Mkdir(dirName, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory"})
			log.Printf("Failed to create directory %s: %v", path, err)
			return
		}

		// Сохраняем имя пользователя в базу данных
		if err := database.SaveUploadedFile(dirName, username.(string)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
			log.Printf("Failed to save info about %s: %v", dirName, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Directory created"})
		log.Printf("✅ Директория %s создана пользователем %s", dirName, username.(string))
	}
}
