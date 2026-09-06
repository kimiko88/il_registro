package metrics

import (
	"strings"
	"testing"
	"time"
)

func TestMetrics_RecordAndGenerate(t *testing.T) {
	reg := NewRegistry()

	reg.RecordRequest("GET", "/api/v1/classes/123e4567-e89b-12d3-a456-426614174000/students", 200, 25*time.Millisecond)
	reg.RecordRequest("GET", "/api/v1/classes/123e4567-e89b-12d3-a456-426614174000/students", 200, 15*time.Millisecond)
	reg.RecordRequest("POST", "/api/v1/grades", 201, 50*time.Millisecond)
	reg.RecordRequest("GET", "/api/v1/users/42", 404, 5*time.Millisecond)

	out := reg.GeneratePrometheus(nil)

	// Check runtime gauges
	if !strings.Contains(out, "go_goroutines") {
		t.Errorf("expected go_goroutines in output, got:\n%s", out)
	}
	if !strings.Contains(out, "go_memstats_alloc_bytes") {
		t.Errorf("expected go_memstats_alloc_bytes in output, got:\n%s", out)
	}

	// Check HTTP request counters with normalized path
	if !strings.Contains(out, `http_requests_total{method="GET",path="/api/v1/classes/:id/students",status="200"} 2`) {
		t.Errorf("expected normalized class path count 2, got:\n%s", out)
	}
	if !strings.Contains(out, `http_requests_total{method="POST",path="/api/v1/grades",status="201"} 1`) {
		t.Errorf("expected post grades count 1, got:\n%s", out)
	}
	if !strings.Contains(out, `http_requests_total{method="GET",path="/api/v1/users/:id",status="404"} 1`) {
		t.Errorf("expected normalized users path count 1, got:\n%s", out)
	}
}
