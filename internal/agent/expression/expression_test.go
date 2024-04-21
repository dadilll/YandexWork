package expression

import (
	"testing"
	"time"
)

func TestParseExpression(t *testing.T) {
	durationMap := map[string]int{
		"+": 1,
		"-": 1,
		"*": 1,
		"/": 1,
	}

	tests := []struct {
		expression string
		expected   float64
	}{
		{"1 + 2", 3},
		{"5 - 3", 2},
		{"2 * 3", 6},
		{"10 / 2", 5},
		{"3 + 4 * 2", 11},
		{"4 * 2 / 2", 4},
		{"6 / 0", 0},   // Тест на деление на ноль
		{"abc + 2", 0}, // Тест на некорректное выражение
	}

	for _, test := range tests {
		start := time.Now() // Measure start time
		result, err := ParseExpression(test.expression, durationMap)
		elapsed := time.Since(start) // Calculate elapsed time
		if err != nil {
			if test.expected != 0 {
				t.Errorf("Unexpected error while parsing expression '%s': %v", test.expression, err)
			}
			continue
		}
		if result != test.expected {
			t.Errorf("Incorrect result for expression '%s'. Expected: %f, Got: %f", test.expression, test.expected, result)
		}
		t.Logf("Expression '%s' executed in %s", test.expression, elapsed)
	}
}

func TestParseExpression_InvalidCharacters(t *testing.T) {
	durationMap := map[string]int{
		"+": 1,
		"-": 1,
		"*": 1,
		"/": 1,
	}

	tests := []struct {
		expression string
	}{
		{"1 + @ 2"}, // Некорректный символ @
		{"5 - 3 ?"}, // Некорректный символ ?
	}

	for _, test := range tests {
		_, err := ParseExpression(test.expression, durationMap)
		if err == nil {
			t.Errorf("Expected error for expression '%s' with invalid characters, but got none", test.expression)
		}
	}
}
