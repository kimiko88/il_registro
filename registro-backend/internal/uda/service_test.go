package uda

import (
	"testing"
)

func TestUdaPlanModel(t *testing.T) {
	plan := UdaPlan{
		ID:           "uda-1",
		Title:        "Unità di Apprendimento 1: La Rivoluzione Industriale",
		Period:       "primo_quadrimestre",
		Competencies: []string{"COMP_CIVIC", "COMP_DIGITAL"},
		Status:       "draft",
	}

	if plan.Title == "" {
		t.Fatalf("title must not be empty")
	}
	if len(plan.Competencies) != 2 {
		t.Fatalf("expected 2 competencies, got %d", len(plan.Competencies))
	}
}
