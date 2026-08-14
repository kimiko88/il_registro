package upload

import (
	"bytes"
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
