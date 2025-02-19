package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginFormHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"error": c.Query("error"), // Для вывода ошибок
	})
}
