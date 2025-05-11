package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/coolorvi/parallel_web_calc/internal/auth"
)

var db *sql.DB

func TestRegister(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, login TEXT, password_hash TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	handler := auth.Register(db)
	requestBody := `{"login": "testuser", "password": "testpassword"}`
	req := httptest.NewRequest("POST", "/register", bytes.NewReader([]byte(requestBody)))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code, "Expected status OK")
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}

func TestLogin(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	auth.JwtKey = []byte(os.Getenv("JWT_SECRET"))

	db := setupTestDB(t)

	creds := auth.Credentials{
		Login:    "user1",
		Password: "pass123",
	}
	body, _ := json.Marshal(creds)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	handler := auth.Register(db)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rec = httptest.NewRecorder()
	handler = auth.Login(db)
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var respBody map[string]string
	err := json.NewDecoder(rec.Body).Decode(&respBody)
	require.NoError(t, err)

	tokenStr, ok := respBody["token"]
	require.True(t, ok, "token not found in response")

	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return auth.JwtKey, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
	require.Greater(t, claims.UserID, 0)
	require.WithinDuration(t, time.Now().Add(24*time.Hour), claims.ExpiresAt.Time, time.Hour)
}

func TestGetUserID(t *testing.T) {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, login TEXT, password_hash TEXT)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	_, err = db.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", "testuser", "$2a$10$B7YsHyR8kA.qkGhpD/ZYp.DdTfmIQ0iIs23dsyKq1d8qlHqfzme7a")
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	claims := &auth.Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(auth.JwtKey)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	req := httptest.NewRequest("GET", "/some_endpoint", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)

	userID := auth.GetUserID(req)
	assert.Equal(t, 1, userID, "Expected userID to be 1")
}
