package middleware

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// AuthMiddleware - проверка авторизации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Декодируем base64
		payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
		parts := strings.SplitN(string(payload), ":", 2)
		if len(parts) != 2 {
			log.Printf("Error decoding base64: %v", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		username := parts[0]
		password := parts[1]

		// Проверяем пользователя в БД

		var hash string
		err = database.DB.QueryRow("SELECT password_hash FROM users WHERE username = ?", username).Scan(&hash)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("❌ Пользователь не найден:", username)
			} else {
				log.Println("Ошибка запроса:", err)
			}
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
			log.Println("❌ Неверный пароль для:", username)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Авторизация успешна
		log.Println("✅ Пользователь авторизован:", username)
		c.Set("username", username)
		c.Next()
	}
}
