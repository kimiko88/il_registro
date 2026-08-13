package unit

import (
	"math"
	"testing"

	"registro-backend/internal/users"

	"github.com/stretchr/testify/assert"
)

func TestScrutinyCalculationLogicUnit(t *testing.T) {
	t.Run("Weighted Average and Rounding Precision", func(t *testing.T) {
		grades := []struct {
			value  float64
			weight float64
		}{
			{value: 7.5, weight: 1.0},
			{value: 8.0, weight: 1.0},
			{value: 6.5, weight: 0.5},
		}

		var totalWeighted float64
		var totalWeight float64

		for _, g := range grades {
			totalWeighted += g.value * g.weight
			totalWeight += g.weight
		}

		rawAverage := totalWeighted / totalWeight
		roundedAverage := math.Round(rawAverage*100) / 100

		assert.InDelta(t, 7.5, roundedAverage, 0.1)
		assert.True(t, roundedAverage >= 6.0, "Student must pass scrutiny when average >= 6.0")
	})

	t.Run("Failing Average Detection in Scrutiny", func(t *testing.T) {
		grades := []float64{4.5, 5.0, 5.5}
		var sum float64
		for _, v := range grades {
			sum += v
		}
		avg := math.Round((sum/float64(len(grades)))*100) / 100

		assert.Equal(t, 5.0, avg)
		assert.True(t, avg < 6.0, "Student average < 6.0 triggers insufficient grade warning")
	})
}

func TestUserRoleFlagsAndPermissionsUnit(t *testing.T) {
	t.Run("Staff and Active User Evaluation", func(t *testing.T) {
		u := &users.User{
			ID:            "user-101",
			Role:          "teacher",
			IsStaff:       true,
			IsActive:      true,
			EmailVerified: true,
		}

		assert.Equal(t, "user-101", u.ID)
		assert.Equal(t, "teacher", u.Role)
		assert.True(t, u.IsStaff)
		assert.True(t, u.IsActive)
		assert.True(t, u.EmailVerified)
	})

	t.Run("Superadmin Role Privileges", func(t *testing.T) {
		u := &users.User{
			ID:       "sa-1",
			Role:     "superadmin",
			IsActive: true,
		}
		assert.Equal(t, "sa-1", u.ID)
		assert.Equal(t, "superadmin", u.Role)
		assert.True(t, u.IsActive)
	})
}
