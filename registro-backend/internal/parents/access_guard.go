package parents

import (
	"context"
)

type AccessGuard struct {
	repo DualSignatureRepository
}

func NewAccessGuard(repo DualSignatureRepository) *AccessGuard {
	return &AccessGuard{repo: repo}
}

// CanAccessStudent verifies whether a parent has legal authorization to access student records.
// If court orders place a restriction (e.g. Tribunale per i Minorenni), access is blocked.
func (g *AccessGuard) CanAccessStudent(ctx context.Context, parentUserID, studentID string) (bool, error) {
	info, err := g.repo.GetCustodyInfo(ctx, parentUserID, studentID)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}

	// Restricted custody completely bars parent access
	if info.CustodyType == CustodyRestricted {
		return false, nil
	}

	return true, nil
}

// CanAuthorizeActivity verifies if the parent has legal rights to sign authorizations (e.g. trips, consent).
func (g *AccessGuard) CanAuthorizeActivity(ctx context.Context, parentUserID, studentID string) (bool, error) {
	info, err := g.repo.GetCustodyInfo(ctx, parentUserID, studentID)
	if err != nil {
		return false, err
	}
	if info == nil {
		return false, nil
	}

	if info.CustodyType == CustodyRestricted {
		return false, nil
	}

	return info.CanAuthorizeActivities, nil
}
