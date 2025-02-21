package handlers

import (
	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/gin-gonic/gin"
)

// IndexHandler раздаёт главную страницу
func IndexHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.File(cfg.Storage.WebDir + "/index.html")
	}
}
