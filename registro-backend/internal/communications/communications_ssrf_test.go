package communications

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// checkAttachmentURLSafe verifica se un URL e sicuro contro SSRF.
// Duplica esattamente la logica di SendMessage per permettere unit test isolati.
func checkAttachmentURLSafe(urlStr string) error {
	urlStr = strings.TrimSpace(urlStr)
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return fmt.Errorf("attachment_url non valido: deve iniziare con http:// o https://")
	}
	parsedAtt, parseErr := url.Parse(urlStr)
	if parseErr != nil || parsedAtt.Host == "" {
		return fmt.Errorf("attachment_url non valido")
	}
	hostname := strings.ToLower(parsedAtt.Hostname())
	if strings.Contains(hostname, "localhost") ||
		hostname == "0.0.0.0" ||
		hostname == "[::1]" || hostname == "::1" ||
		strings.HasPrefix(hostname, "127.") ||
		hostname == "169.254.169.254" ||
		strings.HasPrefix(hostname, "10.") ||
		strings.HasPrefix(hostname, "192.168.") ||
		func() bool {
			parts := strings.Split(hostname, ".")
			if len(parts) != 4 || parts[0] != "172" {
				return false
			}
			var second int
			if _, err := fmt.Sscanf(parts[1], "%d", &second); err != nil {
				return false
			}
			return second >= 16 && second <= 31
		}() {
		return fmt.Errorf("attachment_url non valido: indirizzo interno o privato non consentito (SSRF protection)")
	}
	return nil
}

func TestSSRFProtection_PrivateSubnets_Blocked(t *testing.T) {
	blocked := []string{
		"https://localhost/file.pdf",
		"http://127.0.0.1/secret",
		"https://192.168.1.10/data",
		"http://10.0.0.1/internal",
		"https://172.20.0.5/docker",
		"http://172.16.1.1/docker",
		"https://172.31.255.255/docker",
		"http://0.0.0.0/x",
		"https://169.254.169.254/latest/meta-data/",
	}
	for _, u := range blocked {
		err := checkAttachmentURLSafe(u)
		assert.Error(t, err, "URL=%s deve essere bloccato", u)
		assert.Contains(t, err.Error(), "SSRF")
	}
}

func TestSSRFProtection_PublicURLs_Allowed(t *testing.T) {
	allowed := []string{
		"https://storage.googleapis.com/bucket/file.pdf",
		"http://cdn.example.com/docs/letter.pdf",
		"https://172.15.0.1/ok",
		"http://172.32.0.1/ok",
		"https://11.0.0.1/ok",
	}
	for _, u := range allowed {
		err := checkAttachmentURLSafe(u)
		assert.NoError(t, err, "URL=%s deve essere permesso", u)
	}
}

func TestSSRFProtection_InvalidScheme_Blocked(t *testing.T) {
	err := checkAttachmentURLSafe("ftp://example.com/file.pdf")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deve iniziare con http:// o https://")
}
