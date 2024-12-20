package calculate

import (
	"testing"
)

func TestCalc(t *testing.T) {
	t.Run("success tests", SuccessCases)
	t.Run("fail tests", FailCases)
}

// Тесты без ошибки
func SuccessCases(t *testing.T) {
	cases := []struct {
		name       string
		expression string
		want       float64
	}{
		{
			name:       "simple-1",
			expression: "1+1",
			want:       2,
		},
		{
			name:       "simple-2",
			expression: "3-6",
			want:       -3,
		},
		{
			name:       "priority-1",
			expression: "2 + 2*2",
			want:       6,
		},
		{
			name:       "priority-2",
			expression: "(2+2) * 2",
			want:       8,
		},
		{
			name:       "divizion",
			expression: "1/2",
			want:       0.5,
		},
		{
			name:       "brackets-1",
			expression: "2*(3+5/10)",
			want:       7,
		},
		{
			name:       "combine-1",
			expression: "23*11-37",
			want:       216,
		},
		{
			name:       "brackets-2",
			expression: "(3-6*5)-(4/2+22)",
			want:       -51,
		},
		{
			name:       "combine-2",
			expression: "2- 44 *5",
			want:       -218,
		},
		{
			name:       "combine-3",
			expression: "45/9*23",
			want:       115,
		},
	}

	for _, cs := range cases {
		got, err := Calc(cs.expression)
		if err != nil {
			t.Fatalf("expected %v, but got error: %v", cs.want, err)
		}
		if got != cs.want {
			t.Fatalf("Calc(%s) = %v, expected %v", cs.expression, got, cs.want)
		}
	}
}

// Тесты с ошибкой
func FailCases(t *testing.T) {
	cases := []struct {
		name       string
		expression string
	}{
		{
			name:       "empty",
			expression: "",
		},
		{
			name:       "first operation",
			expression: "-13+3",
		},
		{
			name:       "last operation",
			expression: "23-4+",
		},
		{
			name:       "empty brackets",
			expression: "2*()-3",
		},
		{
			name:       "invalid token",
			expression: "a +3",
		},
		{
			name:       "operations brackets",
			expression: "2(3-4+)-3",
		},
		{
			name:       "invalid brackets",
			expression: "1)*34",
		},
		{
			name:       "2 operations",
			expression: "54/2++5",
		},
		{
			name:       "2 numbers",
			expression: "2 2-345",
		},
		{
			name:       "combine",
			expression: "+23-(34/)-2-",
		},
	}

	for _, cs := range cases {
		got, err := Calc(cs.expression)
		if err == nil {
			t.Fatalf("expected error, but got %v", got)
		}
	}
}
