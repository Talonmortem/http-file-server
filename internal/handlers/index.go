package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

// IndexHandler раздаёт главную страницу
func IndexHandler(webDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Главная страница!", webDir)
		c.File(webDir + "/index.html")
	}
}
