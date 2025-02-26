package router

import (
	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/handlers"
	"github.com/Talonmortem/http-file-server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	// Создаем экземпляр Gin
	router := gin.Default()

	//Загрузка html-шаблонов
	router.LoadHTMLGlob(cfg.Storage.TemplateDir)

	//Публичные роуты
	router.GET("/login", handlers.LoginFormHandler)
	router.POST("/login", handlers.LoginHandler(cfg))
	router.GET("/logout", handlers.LogoutHandler)

	//Защищенные роуты
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware(cfg))
	{
		// Главная стараница
		authGroup.GET("/", handlers.IndexHandler(cfg))
		authGroup.GET("/config", handlers.ConfigHandler(cfg))
		authGroup.GET("/files", handlers.ListFilesHandler(cfg.Storage.UploadDir))

		// register routes
		authGroup.POST("/upload", handlers.UploadHandler(cfg.Storage.UploadDir))
		authGroup.POST("/download-zip", handlers.DownloadFilesHandler(cfg.Storage.UploadDir))
		authGroup.POST("/delete", handlers.DeleteFilesHandler(cfg.Storage.UploadDir))
		authGroup.GET("/download/:filename", handlers.DownloadOnClickHandler(cfg.Storage.UploadDir))
		authGroup.POST("/save-note", handlers.SaveNoteHandler(cfg.Storage.UploadDir))
		authGroup.POST("/create-folder", handlers.CreateFolderHandler(cfg.Storage.UploadDir))
		authGroup.POST("/move-file", handlers.MoveHandler(cfg.Storage.UploadDir))
	}

	return router
}
