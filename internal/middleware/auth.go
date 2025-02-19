package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			return
		}

		claims, err := ValidateToken(tokenString, cfg.JWT.SecretKey)
		if err != nil {
			log.Printf("Ошибка валидации токена: %v", err)
			status := http.StatusUnauthorized
			errorMsg := "Недействительный токен"
			var jwtErr *jwt.ValidationError
			if errors.As(err, &jwtErr) {
				if jwtErr.Errors == jwt.ValidationErrorExpired {
					status = http.StatusUnauthorized
					errorMsg = "Токен просрочен"
				}
			}
			c.AbortWithStatusJSON(status, gin.H{"error": errorMsg})
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if token, err := c.Cookie("jwt_token"); err == nil {
		return token
	}

	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.HasPrefix(bearerToken, "Bearer ") {
		return bearerToken[7:]
	}

	return ""
}

func ValidateToken(tokenString, secretKey string) (*Claims, error) {
	log.Printf("Validating JWT: %v", tokenString)
	log.Printf("Используемый секретный ключ: %s", secretKey)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		log.Printf("Parsing JWT with method: %v", token.Header["alg"])
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		log.Printf("Error parsing JWT: %v", err)
		return nil, err
	}

	if !token.Valid {
		log.Println("Invalid JWT")
		return nil, errors.New("invalid JWT")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		log.Println("Invalid format of claims")
		return nil, errors.New("invalid format of claims")
	}

	log.Printf("Successfully validated JWT: %v", claims)
	return claims, nil
}

func GenerateToken(username string, cfg *config.Config) (string, error) {
	log.Printf("Generating JWT for user: %v", username)

	// Используйте time.Hour для конвертации часов в Duration
	expirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpiresIn) * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	log.Printf("Generated claims: %v", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("Generated JWT: %v", token)

	signedString, err := token.SignedString([]byte(cfg.JWT.SecretKey))
	if err != nil {
		log.Printf("Error signing JWT: %v", err)
		return "", err
	}

	log.Printf("Successfully generated JWT: %v", signedString)
	return signedString, nil
}
