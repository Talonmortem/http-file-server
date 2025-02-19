package database

import "log"

// RunMigrations создаёт таблицы, если их нет
func RunMigrations() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Ошибка миграции БД:", err)
	}

	/*
		_, err = DB.Exec(`
		DROP TABLE IF EXISTS files;
		`)
		if err != nil {
			log.Fatal("Ошибка миграции БД:", err)
		}
	*/

	query = `CREATE TABLE IF NOT EXISTS files (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		path TEXT,
		notes TEXT,
		owner TEXT,
		UNIQUE(path)
	);`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatal("Ошибка миграции БД:", err)
	}

	log.Println("✅ Таблицы проверены/созданы.")
}
