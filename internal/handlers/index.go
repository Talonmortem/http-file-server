package handlers

import (
	"github.com/gin-gonic/gin"
)

// IndexHandler раздаёт главную страницу
func IndexHandler(webDir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(webDir + "/index.html")
	}
}
