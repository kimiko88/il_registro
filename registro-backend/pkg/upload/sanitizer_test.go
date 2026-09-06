package upload

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
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

func TestStripImageMetadata_PNG(t *testing.T) {
	// Create a simple in-memory RGBA image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}

	var rawBuf bytes.Buffer
	err := png.Encode(&rawBuf, img)
	assert.NoError(t, err)

	sanitized, err := StripImageMetadata(&rawBuf, "image/png")
	assert.NoError(t, err)
	assert.NotEmpty(t, sanitized)

	// Verify the sanitized bytes can be decoded as a valid PNG
	decoded, err := png.Decode(bytes.NewReader(sanitized))
	assert.NoError(t, err)
	assert.Equal(t, 10, decoded.Bounds().Dx())
	assert.Equal(t, 10, decoded.Bounds().Dy())
}

func TestStripImageMetadata_JPEG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			img.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}

	var rawBuf bytes.Buffer
	err := jpeg.Encode(&rawBuf, img, &jpeg.Options{Quality: 85})
	assert.NoError(t, err)

	sanitized, err := StripImageMetadata(&rawBuf, "image/jpeg")
	assert.NoError(t, err)
	assert.NotEmpty(t, sanitized)

	decoded, err := jpeg.Decode(bytes.NewReader(sanitized))
	assert.NoError(t, err)
	assert.Equal(t, 16, decoded.Bounds().Dx())
	assert.Equal(t, 16, decoded.Bounds().Dy())
}

func TestStripImageMetadata_Passthrough(t *testing.T) {
	pdfBytes := []byte("%PDF-1.4 dummy content")
	sanitized, err := StripImageMetadata(bytes.NewReader(pdfBytes), "application/pdf")
	assert.NoError(t, err)
	assert.Equal(t, pdfBytes, sanitized)
}
