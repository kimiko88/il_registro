package reports

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeExcelField(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"=cmd|' /C calc'!A0", "\t=cmd|' /C calc'!A0"},
		{"+12345", "\t+12345"},
		{"-SUM(A1:A10)", "\t-SUM(A1:A10)"},
		{"@SUM(A1:A10)", "\t@SUM(A1:A10)"},
		{"Mario Rossi", "Mario Rossi"},
		{"Matematica", "Matematica"},
		{"", ""},
	}

	for _, tc := range cases {
		actual := sanitizeExcelField(tc.input)
		assert.Equal(t, tc.expected, actual, "input=%s", tc.input)
	}
}
