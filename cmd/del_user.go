package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Talonmortem/http-file-server/internal/config"
	"github.com/Talonmortem/http-file-server/internal/database"

	"golang.org/x/crypto/bcrypt"
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
	if len(os.Args) < 3 {
		fmt.Println("Использование: go run cmd/add_user.go <username> <password>")
		return
	}
	username, password := os.Args[1], os.Args[2]

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Ошибка хеширования пароля:", err)
	}

	// Добавляем пользователя в БД
	_, err = database.DB.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", username, string(hash))
	if err != nil {
		log.Fatal("Ошибка добавления пользователя:", err)
	}

	fmt.Println("✅ Пользователь добавлен:", username)
}
