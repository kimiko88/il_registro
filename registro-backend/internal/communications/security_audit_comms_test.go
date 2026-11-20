package communications

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommunications_SendMessage_SSRFProtection(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepoForComms)
	svc := NewService(mockRepo, mockUsers)

	ctx := context.Background()

	// Malicious internal metadata / localhost URLs MUST be rejected
	invalidURLs := []string{
		"http://169.254.169.254/latest/meta-data/",
		"http://localhost:6379/",
		"http://127.0.0.1:8080/admin",
		"http://10.0.0.1/internal",
		"http://192.168.1.1/router",
	}

	for _, badURL := range invalidURLs {
		urlCopy := badURL
		_, err := svc.SendMessage(ctx, "teacher", "school-1", "teacher-1", CreateMessageRequest{
			Subject:       "Test SSRF",
			Body:          "Testing URL",
			Recipients:    []string{"user-1"},
			AttachmentURL: &urlCopy,
		})
		assert.Error(t, err, "expected error for URL: %s", badURL)
		assert.Contains(t, err.Error(), "SSRF protection")
	}
}

func TestCommunications_SendMessage_RejectsForeignPrivateAttachmentKeys(t *testing.T) {
	mockRepo := new(MockRepository)
	mockUsers := new(MockUserRepoForComms)
	svc := NewService(mockRepo, mockUsers)

	foreignKeys := []string{
		"supabase://documents/communications/school-2/teacher-1/file.pdf",
		"supabase://documents/communications/school-1/other-user/file.pdf",
		"supabase://documents/communications/school-1/teacher-1/../secret.pdf",
	}
	for _, key := range foreignKeys {
		attachment := key
		_, err := svc.SendMessage(context.Background(), "teacher", "school-1", "teacher-1", CreateMessageRequest{
			Subject: "Private attachment", Body: "Test", Recipients: []string{"user-1"}, AttachmentURL: &attachment,
		})
		assert.Error(t, err, "foreign or unsafe key should be rejected: %s", key)
	}
}

func TestCommunicationAttachmentRequiresMatchingExtensionAndMIME(t *testing.T) {
	allowed := []struct {
		filename string
		mime     string
	}{
		{"atto.pdf", "application/pdf"},
		{"foto.JPG", "image/jpeg"},
		{"foto.png", "image/png"},
		{"documento.docx", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
		{"foglio.xlsx", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
		{"presentazione.pptx", "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
		{"nota.txt", "text/plain"},
	}
	for _, item := range allowed {
		assert.True(t, isAllowedCommunicationAttachment(item.filename, item.mime), "%s / %s", item.filename, item.mime)
	}

	blocked := []struct {
		filename string
		mime     string
	}{
		{"payload.exe", "application/pdf"},
		{"payload.png", "application/pdf"},
		{"payload.pdf", "image/png"},
		{"payload.svg", "image/svg+xml"},
	}
	for _, item := range blocked {
		assert.False(t, isAllowedCommunicationAttachment(item.filename, item.mime), "%s / %s", item.filename, item.mime)
	}
}

func TestCommunicationAttachmentUploadRoles(t *testing.T) {
	for _, role := range []string{"admin", "superadmin", "teacher", "coordinator", "coordinatore_classe", "secretary", "principal", "vice_principal", "collaboratore_ds"} {
		assert.True(t, isCommunicationSenderRole(role), "role %s should be allowed", role)
	}
	for _, role := range []string{"parent", "student", "company_tutor", "psychologist", ""} {
		assert.False(t, isCommunicationSenderRole(role), "role %s should be denied", role)
	}
}
