package common

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

var errInvalidExpression = errors.New("invalid expression")

// ReplaceCalcSymbols normalizes calculator input symbols.
func ReplaceCalcSymbols(command string) string {
	command = strings.ReplaceAll(command, "x", "*")
	command = strings.ReplaceAll(command, "X", "*")
	command = strings.ReplaceAll(command, "÷", "/")
	return command
}

// EvaluateCalcExpression safely evaluates + - * / and parentheses.
func EvaluateCalcExpression(command string) (string, error) {
	command = strings.TrimSpace(ReplaceCalcSymbols(command))
	if command == "" {
		return "", errInvalidExpression
	}
	for _, r := range command {
		if !unicode.IsDigit(r) && !strings.ContainsRune("+-*/(). ", r) {
			return "", fmt.Errorf("%w: invalid character %q", errInvalidExpression, r)
		}
	}

	val, err := evalCalcExpression(command)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.6f", val), nil
}

func evalCalcExpression(expr string) (float64, error) {
	tokens, err := tokenizeCalc(expr)
	if err != nil {
		return 0, err
	}
	rpn, err := toRPN(tokens)
	if err != nil {
		return 0, err
	}
	return evalRPN(rpn)
}

type calcToken struct {
	kind  byte
	value float64
}

func tokenizeCalc(expr string) ([]calcToken, error) {
	expr = strings.ReplaceAll(expr, " ", "")
	var tokens []calcToken
	for i := 0; i < len(expr); {
		switch expr[i] {
		case '+', '-', '*', '/', '(', ')':
			tokens = append(tokens, calcToken{kind: expr[i]})
			i++
		default:
			j := i
			for j < len(expr) && (unicode.IsDigit(rune(expr[j])) || expr[j] == '.') {
				j++
			}
			if j == i {
				return nil, errInvalidExpression
			}
			num, err := strconv.ParseFloat(expr[i:j], 64)
			if err != nil {
				return nil, errInvalidExpression
			}
			tokens = append(tokens, calcToken{kind: 'n', value: num})
			i = j
		}
	}
	return tokens, nil
}

func toRPN(tokens []calcToken) ([]calcToken, error) {
	var output, stack []calcToken
	precedence := map[byte]int{'+': 1, '-': 1, '*': 2, '/': 2}

	for i, tok := range tokens {
		if tok.kind == 'n' {
			output = append(output, tok)
			continue
		}
		switch tok.kind {
		case '(':
			stack = append(stack, tok)
		case ')':
			for len(stack) > 0 && stack[len(stack)-1].kind != '(' {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errInvalidExpression
			}
			stack = stack[:len(stack)-1]
		case '+', '-', '*', '/':
			if tok.kind == '-' && (i == 0 || tokens[i-1].kind == '(' || tokens[i-1].kind == '+' || tokens[i-1].kind == '-' || tokens[i-1].kind == '*' || tokens[i-1].kind == '/') {
				output = append(output, calcToken{kind: 'n', value: 0})
			}
			for len(stack) > 0 && stack[len(stack)-1].kind != '(' && precedence[stack[len(stack)-1].kind] >= precedence[tok.kind] {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, tok)
		default:
			return nil, errInvalidExpression
		}
	}
	for len(stack) > 0 {
		if stack[len(stack)-1].kind == '(' {
			return nil, errInvalidExpression
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	return output, nil
}

func evalRPN(tokens []calcToken) (float64, error) {
	var stack []float64
	for _, tok := range tokens {
		if tok.kind == 'n' {
			stack = append(stack, tok.value)
			continue
		}
		if len(stack) < 2 {
			return 0, errInvalidExpression
		}
		b := stack[len(stack)-1]
		a := stack[len(stack)-2]
		stack = stack[:len(stack)-2]
		switch tok.kind {
		case '+':
			stack = append(stack, a+b)
		case '-':
			stack = append(stack, a-b)
		case '*':
			stack = append(stack, a*b)
		case '/':
			if b == 0 {
				return 0, errors.New("division by zero")
			}
			stack = append(stack, a/b)
		default:
			return 0, errInvalidExpression
		}
	}
	if len(stack) != 1 {
		return 0, errInvalidExpression
	}
	if math.IsInf(stack[0], 0) || math.IsNaN(stack[0]) {
		return 0, errInvalidExpression
	}
	return stack[0], nil
}
