package upload

import (
	"bytes"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

// StripImageMetadata decodes an image stream and re-encodes it into canonical bytes,
// effectively stripping all EXIF metadata, GPS geotags, camera profiles, and comments.
// If the MIME type is not a supported raster image (JPEG, PNG, GIF), it reads and returns
// the original byte stream intact.
func StripImageMetadata(r io.Reader, mimeType string) ([]byte, error) {
	mimeType = strings.ToLower(strings.TrimSpace(mimeType))
	switch mimeType {
	case "image/jpeg", "image/jpg":
		img, err := jpeg.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decodifica JPEG fallita: %w", err)
		}
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return nil, fmt.Errorf("codifica JPEG fallita: %w", err)
		}
		return buf.Bytes(), nil

	case "image/png":
		img, err := png.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decodifica PNG fallita: %w", err)
		}
		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("codifica PNG fallita: %w", err)
		}
		return buf.Bytes(), nil

	case "image/gif":
		img, err := gif.Decode(r)
		if err != nil {
			return nil, fmt.Errorf("decodifica GIF fallita: %w", err)
		}
		var buf bytes.Buffer
		if err := gif.Encode(&buf, img, nil); err != nil {
			return nil, fmt.Errorf("codifica GIF fallita: %w", err)
		}
		return buf.Bytes(), nil

	default:
		// Non-image format (e.g. PDF, CSV, docx); pass through unaltered
		return io.ReadAll(r)
	}
}

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
