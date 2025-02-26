package handlers

import (
	"net/http"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/gin-gonic/gin"
)

func ConfigHandler(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		rootPath := cfg.Storage.UploadDir[2:]

		c.JSON(http.StatusOK, gin.H{
			"rootPath": rootPath,
			"webDir":   cfg.Storage.WebDir,
		})
	}
}
