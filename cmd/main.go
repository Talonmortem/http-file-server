package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/database"
	"github.com/Talonmortem/http-file-server/internal/router"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Не удалось загрузить конфигурацию")
	}

	cfg.Server.Port, _ = strconv.Atoi(os.Getenv("PORT"))

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

	router := router.SetupRouter(cfg)

	// Запуск сервера

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	} else {
		log.Printf("✅Сервер запущен на %s:%d\n", cfg.Server.Host, cfg.Server.Port)
	}

}
