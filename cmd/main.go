package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/coolorvi/parallel_web_calc/internal/agent"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator"
	"github.com/coolorvi/parallel_web_calc/internal/storage"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbPath := "./db.sqlite"

	applyMigrations(dbPath, "internal/migrations/init.sql")

	_, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	storage.InitDB(dbPath)

	go agent.StartWorker()
	orchestrator.Start()
}

func applyMigrations(dbPath, sqlFile string) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}
	defer db.Close()

	data, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}

	statements := strings.Split(string(data), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, err := db.Exec(stmt)
		if err != nil {
			log.Fatalf("Migration failed: %v\nQuery: %s", err, stmt)
		}
	}

	fmt.Println("Database migrated successfully")
}
