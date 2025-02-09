package database

import (
	"database/sql"
	"log"

	"github.com/Talonmortem/http-file-server/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// Подключение к БД

func ConnectDB(cfg *config.Config) error {
	db, err := sql.Open("sqlite3", cfg.Database.Path)
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	log.Println("✅ Подключение к SQLite успешно!")
	return nil
}
