package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/coolorvi/parallel_web_calc/internal/auth"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/parser"
	"github.com/golang-jwt/jwt/v5"
)

func TestAddExpression_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	expression := "2+3"
	mock.ExpectExec("INSERT INTO expressions").
		WithArgs(sqlmock.AnyArg(), 1, expression, "pending").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("UPDATE expressions SET result").WillReturnResult(sqlmock.NewResult(1, 1))

	orch := NewOrchestrator(db)

	body := map[string]string{"expression": expression}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/expressions", bytes.NewReader(b))
	req = SetUserID(req, 1)
	rec := httptest.NewRecorder()

	go func() {
		task := <-orch.tasks
		orch.mu.Lock()
		orch.results[task.ID] = task.Arg1 + task.Arg2
		orch.mu.Unlock()
	}()

	orch.AddExpression(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	if resp["id"] == "" {
		t.Fatal("expected expression ID in response")
	}
}

func TestAddExpression_InvalidJSON(t *testing.T) {
	orch := NewOrchestrator(&sql.DB{})

	req := httptest.NewRequest(http.MethodPost, "/expressions", bytes.NewReader([]byte(`invalid`)))
	req = SetUserID(req, 1)
	rec := httptest.NewRecorder()

	orch.AddExpression(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", rec.Code)
	}
}

func TestAddExpression_Unauthorized(t *testing.T) {
	orch := NewOrchestrator(&sql.DB{})
	req := httptest.NewRequest(http.MethodPost, "/expressions", nil)
	rec := httptest.NewRecorder()

	orch.AddExpression(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestEvaluateNode_Recursive(t *testing.T) {
	orch := NewOrchestrator(&sql.DB{})
	node, _ := parser.Parse("2+3")
	expr := &Expression{}
	go func() {
		task := <-orch.tasks
		orch.mu.Lock()
		orch.results[task.ID] = task.Arg1 + task.Arg2
		orch.mu.Unlock()
	}()
	result := orch.evaluateNode(node, expr)

	if result != 5 {
		t.Fatalf("expected result 5, got %.2f", result)
	}
}

func TestGetOperationTime_EnvOverride(t *testing.T) {
	os.Setenv("TIME_ADDITION_MS", "2500")
	defer os.Unsetenv("TIME_ADDITION_MS")

	orch := NewOrchestrator(&sql.DB{})
	got := orch.getOperationTime("+")

	if got != 2500 {
		t.Fatalf("expected 2500 from env, got %d", got)
	}
}

func SetUserID(r *http.Request, userID int) *http.Request {
	claims := &auth.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(auth.JwtKey)

	r.Header.Set("Authorization", "Bearer "+tokenStr)
	return r
}
