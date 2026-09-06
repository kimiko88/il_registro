package cache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryCache_BasicOperations(t *testing.T) {
	mc := NewMemoryCache()
	defer func() { _ = mc.Close() }()

	ctx := context.Background()

	// 1. Get non-existent key
	_, err := mc.Get(ctx, "non_existent")
	if err != ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss, got %v", err)
	}

	// 2. Set and Get
	err = mc.Set(ctx, "school:1", `{"name":"Liceo Galilei"}`, 1*time.Minute)
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	val, err := mc.Get(ctx, "school:1")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if val != `{"name":"Liceo Galilei"}` {
		t.Fatalf("unexpected value: %s", val)
	}

	// 3. Delete
	err = mc.Delete(ctx, "school:1")
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = mc.Get(ctx, "school:1")
	if err != ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after delete, got %v", err)
	}
}

func TestMemoryCache_TTLExpiration(t *testing.T) {
	mc := NewMemoryCache()
	defer func() { _ = mc.Close() }()

	ctx := context.Background()

	// Set with very short TTL
	err := mc.Set(ctx, "temp_key", "temporary_val", 50*time.Millisecond)
	if err != nil {
		t.Fatalf("set failed: %v", err)
	}

	// Immediate get
	val, err := mc.Get(ctx, "temp_key")
	if err != nil || val != "temporary_val" {
		t.Fatalf("immediate get failed, got: %s, err: %v", val, err)
	}

	// Wait for expiration
	time.Sleep(70 * time.Millisecond)

	_, err = mc.Get(ctx, "temp_key")
	if err != ErrCacheMiss {
		t.Fatalf("expected ErrCacheMiss after expiration, got %v", err)
	}
}

func TestMemoryCache_DeletePrefix(t *testing.T) {
	mc := NewMemoryCache()
	defer func() { _ = mc.Close() }()

	ctx := context.Background()

	_ = mc.Set(ctx, "class:1:students", "list1", 1*time.Minute)
	_ = mc.Set(ctx, "class:1:teachers", "list2", 1*time.Minute)
	_ = mc.Set(ctx, "class:2:students", "list3", 1*time.Minute)
	_ = mc.Set(ctx, "subject:1", "math", 1*time.Minute)

	err := mc.DeletePrefix(ctx, "class:1:")
	if err != nil {
		t.Fatalf("delete prefix failed: %v", err)
	}

	// class:1 keys should be gone
	if _, err := mc.Get(ctx, "class:1:students"); err != ErrCacheMiss {
		t.Fatalf("expected class:1:students to be deleted")
	}
	if _, err := mc.Get(ctx, "class:1:teachers"); err != ErrCacheMiss {
		t.Fatalf("expected class:1:teachers to be deleted")
	}

	// other keys must remain intact
	if val, err := mc.Get(ctx, "class:2:students"); err != nil || val != "list3" {
		t.Fatalf("class:2:students was unexpectedly modified")
	}
	if val, err := mc.Get(ctx, "subject:1"); err != nil || val != "math" {
		t.Fatalf("subject:1 was unexpectedly modified")
	}
}

func TestMemoryCache_Concurrency(t *testing.T) {
	mc := NewMemoryCache()
	defer func() { _ = mc.Close() }()

	ctx := context.Background()
	var wg sync.WaitGroup

	// Run 20 concurrent readers and writers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", id%5)
			_ = mc.Set(ctx, key, fmt.Sprintf("val_%d", id), 1*time.Minute)
			_, _ = mc.Get(ctx, key)
		}(i)
	}

	wg.Wait()
}
