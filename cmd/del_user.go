package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/database"
)

//TODO: del user

func main() {
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
		fmt.Println("Использование: go run cmd/del_user.go <username>")
		return
	}
	username := os.Args[1]

	// Удаляем пользователя в БД
	_, err = database.DB.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		log.Fatal("Ошибка удаления пользователя:", err)
	}

	fmt.Println("✅ Пользователь удалён:", username)
}
