package tests

import (
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/parser"
)

func TestParseSimpleAddition(t *testing.T) {
	node, err := parser.Parse("2+3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Operator != "+" {
		t.Errorf("expected operator '+', got %v", node.Operator)
	}
	if node.Left.Value != 2 || node.Right.Value != 3 {
		t.Errorf("expected 2 + 3, got %v + %v", node.Left.Value, node.Right.Value)
	}
}

func TestParseWithPrecedence(t *testing.T) {
	node, err := parser.Parse("2+3*4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Operator != "+" {
		t.Errorf("expected root operator '+', got %v", node.Operator)
	}
	if node.Left.Value != 2 {
		t.Errorf("expected left to be 2, got %v", node.Left.Value)
	}
	if node.Right.Operator != "*" {
		t.Errorf("expected right operator '*', got %v", node.Right.Operator)
	}
}

func TestParseWithParentheses(t *testing.T) {
	node, err := parser.Parse("(2+3)*4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Operator != "*" {
		t.Errorf("expected root operator '*', got %v", node.Operator)
	}
	if node.Left.Operator != "+" {
		t.Errorf("expected left operator '+', got %v", node.Left.Operator)
	}
	if node.Right.Value != 4 {
		t.Errorf("expected right value 4, got %v", node.Right.Value)
	}
}

func TestParseInvalidExpression(t *testing.T) {
	_, err := parser.Parse("2+*3")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestParseEmptyExpression(t *testing.T) {
	_, err := parser.Parse("")
	if err == nil {
		t.Error("expected error for empty expression, got nil")
	}
}

func TestParseFloatNumbers(t *testing.T) {
	node, err := parser.Parse("1.5*2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if node.Operator != "*" {
		t.Errorf("expected '*', got %v", node.Operator)
	}
	if node.Left.Value != 1.5 || node.Right.Value != 2 {
		t.Errorf("expected 1.5 * 2, got %v * %v", node.Left.Value, node.Right.Value)
	}
}
