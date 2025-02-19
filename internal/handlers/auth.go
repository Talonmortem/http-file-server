package handlers

import (
	"log"
	"net/http"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/Talonmortem/http-file-server/internal/middleware"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// handlers/auth.go
func LoginHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем данные из формы
		username := c.PostForm("username")
		password := c.PostForm("password")

		log.Printf("Попытка входа пользователя: %s", username)

		// Проверка учетных данных
		var hash string
		err := database.DB.QueryRow(
			"SELECT password_hash FROM users WHERE username = ?",
			username,
		).Scan(&hash)

		if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
			// Перенаправляем с ошибкой
			c.Redirect(http.StatusFound, "/login?error=invalid_credentials")
			return
		}

		// Генерация токена
		token, err := middleware.GenerateToken(username, cfg)
		if err != nil {
			c.Redirect(http.StatusFound, "/login?error=internal_error")
			return
		}

		// Сохраняем токен в куки
		c.SetCookie(
			"jwt_token",
			token,
			cfg.JWT.ExpiresIn*3600, // В секундах
			"/",
			"",
			false, // Только HTTPS в продакшене
			true,  // HttpOnly
		)

		log.Printf("Токен для пользователя %s: %s", username, token)

		c.Redirect(http.StatusFound, "/")
	}
}

func LogoutHandler(c *gin.Context) {
	// Удаляем куку
	c.SetCookie(
		"jwt_token",
		"",
		-1, // Удалить куку
		"/",
		"",
		false,
		true,
	)

	c.Redirect(http.StatusFound, "/login")
}
