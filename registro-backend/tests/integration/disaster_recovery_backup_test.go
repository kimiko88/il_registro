package integration

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type AuditLogEntry struct {
	ID        string `json:"id"`
	Action    string `json:"action"`
	PrevHash  string `json:"prev_hash"`
	EntryHash string `json:"entry_hash"`
}

func calculateHash(prevHash, action string) string {
	h := sha256.Sum256([]byte(prevHash + ":" + action))
	return fmt.Sprintf("%x", h)
}

func TestIntegration_Disaster_Recovery_Backup(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auditChain := make([]AuditLogEntry, 0)
	genesisHash := "0000000000000000000000000000000000000000000000000000000000000000"

	// Create initial chain entry
	hash1 := calculateHash(genesisHash, "USER_LOGIN")
	auditChain = append(auditChain, AuditLogEntry{
		ID:        "log-1",
		Action:    "USER_LOGIN",
		PrevHash:  genesisHash,
		EntryHash: hash1,
	})

	// Add second entry
	hash2 := calculateHash(hash1, "GRADE_MODIFIED")
	auditChain = append(auditChain, AuditLogEntry{
		ID:        "log-2",
		Action:    "GRADE_MODIFIED",
		PrevHash:  hash1,
		EntryHash: hash2,
	})

	r := gin.New()
	r.GET("/audit/verify", func(c *gin.Context) {
		// Verify hash chain continuity
		currPrevHash := genesisHash
		for _, entry := range auditChain {
			if entry.PrevHash != currPrevHash {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "audit log hash chain broken"})
				return
			}
			expectedHash := calculateHash(currPrevHash, entry.Action)
			if entry.EntryHash != expectedHash {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "audit log entry hash mismatch"})
				return
			}
			currPrevHash = entry.EntryHash
		}
		c.JSON(http.StatusOK, gin.H{"status": "intact", "chain_length": len(auditChain)})
	})

	// 1. Verify unbroken audit log chain
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/audit/verify", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

	// 2. Simulate node restart / recovery and add post-recovery entry
	hash3 := calculateHash(hash2, "RECOVERY_NODE_RESTART")
	auditChain = append(auditChain, AuditLogEntry{
		ID:        "log-3",
		Action:    "RECOVERY_NODE_RESTART",
		PrevHash:  hash2,
		EntryHash: hash3,
	})

	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/audit/verify", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)
}
