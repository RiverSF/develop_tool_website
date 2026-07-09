package common_test

import (
	"testing"

	"develop_tools/pkg/common"
)

func TestEvaluateCalcExpression(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1+2*3", "7.000000"},
		{"(1+2)*3", "9.000000"},
		{"10/4", "2.500000"},
	}
	for _, tc := range tests {
		got, err := common.EvaluateCalcExpression(tc.in)
		if err != nil {
			t.Fatalf("input %q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("input %q: got %s want %s", tc.in, got, tc.want)
		}
	}
}

func TestEvaluateCalcExpressionRejectsInvalidInput(t *testing.T) {
	if _, err := common.EvaluateCalcExpression("1;rm"); err == nil {
		t.Fatal("expected error for injected command")
	}
}
