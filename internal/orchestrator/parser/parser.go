package parser

import (
	"encoding/json"
	"strings"
	"unicode"
)

type Expression struct {
	Operator     string `json:"operator"`
	LeftOperand  string `json:"left_operand"`
	RightOperand string `json:"right_operand"`
}

type InputJSON struct {
	Expression string `json:"expression"`
}

type OutputJSON struct {
	SubExpressions []Expression `json:"sub_expressions"`
}

func ParseExpression(input []byte) ([]byte, error) {
	var inputJSON InputJSON
	if err := json.Unmarshal(input, &inputJSON); err != nil {
		return nil, err
	}

	tokens := tokenize(inputJSON.Expression)
	subExpressions := shuntingYard(tokens)

	output := OutputJSON{SubExpressions: subExpressions}
	return json.Marshal(output)
}

func tokenize(expr string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, c := range expr {
		if unicode.IsDigit(c) || c == '.' {
			currentToken.WriteRune(c)
		} else if strings.ContainsRune("+-*/()", c) {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(c))
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}

	return tokens
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func shuntingYard(tokens []string) []Expression {
	var output []Expression
	var stack []string
	var operandStack []string

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				right := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]
				left := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]
				op := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				expr := Expression{Operator: op, LeftOperand: left, RightOperand: right}
				output = append(output, expr)
				operandStack = append(operandStack, "res")
			}
			stack = append(stack, token)
		case "(":
			stack = append(stack, token)
		case ")":
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				right := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]
				left := operandStack[len(operandStack)-1]
				operandStack = operandStack[:len(operandStack)-1]
				op := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				expr := Expression{Operator: op, LeftOperand: left, RightOperand: right}
				output = append(output, expr)
				operandStack = append(operandStack, "res")
			}
			stack = stack[:len(stack)-1]
		default:
			operandStack = append(operandStack, token)
		}
	}

	for len(stack) > 0 {
		right := operandStack[len(operandStack)-1]
		operandStack = operandStack[:len(operandStack)-1]
		left := operandStack[len(operandStack)-1]
		operandStack = operandStack[:len(operandStack)-1]
		op := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		expr := Expression{Operator: op, LeftOperand: left, RightOperand: right}
		output = append(output, expr)
		operandStack = append(operandStack, "res")
	}

	return output
}
