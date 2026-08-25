package textbooks

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateTextbook_CrossTenant_Blocked(t *testing.T) {
	repo := newMockRepo()
	repo.textbooks["tb-1"] = &Textbook{
		ID:       "tb-1",
		SchoolID: "school-1",
		Title:    "Libro Scuola 1",
		Price:    20.0,
	}

	svc := NewService(repo)

	// Teacher from school-2 attempts to update textbook in school-1
	err := svc.UpdateTextbook(context.Background(), "teacher", "school-2", "tb-1", CreateTextbookRequest{
		Title: "Libro Modificato",
		Price: 25.0,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestDeleteTextbook_CrossTenant_Blocked(t *testing.T) {
	repo := newMockRepo()
	repo.textbooks["tb-1"] = &Textbook{
		ID:       "tb-1",
		SchoolID: "school-1",
		Title:    "Libro Scuola 1",
		Price:    20.0,
	}

	svc := NewService(repo)

	// Teacher from school-2 attempts to delete textbook in school-1
	err := svc.DeleteTextbook(context.Background(), "teacher", "school-2", "tb-1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
	assert.Len(t, repo.textbooks, 1)
}

func TestDeleteTextbook_SameSchool_Allowed(t *testing.T) {
	repo := newMockRepo()
	repo.textbooks["tb-1"] = &Textbook{
		ID:       "tb-1",
		SchoolID: "school-1",
		Title:    "Libro Scuola 1",
		Price:    20.0,
	}

	svc := NewService(repo)

	// Teacher from school-1 deletes textbook in school-1
	err := svc.DeleteTextbook(context.Background(), "teacher", "school-1", "tb-1")
	assert.NoError(t, err)
	assert.Len(t, repo.textbooks, 0)
}
