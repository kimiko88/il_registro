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
