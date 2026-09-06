package upload

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"strings"
	"testing"
)

// mockFile wraps *bytes.Reader to satisfy multipart.File interface (io.Reader, io.ReaderAt, io.Seeker, io.Closer).
type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}

func newMockFile(content []byte) multipart.File {
	return &mockFile{Reader: bytes.NewReader(content)}
}

func TestValidateUpload_ValidPDF(t *testing.T) {
	// Magic bytes for PDF "%PDF-1.4..."
	pdfContent := []byte("%PDF-1.4\n%âãÏÓ\n1 0 obj\n<< /Type /Catalog >>\nendobj\n")
	file := newMockFile(pdfContent)
	header := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     int64(len(pdfContent)),
	}

	err := ValidateUpload(file, header)
	if err != nil {
		t.Fatalf("expected valid PDF to pass validation, got error: %v", err)
	}
}

func TestValidateUpload_ValidPNG(t *testing.T) {
	// Magic bytes for PNG: 89 50 4E 47 0D 0A 1A 0A
	pngContent := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D}
	file := newMockFile(pngContent)
	header := &multipart.FileHeader{
		Filename: "image.png",
		Size:     int64(len(pngContent)),
	}

	err := ValidateUpload(file, header)
	if err != nil {
		t.Fatalf("expected valid PNG to pass validation, got error: %v", err)
	}
}

func TestValidateUpload_ValidJPEG(t *testing.T) {
	// Magic bytes for JPEG: FF D8 FF E0
	jpegContent := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01}
	file := newMockFile(jpegContent)
	header := &multipart.FileHeader{
		Filename: "photo.jpg",
		Size:     int64(len(jpegContent)),
	}

	err := ValidateUpload(file, header)
	if err != nil {
		t.Fatalf("expected valid JPEG to pass validation, got error: %v", err)
	}
}

func TestValidateUpload_DisallowedExecutable(t *testing.T) {
	// Disallowed shell script or binary content
	scriptContent := []byte("#!/bin/bash\necho 'malicious code'\n")
	file := newMockFile(scriptContent)
	header := &multipart.FileHeader{
		Filename: "script.sh",
		Size:     int64(len(scriptContent)),
	}

	err := ValidateUpload(file, header)
	if err == nil {
		t.Fatalf("expected script file to be rejected, but passed validation")
	}
	if err != ErrFileTypeNotAllowed {
		t.Errorf("expected ErrFileTypeNotAllowed, got: %v", err)
	}
}

func TestValidateUpload_DisallowedHTML(t *testing.T) {
	// HTML content should be rejected to prevent stored XSS via file view
	htmlContent := []byte("<html><body><script>alert(1)</script></body></html>")
	file := newMockFile(htmlContent)
	header := &multipart.FileHeader{
		Filename: "page.html",
		Size:     int64(len(htmlContent)),
	}

	err := ValidateUpload(file, header)
	if err == nil {
		t.Fatalf("expected HTML file to be rejected, but passed validation")
	}
	if err != ErrFileTypeNotAllowed {
		t.Errorf("expected ErrFileTypeNotAllowed, got: %v", err)
	}
}

func TestValidateUpload_OversizedFile(t *testing.T) {
	pdfContent := []byte("%PDF-1.4 test")
	file := newMockFile(pdfContent)
	header := &multipart.FileHeader{
		Filename: "large.pdf",
		Size:     MaxUploadSize + 100, // Exceeds 20 MB
	}

	err := ValidateUpload(file, header)
	if err == nil {
		t.Fatalf("expected oversized file to be rejected, but passed validation")
	}
	if err != ErrFileTooLarge {
		t.Errorf("expected ErrFileTooLarge, got: %v", err)
	}
}

func TestValidateUpload_SeekReset(t *testing.T) {
	pdfContent := []byte("%PDF-1.4\n%âãÏÓ\n1 0 obj\n<< /Type /Catalog >>\nendobj\n")
	file := newMockFile(pdfContent)
	header := &multipart.FileHeader{
		Filename: "test.pdf",
		Size:     int64(len(pdfContent)),
	}

	err := ValidateUpload(file, header)
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	// Verify seeker was reset to position 0 so caller can read full content
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		t.Fatalf("failed to read file after validation: %v", err)
	}

	if !strings.HasPrefix(buf.String(), "%PDF-1.4") {
		t.Errorf("expected file pointer to be reset to 0, but content was: %s", buf.String())
	}
}

func TestValidateAndSanitize_ValidPNG(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	var rawBuf bytes.Buffer
	if err := png.Encode(&rawBuf, img); err != nil {
		t.Fatalf("png encode failed: %v", err)
	}

	rawBytes := rawBuf.Bytes()
	file := newMockFile(rawBytes)
	header := &multipart.FileHeader{
		Filename: "avatar.png",
		Size:     int64(len(rawBytes)),
	}

	sanitized, mime, err := ValidateAndSanitize(file, header)
	if err != nil {
		t.Fatalf("expected valid PNG to pass ValidateAndSanitize, got: %v", err)
	}
	if mime != "image/png" {
		t.Errorf("expected mime image/png, got %s", mime)
	}
	if len(sanitized) == 0 {
		t.Errorf("expected non-empty sanitized bytes")
	}

	// Verify sanitized bytes decode cleanly as image
	decoded, err := png.Decode(bytes.NewReader(sanitized))
	if err != nil {
		t.Fatalf("failed to decode sanitized image: %v", err)
	}
	if decoded.Bounds().Dx() != 8 || decoded.Bounds().Dy() != 8 {
		t.Errorf("unexpected image bounds: %v", decoded.Bounds())
	}
}

func TestValidateAndSanitize_DisallowedFile(t *testing.T) {
	scriptContent := []byte("#!/bin/bash\nrm -rf /")
	file := newMockFile(scriptContent)
	header := &multipart.FileHeader{
		Filename: "script.sh",
		Size:     int64(len(scriptContent)),
	}

	_, _, err := ValidateAndSanitize(file, header)
	if err == nil {
		t.Fatal("expected script file to be rejected by ValidateAndSanitize")
	}
	if err != ErrFileTypeNotAllowed {
		t.Errorf("expected ErrFileTypeNotAllowed, got %v", err)
	}
}
