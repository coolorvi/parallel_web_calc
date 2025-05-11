package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"

	handler "github.com/coolorvi/parallel_web_calc/internal/auth"
	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/parser"
)

type Task struct {
	ID            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

type Expression struct {
	ID     string       `json:"id"`
	Status string       `json:"status"`
	Result *float64     `json:"result"`
	Node   *parser.Node `json:"-"`
	Tasks  []Task       `json:"-"`
	UserID int          `json:"-"`
}

type Orchestrator struct {
	DB          *sql.DB
	expressions map[string]*Expression
	tasks       chan Task
	results     map[int]float64
	taskID      int
	mu          sync.Mutex
}

func NewOrchestrator(db *sql.DB) *Orchestrator {
	return &Orchestrator{
		DB:          db,
		expressions: make(map[string]*Expression),
		tasks:       make(chan Task, 100),
		results:     make(map[int]float64),
	}
}

func (o *Orchestrator) AddExpression(w http.ResponseWriter, r *http.Request) {
	userID := handler.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	node, err := parser.Parse(req.Expression)
	if err != nil {
		http.Error(w, "Invalid expression", http.StatusUnprocessableEntity)
		return
	}

	id := uuid.New().String()
	expr := &Expression{
		ID:     id,
		Status: "pending",
		Node:   node,
		UserID: userID,
	}
	o.mu.Lock()
	o.expressions[id] = expr
	o.mu.Unlock()

	_, err = o.DB.Exec(`INSERT INTO expressions (id, user_id, expression, status) VALUES (?, ?, ?, ?)`,
		id, userID, req.Expression, "pending")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	go o.processExpression(expr)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (o *Orchestrator) processExpression(expr *Expression) {
	result := o.evaluateNode(expr.Node, expr)

	o.mu.Lock()
	expr.Result = &result
	expr.Status = "completed"
	o.mu.Unlock()

	_, err := o.DB.Exec(`UPDATE expressions SET result = ?, status = ? WHERE id = ?`,
		result, "completed", expr.ID)
	if err != nil {
	}
}

func (o *Orchestrator) evaluateNode(node *parser.Node, expr *Expression) float64 {
	if node.Operator == "" {
		return node.Value
	}

	var leftVal, rightVal float64
	if node.Left.Operator != "" {
		leftVal = o.evaluateNode(node.Left, expr)
	} else {
		leftVal = node.Left.Value
	}
	if node.Right.Operator != "" {
		rightVal = o.evaluateNode(node.Right, expr)
	} else {
		rightVal = node.Right.Value
	}

	o.mu.Lock()
	o.taskID++
	task := Task{
		ID:            o.taskID,
		Arg1:          leftVal,
		Arg2:          rightVal,
		Operation:     node.Operator,
		OperationTime: o.getOperationTime(node.Operator),
	}
	expr.Tasks = append(expr.Tasks, task)
	o.mu.Unlock()

	o.tasks <- task

	for {
		o.mu.Lock()
		if result, ok := o.results[task.ID]; ok {
			delete(o.results, task.ID)
			o.mu.Unlock()
			return result
		}
		o.mu.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (o *Orchestrator) getOperationTime(op string) int {
	switch op {
	case "+":
		return getEnvInt("TIME_ADDITION_MS", 1000)
	case "-":
		return getEnvInt("TIME_SUBTRACTION_MS", 1000)
	case "*":
		return getEnvInt("TIME_MULTIPLICATIONS_MS", 1000)
	case "/":
		return getEnvInt("TIME_DIVISIONS_MS", 1000)
	default:
		return 1000
	}
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
