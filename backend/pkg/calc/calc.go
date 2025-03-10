package calc

import (
	"calc-website/pkg/utils"
	"errors"
	"unicode"
)

var (
	ErrDivisionByZero    = errors.New("division by zero")
	ErrUnknownOperator   = errors.New("unknown operator")
	ErrExpressionInvalid = errors.New("expression is invalid")
)

type Node struct {
	Value string
	Right *Node
	Left  *Node
}

var OperationPriorities = map[string]int{
	"*": 1,
	"/": 1,
	"-": 2,
	"+": 2,
}

func Compute(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrUnknownOperator
	}
}

func checkExpression(expression string) bool {
	for _, symbol := range expression {
		_, check := OperationPriorities[string(symbol)]
		if !unicode.IsDigit(symbol) &&
			!check &&
			symbol != ' ' &&
			symbol != '(' &&
			symbol != ')' {
			return false
		}
	}
	return true
}

func hasDivisionByZero(root *Node) bool {
	return root.Right.Value == "0" && root.Value == "/"
}

func tokenize(expression string) []string {
	index := 0
	expressionLength := len(expression)

	var output []string

	for index < expressionLength {
		symbol := rune(expression[index])
		if unicode.IsDigit(symbol) {
			cur := index
			for cur < expressionLength-1 && unicode.IsDigit(rune(expression[cur+1])) {
				cur++
			}

			if cur == expressionLength-1 {
				output = append(output, expression[index:])
			} else {
				output = append(output, expression[index:cur+1])
			}

			index = cur + 1
			continue
		} else if symbol != ' ' {
			output = append(output, string(symbol))
		}
		index++
	}

	return output
}

func ToTree(expression string) (Node, error) {
	var operands []Node
	var operators []string

	if !checkExpression(expression) {
		return Node{}, ErrExpressionInvalid
	}
	tokens := tokenize(expression)
	if len(tokens) < 3 {
		return Node{}, ErrExpressionInvalid
	}

	for _, token := range tokens {
		priority, op := OperationPriorities[token]
		if op {
			for len(operators) > 0 &&
				operators[len(operators)-1] != "(" &&
				OperationPriorities[operators[len(operators)-1]] <= priority {

				op, err := utils.Pop(&operators)
				if err != nil {
					return Node{}, err
				}

				right, err := utils.Pop(&operands)
				if err != nil {
					return Node{}, err
				}

				left, err := utils.Pop(&operands)
				if err != nil {
					return Node{}, err
				}

				node := &Node{Value: op}
				node.Left = &left
				node.Right = &right
				operands = append(operands, *node)
			}
			operators = append(operators, token)
		} else if utils.IsNumber(token) {
			operands = append(operands, Node{Value: token})
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 &&
				operators[len(operators)-1] != "(" {

				op, err := utils.Pop(&operators)
				if err != nil {
					return Node{}, err
				}

				right, err := utils.Pop(&operands)
				if err != nil {
					return Node{}, err
				}

				left, err := utils.Pop(&operands)
				if err != nil {
					return Node{}, err
				}
				node := &Node{Value: op}
				node.Left = &left
				node.Right = &right
				operands = append(operands, *node)
			}

			_, err := utils.Pop(&operators)

			if err != nil {
				return Node{}, err
			}
		}
	}

	for len(operators) > 0 {
		op, err := utils.Pop(&operators)
		if err != nil {
			return Node{}, err
		}
		right, err := utils.Pop(&operands)
		if err != nil {
			return Node{}, err
		}
		left, err := utils.Pop(&operands)
		if err != nil {
			return Node{}, err
		}
		node := &Node{Value: op}
		node.Left = &left
		node.Right = &right
		operands = append(operands, *node)
	}

	if size := len(operands); size > 0 {
		root := operands[len(operands)-1]

		if hasDivisionByZero(&root) {
			return Node{}, ErrDivisionByZero
		}

		return root, nil
	}
	return Node{}, nil
}

func (n *Node) Infix() string {
	if n == nil {
		return ""
	}
	if n.Left == nil && n.Right == nil {
		return n.Value
	}
	return "(" + n.Left.Infix() + " " + n.Value + " " + n.Right.Infix() + ")"
}
