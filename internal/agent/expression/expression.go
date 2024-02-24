package expression

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Token struct {
	Type  string
	Value string
}

type Operation struct {
	Operator        string
	DurationSeconds int
}

func TokenizeExpression(expression string) ([]Token, error) {
	var tokens []Token
	var buffer strings.Builder

	for _, char := range expression {
		switch {
		case char == '+' || char == '-' || char == '*' || char == '/':
			if buffer.Len() > 0 {
				tokens = append(tokens, Token{Type: "number", Value: buffer.String()})
				buffer.Reset()
			}
			tokens = append(tokens, Token{Type: "operator", Value: string(char)})
		case char >= '0' && char <= '9' || char == '.':
			buffer.WriteRune(char)
		case char == ' ':
			continue
		default:
			return nil, fmt.Errorf("invalid character in expression: %c", char)
		}
	}

	if buffer.Len() > 0 {
		tokens = append(tokens, Token{Type: "number", Value: buffer.String()})
	}

	return tokens, nil
}

func ParseExpression(expression string) (float64, error) {
	tokens, err := TokenizeExpression(expression)
	if err != nil {
		return 0, err
	}

	for _, op := range []string{"*/", "+-"} {
		for {
			found := false
			for i, token := range tokens {
				if token.Type == "operator" && strings.Contains(op, token.Value) {
					op1, err := strconv.ParseFloat(tokens[i-1].Value, 64)
					if err != nil {
						return 0, fmt.Errorf("invalid operand 1: %s", tokens[i-1].Value)
					}
					op2, err := strconv.ParseFloat(tokens[i+1].Value, 64)
					if err != nil {
						return 0, fmt.Errorf("invalid operand 2: %s", tokens[i+1].Value)
					}

					result, err := EvaluateExpression(op1, op2, token.Value, 0)
					if err != nil {
						return 0, err
					}

					// Заменяем оператор и его операнды результатом вычисления
					tokens = append(tokens[:i-1], append([]Token{{Type: "number", Value: fmt.Sprintf("%f", result)}}, tokens[i+2:]...)...)
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
	}

	// Возвращаем результат
	return strconv.ParseFloat(tokens[0].Value, 64)
}

func EvaluateExpression(op1, op2 float64, operator string, duration int) (float64, error) {
	time.Sleep(time.Duration(duration) * time.Second)

	var result float64
	switch operator {
	case "+":
		result = op1 + op2
	case "-":
		result = op1 - op2
	case "*":
		result = op1 * op2
	case "/":
		if op2 == 0 {
			return 0, errors.New("division by zero")
		}
		result = op1 / op2
	default:
		return 0, fmt.Errorf("unsupported operator: %s", operator)
	}

	return result, nil
}
