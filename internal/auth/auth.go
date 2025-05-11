package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET"))

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid input")
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Server error")
			return
		}

		_, err = db.Exec("INSERT INTO users (login, password_hash) VALUES (?, ?)", creds.Login, string(hash))
		if err != nil {
			writeJSONError(w, http.StatusConflict, "User already exists")
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
	}
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			writeJSONError(w, http.StatusBadRequest, "Invalid input")
			return
		}

		var id int
		var hashed string
		err := db.QueryRow("SELECT id, password_hash FROM users WHERE login = ?", creds.Login).Scan(&id, &hashed)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "Invalid login or password")
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(creds.Password)) != nil {
			writeJSONError(w, http.StatusUnauthorized, "Invalid login or password")
			return
		}

		exp := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			UserID: id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(exp),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(JwtKey)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "Could not create token")
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
	}
}

func GetUserID(r *http.Request) int {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0
	}

	tokenStr := authHeader[len("Bearer "):]
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0
	}
	return claims.UserID
}
