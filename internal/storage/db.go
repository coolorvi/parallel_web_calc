package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(path string) {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Ошибка при открытии БД: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}

	log.Println("База данных успешно подключена.")
}
