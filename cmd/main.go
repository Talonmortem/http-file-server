package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/config/middleware"
	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/Talonmortem/http-file-server/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Не удалось загрузить конфигурацию")
	}

	if err := os.MkdirAll(cfg.Storage.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("Не удалось создать директорию для загрузок: %v", err)
	}

	// Подключаемся к базе данных
	err = database.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer database.DB.Close()

	database.RunMigrations()
	log.Println("✅ База данных подключена.")

	// Создаем экземпляр Gin
	router := gin.Default()

	// Включаем авторизацию
	router.Use(middleware.AuthMiddleware())

	// Главная стараница
	router.GET("/", handlers.IndexHandler())

	// register routes
	router.POST("/upload", handlers.UploadHandler(cfg.Storage.UploadDir))
	router.POST("/download-zip", handlers.DownloadFilesHandler(cfg.Storage.UploadDir))
	router.GET("/files", handlers.ListFilesHandler(cfg.Storage.UploadDir, cfg.Storage.PublicDir))
	router.POST("/delete", handlers.DeleteFilesHandler(cfg.Storage.UploadDir, cfg.Storage.PublicDir))
	router.POST("/download/:filename", handlers.DownloadOnClickHandler(cfg.Storage.UploadDir))

	// Запуск сервера

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
	log.Printf("Сервер запущен на %s:%d\n", cfg.Server.Host, cfg.Server.Port)

}
