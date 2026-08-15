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
