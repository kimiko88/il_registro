package metrics

import (
	"bytes"
	"database/sql"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Registry struct {
	mu           sync.RWMutex
	requestCount map[string]*uint64
	latencies    map[string]*uint64
}

var DefaultRegistry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		requestCount: make(map[string]*uint64),
		latencies:    make(map[string]*uint64),
	}
}

func (r *Registry) RecordRequest(method, path string, status int, duration time.Duration) {
	normPath := normalizePath(path)
	key := fmt.Sprintf(`method="%s",path="%s",status="%d"`, method, normPath, status)

	r.mu.RLock()
	cnt, ok := r.requestCount[key]
	r.mu.RUnlock()

	if !ok {
		r.mu.Lock()
		if cnt, ok = r.requestCount[key]; !ok {
			var newCnt uint64
			var newLat uint64
			r.requestCount[key] = &newCnt
			r.latencies[key] = &newLat
			cnt = &newCnt
		}
		r.mu.Unlock()
	}

	atomic.AddUint64(cnt, 1)
	r.mu.RLock()
	lat := r.latencies[key]
	r.mu.RUnlock()
	if lat != nil {
		atomic.AddUint64(lat, uint64(duration.Milliseconds()))
	}
}

// normalizePath replaces UUID segments and numeric IDs with :id to prevent high cardinality
func normalizePath(p string) string {
	parts := strings.Split(p, "/")
	for i, part := range parts {
		if len(part) == 36 && strings.Count(part, "-") == 4 {
			parts[i] = ":id"
		} else if isNumeric(part) {
			parts[i] = ":id"
		}
	}
	return strings.Join(parts, "/")
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func (r *Registry) GeneratePrometheus(db *sql.DB) string {
	var buf bytes.Buffer

	// 1. Go Runtime Stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	buf.WriteString("# HELP go_goroutines Number of goroutines currently existing.\n")
	buf.WriteString("# TYPE go_goroutines gauge\n")
	fmt.Fprintf(&buf, "go_goroutines %d\n", runtime.NumGoroutine())

	buf.WriteString("# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.\n")
	buf.WriteString("# TYPE go_memstats_alloc_bytes gauge\n")
	fmt.Fprintf(&buf, "go_memstats_alloc_bytes %d\n", m.Alloc)

	buf.WriteString("# HELP go_memstats_sys_bytes Total bytes of memory obtained from the OS.\n")
	buf.WriteString("# TYPE go_memstats_sys_bytes gauge\n")
	fmt.Fprintf(&buf, "go_memstats_sys_bytes %d\n", m.Sys)

	// 2. Database Stats
	if db != nil {
		stats := db.Stats()
		buf.WriteString("# HELP db_open_connections The number of established connections both in use and idle.\n")
		buf.WriteString("# TYPE db_open_connections gauge\n")
		fmt.Fprintf(&buf, "db_open_connections %d\n", stats.OpenConnections)

		buf.WriteString("# HELP db_in_use_connections The number of connections currently in use.\n")
		buf.WriteString("# TYPE db_in_use_connections gauge\n")
		fmt.Fprintf(&buf, "db_in_use_connections %d\n", stats.InUse)

		buf.WriteString("# HELP db_idle_connections The number of idle connections.\n")
		buf.WriteString("# TYPE db_idle_connections gauge\n")
		fmt.Fprintf(&buf, "db_idle_connections %d\n", stats.Idle)

		buf.WriteString("# HELP db_wait_count The total number of connections waited for.\n")
		buf.WriteString("# TYPE db_wait_count counter\n")
		fmt.Fprintf(&buf, "db_wait_count %d\n", stats.WaitCount)
	}

	// 3. HTTP Request Metrics
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.requestCount) > 0 {
		buf.WriteString("# HELP http_requests_total Total number of HTTP requests made.\n")
		buf.WriteString("# TYPE http_requests_total counter\n")
		for label, cnt := range r.requestCount {
			fmt.Fprintf(&buf, "http_requests_total{%s} %d\n", label, atomic.LoadUint64(cnt))
		}

		buf.WriteString("# HELP http_request_duration_ms_total Total duration of HTTP requests in milliseconds.\n")
		buf.WriteString("# TYPE http_request_duration_ms_total counter\n")
		for label, lat := range r.latencies {
			fmt.Fprintf(&buf, "http_request_duration_ms_total{%s} %d\n", label, atomic.LoadUint64(lat))
		}
	}

	return buf.String()
}
