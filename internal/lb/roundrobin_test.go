package lb

import (
	"testing"
)

func TestRoundRobin_Distribution(t *testing.T) {
	r := NewRoundRobin([]string{"a", "b", "c"})
	seen := map[string]int{}
	for range 6 {
		s := r.Next()
		seen[s]++
	}
	if len(seen) != 3 {
		t.Fatalf("expected 3 unique backends, got %d", len(seen))
	}
	for _, count := range seen {
		if count != 2 {
			t.Errorf("expected each backend 2 times in 6 calls, got %v", seen)
			break
		}
	}
}

func TestRoundRobin_EmptyReturnsEmpty(t *testing.T) {
	r := NewRoundRobin(nil)
	if s := r.Next(); s != "" {
		t.Errorf("Next() = %q, want \"\"", s)
	}
	r2 := NewRoundRobin([]string{})
	if s := r2.Next(); s != "" {
		t.Errorf("Next() = %q, want \"\"", s)
	}
}

func TestRoundRobin_SingleBackend(t *testing.T) {
	r := NewRoundRobin([]string{"only"})
	for range 3 {
		if s := r.Next(); s != "only" {
			t.Errorf("Next() = %q, want \"only\"", s)
		}
	}
}
