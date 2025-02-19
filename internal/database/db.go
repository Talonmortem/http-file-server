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

func SaveNotes(request struct {
	Path string `json:"path"`
	Note string `json:"note"`
}) error {

	// insert into database replace on conflict
	_, err := DB.Exec("INSERT INTO files (path, notes) VALUES ($1, $2) ON CONFLICT (path) DO UPDATE SET notes = $2", request.Path, request.Note)
	if err != nil {
		return err
	}
	return nil
}

func GetNotes(path string) string {
	var notes string
	err := DB.QueryRow("SELECT notes FROM files WHERE path = $1", path).Scan(&notes)
	if err != nil {
		return ""
	}
	return notes
}

func SaveUploadedFile(path string, username string) error {
	_, err := DB.Exec("INSERT INTO files (path, owner) VALUES ($1, $2) ON CONFLICT (path) DO UPDATE SET owner = $2", path, username)
	if err != nil {
		return err
	}
	return nil
}

func GetOwner(path string) string {
	var owner string
	err := DB.QueryRow("SELECT owner FROM files WHERE path = $1", path).Scan(&owner)
	if err != nil {
		return ""
	}
	return owner
}
