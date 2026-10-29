package testinghelpers

import "testing"

func TestRandomIpBlockPrefixInPool(t *testing.T) {
	for range 1000 {
		if p := RandomIpBlockPrefix(); !IsIpBlockTestPrefix(p) {
			t.Fatalf("random prefix %s is outside the test pool", p)
		}
	}
}

func TestIsIpBlockTestPrefix(t *testing.T) {
	for prefix, want := range map[string]bool{
		"198.18.0.0/24":   true,
		"198.19.255.0/24": true,
		"198.20.0.0/24":   false,
		"203.0.113.0/24":  false,
		"198.18.0.0/14":   false,
		"":                false,
	} {
		if got := IsIpBlockTestPrefix(prefix); got != want {
			t.Errorf("IsIpBlockTestPrefix(%q) = %v, want %v", prefix, got, want)
		}
	}
}
