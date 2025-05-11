package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	if err := runMigrations(DB); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("База данных успешно подключена.")
	return nil
}

func runMigrations(db *sql.DB) error {
	createTables := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS expressions (
		id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expression TEXT NOT NULL,
		result TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`
	if _, err := db.Exec(createTables); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	if !ColumnExists(db, "expressions", "status") {
		_, err := db.Exec(`ALTER TABLE expressions ADD COLUMN status TEXT DEFAULT 'pending'`)
		if err != nil {
			return fmt.Errorf("failed to add status column: %w", err)
		}
		log.Println("Столбец 'status' добавлен в таблицу expressions.")
	} else {
		log.Println("Столбец 'status' уже существует.")
	}

	log.Println("Database migrated successfully")
	return nil
}

func ColumnExists(db *sql.DB, tableName, columnName string) bool {
	query := fmt.Sprintf(`PRAGMA table_info(%s);`, tableName)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Ошибка при выполнении PRAGMA table_info: %v", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			log.Printf("Ошибка при чтении результата PRAGMA: %v", err)
			return false
		}
		if name == columnName {
			return true
		}
	}
	return false
}
