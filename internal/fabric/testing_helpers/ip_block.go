package testinghelpers

import (
	"fmt"
	"math/rand/v2"
	"net/netip"
)

// IpBlockTestPool is the range reserved for customer-owned IP blocks created by acceptance tests.
// The sweeper deletes only blocks whose prefix lies within it.
var IpBlockTestPool = netip.MustParsePrefix("198.18.0.0/15")

// RandomIpBlockPrefix returns a random /24 from IpBlockTestPool.
func RandomIpBlockPrefix() string {
	return fmt.Sprintf("198.%d.%d.0/24", 18+rand.IntN(2), rand.IntN(256))
}

// IsIpBlockTestPrefix reports whether prefix lies within IpBlockTestPool.
func IsIpBlockTestPrefix(prefix string) bool {
	p, err := netip.ParsePrefix(prefix)
	if err != nil {
		return false
	}
	return IpBlockTestPool.Contains(p.Addr()) && p.Bits() >= IpBlockTestPool.Bits()
}
