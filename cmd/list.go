package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/database"
)

//TODO: del user

func listUsers() {
	// Загружаем конфиг
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Подключаемся к БД
	err = database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer database.DB.Close()

	// Запрашиваем логин и пароль
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run cmd/list_users.go")
		return
	}

	// Удаляем пользователя в БД
	users, err := database.DB.Exec("SELECT * FROM users")
	if err != nil {
		log.Fatal("Ошибка поиска пользователя:", err)
	}

	fmt.Println("✅ Пользователи:", users)
}
