package handlers

import (
	"log"
	"net/http"

	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/gin-gonic/gin"
)

func SaveNoteHandler(uploadDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем имя пользователя из контекста
		username, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var request struct {
			Path string `json:"path"`
			Note string `json:"note"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := database.SaveNotes(request); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save note"})
			log.Println("Failed to save note:", err)
			return
		}
		log.Printf("✅ Заметка для %s сохранена пользователем %s: %s", request.Path, username.(string), request.Note)
		c.JSON(http.StatusOK, gin.H{"message": "Note saved successfully"})
	}
}
