package calc

import (
	"errors"
	"testing"
)

func TestCompute(t *testing.T) {
	tests := []struct {
		a, b     float64
		op       string
		expected float64
		err      error
	}{
		{3, 2, "+", 5, nil},
		{3, 2, "-", 1, nil},
		{3, 2, "*", 6, nil},
		{6, 2, "/", 3, nil},
		{3, 0, "/", 0, ErrDivisionByZero},
		{3, 2, "%", 0, ErrUnknownOperator},
	}

	for _, tc := range tests {
		res, err := Compute(tc.a, tc.b, tc.op)
		if !errors.Is(err, tc.err) {
			t.Errorf("Compute(%v, %v, %q) error = %v, expected %v", tc.a, tc.b, tc.op, err, tc.err)
		}
		if err == nil && res != tc.expected {
			t.Errorf("Compute(%v, %v, %q) = %v, expected %v", tc.a, tc.b, tc.op, res, tc.expected)
		}
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"3+4", []string{"3", "+", "4"}},
		{"12 * 3", []string{"12", "*", "3"}},
		{"(1+2)", []string{"(", "1", "+", "2", ")"}},
		{"3 + 4 * (2 - 1)", []string{"3", "+", "4", "*", "(", "2", "-", "1", ")"}},
	}

	for _, tc := range tests {
		tokens := tokenize(tc.input)
		if len(tokens) != len(tc.expected) {
			t.Errorf("tokenize(%q) = %v, expected %v", tc.input, tokens, tc.expected)
			continue
		}
		for i, token := range tokens {
			if token != tc.expected[i] {
				t.Errorf("tokenize(%q)[%d] = %q, expected %q", tc.input, i, token, tc.expected[i])
			}
		}
	}
}

func TestCheckExpression(t *testing.T) {
	tests := []struct {
		expr  string
		valid bool
	}{
		{"3 + 4", true},
		{"3a+4", false},
		{"3 + 4 * (2 - 1)", true},
		{"3 + 4 & 5", false},
	}

	for _, tc := range tests {
		if got := checkExpression(tc.expr); got != tc.valid {
			t.Errorf("checkExpression(%q) = %v, expected %v", tc.expr, got, tc.valid)
		}
	}
}

func TestToTree(t *testing.T) {
	tests := []struct {
		expression    string
		expectedInfix string
		shouldError   bool
	}{
		{"3+4*2", "(3 + (4 * 2))", false},
		{"(1+2)*3", "((1 + 2) * 3)", false},
		{"3+(4*2)", "(3 + (4 * 2))", false},
		{"3++4", "", true},
		{"3+(4", "", true},
		{"3+4)", "", true},
	}

	for _, tc := range tests {
		tree, err := ToTree(tc.expression)
		if tc.shouldError {
			if err == nil {
				t.Errorf("ToTree(%q) expected error, got nil", tc.expression)
			}
		} else {
			if err != nil {
				t.Errorf("ToTree(%q) returned error: %v", tc.expression, err)
			} else {
				infix := tree.Infix()
				if infix != tc.expectedInfix {
					t.Errorf("ToTree(%q).Infix() = %q, expected %q", tc.expression, infix, tc.expectedInfix)
				}
			}
		}
	}
}
