package lb

import "testing"

func TestRoundRobin(t *testing.T) {
	r := NewRoundRobin([]string{"a", "b", "c"})
	seen := map[string]bool{}
	for range 6 {
		s := r.Next()
		seen[s] = true
	}
	if len(seen) != 3 {
		t.Fatalf("expected 3 unique, got %d", len(seen))
	}
}
