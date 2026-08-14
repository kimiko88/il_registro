// Package upload provides utilities for validating file uploads by inspecting
// the actual file content (magic bytes) rather than relying on the HTTP
// Content-Type header or file extension, which an attacker can trivially fake.
package upload

import (
	"errors"
	"io"
	"mime/multipart"

	"github.com/gabriel-vasile/mimetype"
)

// MaxUploadSize is the maximum allowed file size (20 MB).
const MaxUploadSize = 20 << 20 // 20 MiB

// ErrFileTooLarge is returned when the uploaded file exceeds MaxUploadSize.
var ErrFileTooLarge = errors.New("il file supera la dimensione massima consentita (20 MB)")

// ErrFileTypeNotAllowed is returned when the detected MIME type is not in the allowlist.
var ErrFileTypeNotAllowed = errors.New("tipo di file non consentito: solo PDF, immagini e documenti Office sono accettati")

// allowedMIMETypes is the allowlist of accepted MIME types for document uploads.
var allowedMIMETypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	// Microsoft Office Open XML formats (DOCX, XLSX, PPTX)
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	// LibreOffice / OpenDocument formats
	"application/vnd.oasis.opendocument.text":         true,
	"application/vnd.oasis.opendocument.spreadsheet":  true,
	"application/vnd.oasis.opendocument.presentation": true,
}

// ValidateUpload checks both the size and the real MIME type of the uploaded file.
// It reads the first 512 bytes to detect the magic number, then seeks back to the
// beginning so the caller can still read the full file content.
func ValidateUpload(file multipart.File, header *multipart.FileHeader) error {
	// 1. Size check from header (fast rejection if header states oversized)
	if header != nil && header.Size > MaxUploadSize {
		return ErrFileTooLarge
	}

	// 2. Read the first 512 bytes using LimitReader to ensure actual payload doesn't bypass limit
	limitedReader := io.LimitReader(file, MaxUploadSize+1)
	head := make([]byte, 512)
	n, err := limitedReader.Read(head)
	if err != nil && err != io.EOF {
		return errors.New("errore nella lettura del file")
	}
	head = head[:n]

	// Seek back to the beginning so the caller can read the full file.
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return errors.New("errore nel riposizionamento del file")
		}
	}

	// 3. Detect MIME from magic bytes (not from header or extension).
	mtype := mimetype.Detect(head)
	detected := mtype.String()

	if !allowedMIMETypes[detected] {
		return ErrFileTypeNotAllowed
	}

	return nil
}

// DetectMIME inspects the first 512 bytes of file to detect real MIME type using magic bytes, then seeks back.
func DetectMIME(file multipart.File) (string, error) {
	if file == nil {
		return "", errors.New("file is nil")
	}
	head := make([]byte, 512)
	n, err := file.Read(head)
	if err != nil && err != io.EOF {
		return "", err
	}
	head = head[:n]

	if seeker, ok := file.(io.Seeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
	}

	mtype := mimetype.Detect(head)
	return mtype.String(), nil
}
