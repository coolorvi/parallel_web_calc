package tests

import (
	"testing"

	"github.com/coolorvi/parallel_web_calc/internal/orchestrator/parser"
)

func TestParseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `{"expression": "3+5*2-8/4"}`,
			expected: `{"sub_expressions":[{"operator":"+","left_operand":"3","right_operand":"5"},{"operator":"*","left_operand":"5","right_operand":"2"},{"operator":"-","left_operand":"2","right_operand":"8"},{"operator":"/","left_operand":"8","right_operand":"4"}]}`,
		},
		{
			input:    `{"expression": "10-2*3"}`,
			expected: `{"sub_expressions":[{"operator":"-","left_operand":"10","right_operand":"2"},{"operator":"*","left_operand":"2","right_operand":"3"}]}`,
		},
		{
			input:    `{"expression": "7+8"}`,
			expected: `{"sub_expressions":[{"operator":"+","left_operand":"7","right_operand":"8"}]}`,
		},
	}

	for _, tt := range tests {
		output, err := parser.ParseExpression([]byte(tt.input))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(output) != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, string(output))
		}
	}
}
