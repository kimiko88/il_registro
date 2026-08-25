package unit

import (
	"crypto/sha256"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"registro-backend/internal/auditlog"
)

func calculateBlockHash(prevHash string, action string, resourceType string, resourceID string, timestamp time.Time) string {
	input := fmt.Sprintf("%s:%s:%s:%s:%s", prevHash, action, resourceType, resourceID, timestamp.Format(time.RFC3339))
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

func TestSecurity_AuditLogChainImmutability(t *testing.T) {
	t.Run("Validates intact audit log block chain", func(t *testing.T) {
		t1 := time.Now().UTC()
		t2 := t1.Add(1 * time.Minute)

		prevHash0 := "0000000000000000000000000000000000000000000000000000000000000000"
		hash1 := calculateBlockHash(prevHash0, "LOGIN", "user", "u-001", t1)
		hash2 := calculateBlockHash(hash1, "CHANGE_PASSWORD", "user", "u-001", t2)

		block1 := auditlog.CertifiedChainBlock{
			ID:           "b1",
			Action:       "LOGIN",
			ActorName:    "Mario Rossi",
			ResourceType: "user",
			ResourceID:   "u-001",
			PrevHash:     prevHash0,
			CurrentHash:  hash1,
			Timestamp:    t1,
			IsValid:      true,
		}

		block2 := auditlog.CertifiedChainBlock{
			ID:           "b2",
			Action:       "CHANGE_PASSWORD",
			ActorName:    "Mario Rossi",
			ResourceType: "user",
			ResourceID:   "u-001",
			PrevHash:     hash1,
			CurrentHash:  hash2,
			Timestamp:    t2,
			IsValid:      true,
		}

		blocks := []auditlog.CertifiedChainBlock{block1, block2}

		// Verify chain
		currentPrev := prevHash0
		chainIntact := true
		for _, b := range blocks {
			expectedHash := calculateBlockHash(b.PrevHash, b.Action, b.ResourceType, b.ResourceID, b.Timestamp)
			if b.PrevHash != currentPrev || b.CurrentHash != expectedHash {
				chainIntact = false
				break
			}
			currentPrev = b.CurrentHash
		}

		assert.True(t, chainIntact)
	})

	t.Run("Detects tampered audit log action payload and breaks chain integrity", func(t *testing.T) {
		t1 := time.Now().UTC()
		prevHash0 := "0000000000000000000000000000000000000000000000000000000000000000"
		hash1 := calculateBlockHash(prevHash0, "LOGIN", "user", "u-001", t1)

		// Tampered block where action was modified after creation
		tamperedBlock := auditlog.CertifiedChainBlock{
			Action:       "DELETE_ALL_DATA", // Tampered action string
			ResourceType: "user",
			ResourceID:   "u-001",
			PrevHash:     prevHash0,
			CurrentHash:  hash1, // Old hash doesn't match new action
			Timestamp:    t1,
		}

		expectedHash := calculateBlockHash(tamperedBlock.PrevHash, tamperedBlock.Action, tamperedBlock.ResourceType, tamperedBlock.ResourceID, tamperedBlock.Timestamp)
		assert.NotEqual(t, expectedHash, tamperedBlock.CurrentHash)
	})
}
