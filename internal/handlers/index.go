package handlers

import (
	"github.com/gin-gonic/gin"
)

// IndexHandler раздаёт главную страницу
func IndexHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File("web/index.html")
	}
}
