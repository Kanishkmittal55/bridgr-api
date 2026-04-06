package bridgr_worker

import "testing"

func TestDiscoveryPageMaxResults(t *testing.T) {
	if n := discoveryPageMaxResults(3, 30, 20); n != 8 {
		t.Fatalf("need=3 ceiling=30 cap=20: got %d want 8 (floor)", n)
	}
	if n := discoveryPageMaxResults(25, 30, 20); n != 20 {
		t.Fatalf("need=25: got %d want 20 (page cap)", n)
	}
	if n := discoveryPageMaxResults(1, 0, 0); n >= 8 && n <= 20 {
		// default cap 20, need+4=5 floored to 8
	} else {
		t.Fatalf("unexpected %d", n)
	}
}
