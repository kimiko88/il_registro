package upload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeHeaderFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal filename",
			input:    "report_2026.pdf",
			expected: "report_2026.pdf",
		},
		{
			name:     "CRLF Header Injection attempt",
			input:    "report\r\nSet-Cookie: session=evil.pdf",
			expected: "report__Set-Cookie: session=evil.pdf",
		},
		{
			name:     "Double quotes and semicolon injection",
			input:    "file\"; filename=\"injected.exe",
			expected: "file__ filename=_injected.exe",
		},
		{
			name:     "Directory traversal attempt",
			input:    "../../../../etc/passwd",
			expected: "passwd",
		},
		{
			name:     "Windows path traversal attempt",
			input:    `..\..\..\boot.ini`,
			expected: "boot.ini",
		},
		{
			name:     "Empty filename fallback",
			input:    "",
			expected: "download.bin",
		},
		{
			name:     "All unsafe characters",
			input:    "\r\n\"';\\/",
			expected: "download.bin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeHeaderFilename(tt.input)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestFormatContentDisposition(t *testing.T) {
	header := FormatContentDisposition("pagella_stud-123.pdf")
	assert.Equal(t, "attachment; filename=\"pagella_stud-123.pdf\"", header)

	malicious := FormatContentDisposition("pagella\r\nInjected-Header: 123.pdf")
	assert.NotContains(t, malicious, "\r")
	assert.NotContains(t, malicious, "\n")
}
