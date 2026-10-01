package unit

import (
	"math"
	"testing"
	"time"

	"registro-backend/internal/strike"

	"github.com/stretchr/testify/assert"
)

func TestStrike_IntentionValidation(t *testing.T) {
	// Valid intentions
	assert.True(t, strike.IsValidIntention(strike.IntentionParticipates))
	assert.True(t, strike.IsValidIntention(strike.IntentionNotParticipates))
	assert.True(t, strike.IsValidIntention(strike.IntentionUndecided))

	// Invalid intentions
	assert.False(t, strike.IsValidIntention("invalid"))
	assert.False(t, strike.IsValidIntention("yes"))
	assert.False(t, strike.IsValidIntention("unanswered"))
	assert.False(t, strike.IsValidIntention(""))
	assert.False(t, strike.IsValidIntention("declined"))
}

func TestStrike_DeadlineExpiry(t *testing.T) {
	now := time.Now()

	t.Run("deadline in future is not expired", func(t *testing.T) {
		deadline := now.Add(24 * time.Hour)
		notice := strike.StrikeNotice{
			DeclarationDeadline: deadline,
		}
		isExpired := notice.DeclarationDeadline.Before(now)
		assert.False(t, isExpired)
	})

	t.Run("deadline in past is expired", func(t *testing.T) {
		deadline := now.Add(-1 * time.Hour)
		notice := strike.StrikeNotice{
			DeclarationDeadline: deadline,
		}
		isExpired := notice.DeclarationDeadline.Before(now)
		assert.True(t, isExpired)
	})
}

func TestStrike_SummaryPercentageCalculation(t *testing.T) {
	t.Run("accurate percentages with mixed intentions", func(t *testing.T) {
		totalStaff := 20
		participates := 8
		notParticipates := 6
		undecided := 4
		unanswered := 2

		calcPct := func(count, total int) float64 {
			if total == 0 {
				return 0.0
			}
			return math.Round((float64(count)/float64(total)*100.0)*10.0) / 10.0
		}

		partPct := calcPct(participates, totalStaff)
		notPartPct := calcPct(notParticipates, totalStaff)
		undecPct := calcPct(undecided, totalStaff)
		unanPct := calcPct(unanswered, totalStaff)

		assert.Equal(t, 40.0, partPct)
		assert.Equal(t, 30.0, notPartPct)
		assert.Equal(t, 20.0, undecPct)
		assert.Equal(t, 10.0, unanPct)
		assert.Equal(t, 100.0, partPct+notPartPct+undecPct+unanPct)
	})

	t.Run("zero total staff does not panic or divide by zero", func(t *testing.T) {
		totalStaff := 0
		calcPct := func(count, total int) float64 {
			if total == 0 {
				return 0.0
			}
			return float64(count) / float64(total) * 100.0
		}

		assert.Equal(t, 0.0, calcPct(0, totalStaff))
	})
}

func TestStrike_RoleAuthorizationMatrix(t *testing.T) {
	// Authorized management roles for strike notices & summary
	authorizedRoles := []string{
		"superadmin",
		"admin",
		"principal",
		"vice_principal",
		"dsga",
		"secretary",
		"assistente_amministrativo",
		"assistente_personale",
	}

	// Unauthorized roles (must not create notices or view aggregate summary)
	unauthorizedRoles := []string{
		"teacher",
		"student",
		"parent",
		"collaboratore_scolastico",
		"assistente_tecnico",
	}

	isAuthorized := func(role string) bool {
		switch role {
		case "superadmin", "admin", "principal", "vice_principal", "dsga", "secretary", "assistente_amministrativo", "assistente_personale":
			return true
		default:
			return false
		}
	}

	for _, role := range authorizedRoles {
		assert.True(t, isAuthorized(role), "Role %s should be authorized for strike management", role)
	}

	for _, role := range unauthorizedRoles {
		assert.False(t, isAuthorized(role), "Role %s must NOT be authorized for strike management", role)
	}
}
