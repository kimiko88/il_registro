package pcto

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegression_CreateProject_Validation(t *testing.T) {
	repo := new(MockRepo)
	svc := NewService(repo)

	t.Run("Invalid Date Format", func(t *testing.T) {
		req := CreateProjectRequest{
			Title: "Test", Type: "Internal",
			StartDate: "invalid", EndDate: "2025-01-01",
		}
		err := svc.CreateProject(context.Background(), "school1", "secretary", "t1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid start_date")
	})

	t.Run("End Before Start", func(t *testing.T) {
		req := CreateProjectRequest{
			Title: "Test", Type: "Internal",
			StartDate: "2025-02-01", EndDate: "2025-01-01",
		}
		err := svc.CreateProject(context.Background(), "school1", "secretary", "t1", req)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "end_date cannot be before")
	})
}
