package upload

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var unsafeHeaderChars = regexp.MustCompile(`[\r\n"';\\/]`)

// SanitizeHeaderFilename removes control characters, newlines, quotes, path separators,
// and other dangerous characters from filenames used in HTTP Content-Disposition headers.
func SanitizeHeaderFilename(filename string) string {
	// Extract only the base filename to prevent directory traversal
	cleaned := filepath.Base(filename)

	// Remove CRLF, double quotes, single quotes, semicolons, backslashes and slashes
	cleaned = unsafeHeaderChars.ReplaceAllString(cleaned, "_")

	// Filter out non-printable ASCII / control characters (< 32 or 127)
	cleaned = strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1
		}
		return r
	}, cleaned)

	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.Trim(cleaned, "._-")

	if cleaned == "" {
		return "download.bin"
	}

	// Limit length to 255 characters to comply with standard filesystem/header limits
	if len(cleaned) > 255 {
		ext := filepath.Ext(cleaned)
		base := cleaned[:255-len(ext)]
		cleaned = base + ext
	}

	return cleaned
}

// FormatContentDisposition returns a safe Content-Disposition header value with attachment.
func FormatContentDisposition(filename string) string {
	safeName := SanitizeHeaderFilename(filename)
	return fmt.Sprintf("attachment; filename=\"%s\"", safeName)
}
