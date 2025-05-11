package tests

import (
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/storage"
)

func TestInitDB_CreatesTablesAndStatusColumn(t *testing.T) {
	err := storage.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	tables := []string{"users", "expressions"}
	for _, table := range tables {
		exists := tableExists(table)
		if !exists {
			t.Errorf("expected table %q to exist", table)
		}
	}

	if !storage.ColumnExists(storage.DB, "expressions", "status") {
		t.Errorf("expected column 'status' to exist in 'expressions' table")
	}
}

func tableExists(name string) bool {
	row := storage.DB.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name=?`, name)
	var tableName string
	err := row.Scan(&tableName)
	return err == nil && tableName == name
}

func TestColumnExists_ReturnsTrueForExistingColumn(t *testing.T) {
	err := storage.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	if !storage.ColumnExists(storage.DB, "users", "login") {
		t.Errorf("expected column 'login' to exist in 'users'")
	}
}

func TestColumnExists_ReturnsFalseForMissingColumn(t *testing.T) {
	err := storage.InitDB(":memory:")
	if err != nil {
		t.Fatalf("InitDB failed: %v", err)
	}

	if storage.ColumnExists(storage.DB, "users", "nonexistent_column") {
		t.Errorf("did not expect 'nonexistent_column' to exist in 'users'")
	}
}
