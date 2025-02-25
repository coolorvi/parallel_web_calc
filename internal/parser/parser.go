package orchestrator

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
	subExpressions := extractSubExpressions(tokens)

	output := OutputJSON{SubExpressions: subExpressions}
	return json.Marshal(output)
}

func tokenize(expr string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, c := range expr {
		if unicode.IsDigit(c) {
			currentToken.WriteRune(c)
		} else if strings.ContainsRune("+-*/", c) {
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

func extractSubExpressions(tokens []string) []Expression {
	var subExpressions []Expression
	for i := 0; i < len(tokens)-2; i += 2 {
		subExpressions = append(subExpressions, Expression{
			Operator:     tokens[i+1],
			LeftOperand:  tokens[i],
			RightOperand: tokens[i+2],
		})
	}
	return subExpressions
}
