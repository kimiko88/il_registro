package unit

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"registro-backend/internal/grades"
	"registro-backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func TestGetLimiter_RaceConditionSafety(t *testing.T) {
	limiter := middleware.NewIPRateLimiter(rate.Every(time.Second), 10)
	defer limiter.Close()

	var wg sync.WaitGroup
	ips := []string{"192.168.1.1", "192.168.1.2", "192.168.1.3", "10.0.0.1"}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ip := ips[id%len(ips)]
			l := limiter.GetLimiter(ip)
			if l == nil {
				t.Errorf("got nil rate limiter for IP %s", ip)
			}
		}(i)
	}

	wg.Wait()
}

func TestSecurityHeaders_CSPNoUnsafeInlineAndSanitizedBackendURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("BACKEND_URL", "http://api.example.com; script-src 'unsafe-inline'")
	os.Setenv("BACKEND_WS_URL", "ws://api.example.com/ws\n")
	defer func() {
		os.Unsetenv("BACKEND_URL")
		os.Unsetenv("BACKEND_WS_URL")
	}()

	r := gin.New()
	r.Use(middleware.SecurityHeadersMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	csp := w.Header().Get("Content-Security-Policy")

	// 1. Verify 'unsafe-inline' is NOT present in style-src
	if len(csp) == 0 {
		t.Fatal("Content-Security-Policy header is missing")
	}
	if containsSubstring(csp, "'unsafe-inline'") {
		t.Errorf("Content-Security-Policy should NOT contain 'unsafe-inline', got: %s", csp)
	}

	// 2. Verify BACKEND_URL was sanitized (no semicolon, no newline)
	if containsSubstring(csp, "; script-src 'unsafe-inline'") {
		t.Errorf("Content-Security-Policy contains unsanitized env injection: %s", csp)
	}
}

func TestCalculateAverage_ExcludesGradeZero(t *testing.T) {
	calc := grades.NewCalculator()

	gradeList := []grades.Grade{
		{ID: "g1", GradeValue: 0.0, GradeType: grades.GradeTypeNumeric},
		{ID: "g2", GradeValue: 6.0, GradeType: grades.GradeTypeNumeric},
		{ID: "g3", GradeValue: 9.0, GradeType: grades.GradeTypeNumeric},
	}

	// 0.0 is unrated/unset and excluded; average is (6 + 9) / 2 = 7.5
	avg := calc.CalculateAverage(gradeList)
	if avg != 7.5 {
		t.Errorf("expected average to be 7.5 excluding grade 0.0, got %f", avg)
	}
}

func TestCalculateWeightedAverage_ExcludesGradeZero(t *testing.T) {
	calc := grades.NewCalculator()

	gradeList := []grades.Grade{
		{ID: "g1", GradeValue: 0.0, Weight: 1.0, GradeType: grades.GradeTypeNumeric},
		{ID: "g2", GradeValue: 10.0, Weight: 1.0, GradeType: grades.GradeTypeNumeric},
	}

	// 0.0 is unrated/unset and excluded; weighted average is 10.0
	avg := calc.CalculateWeightedAverage(gradeList)
	if avg != 10.0 {
		t.Errorf("expected weighted average to be 10.0 excluding grade 0.0, got %f", avg)
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && func() bool {
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	}()
}
