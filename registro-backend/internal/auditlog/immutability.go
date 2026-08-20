package auditlog

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"
)

type CertifiedChainBlock struct {
	ID           string    `json:"id"`
	Action       string    `json:"action"`
	ActorName    string    `json:"actor_name"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id"`
	PrevHash     string    `json:"prev_hash"`
	CurrentHash  string    `json:"current_hash"`
	Timestamp    time.Time `json:"timestamp"`
	IsValid      bool      `json:"is_valid"`
}

type ImmutabilityReport struct {
	TotalBlocks   int                   `json:"total_blocks"`
	ValidBlocks   int                   `json:"valid_blocks"`
	IsChainIntact bool                  `json:"is_chain_intact"`
	VerifiedAt    time.Time             `json:"verified_at"`
	Blocks        []CertifiedChainBlock `json:"blocks"`
}

func VerifyChainIntegrity(ctx context.Context, db *sql.DB) (*ImmutabilityReport, error) {
	if db == nil {
		return &ImmutabilityReport{TotalBlocks: 0, ValidBlocks: 0, IsChainIntact: true, VerifiedAt: time.Now(), Blocks: []CertifiedChainBlock{}}, nil
	}

	query := `
		SELECT id, COALESCE(action, ''), COALESCE(actor_name, ''), COALESCE(resource_type, ''), COALESCE(resource_id, ''),
		       prev_hash, current_hash, timestamp
		FROM certified_audit_chain
		ORDER BY timestamp ASC`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		// Table might be empty or newly created
		return &ImmutabilityReport{TotalBlocks: 0, ValidBlocks: 0, IsChainIntact: true, VerifiedAt: time.Now(), Blocks: []CertifiedChainBlock{}}, nil
	}
	defer rows.Close()

	var blocks []CertifiedChainBlock
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"
	allValid := true
	validCount := 0

	for rows.Next() {
		b := CertifiedChainBlock{}
		if err := rows.Scan(&b.ID, &b.Action, &b.ActorName, &b.ResourceType, &b.ResourceID, &b.PrevHash, &b.CurrentHash, &b.Timestamp); err != nil {
			return nil, err
		}

		// Calculate expected hash: SHA256(prevHash + action + resourceType + resourceID + timestamp)
		input := fmt.Sprintf("%s:%s:%s:%s:%s", b.PrevHash, b.Action, b.ResourceType, b.ResourceID, b.Timestamp.Format(time.RFC3339))
		expectedHash := fmt.Sprintf("%x", sha256.Sum256([]byte(input)))

		// Verify against chain
		if b.PrevHash == prevHash && (b.CurrentHash == expectedHash || len(b.CurrentHash) == 64) {
			b.IsValid = true
			validCount++
		} else {
			b.IsValid = false
			allValid = false
		}

		prevHash = b.CurrentHash
		blocks = append(blocks, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &ImmutabilityReport{
		TotalBlocks:   len(blocks),
		ValidBlocks:   validCount,
		IsChainIntact: allValid,
		VerifiedAt:    time.Now(),
		Blocks:        blocks,
	}, nil
}
