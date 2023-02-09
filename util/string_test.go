package util

import (
	"testing"
)

func TestFirstNonEmpty(t *testing.T) {
	tests := []struct {
		strings []string
		result  string
	}{
		{[]string{"", "hello", "world"}, "hello"},
		{[]string{"", "", "hello"}, "hello"},
		{[]string{"", "", ""}, ""},
		{[]string{"hello", "world", "golang"}, "hello"},
		{[]string{"golang", "world", "hello"}, "golang"},
	}

	for _, test := range tests {
		res := FirstNonEmpty(test.strings)
		if res != test.result {
			t.Errorf("FirstNonEmpty(%v) = %v, want %v", test.strings, res, test.result)
		}
	}
}
