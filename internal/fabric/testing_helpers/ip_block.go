package testinghelpers

import (
	"fmt"
	"math/rand/v2"
	"net/netip"
)

var IpBlockTestPool = netip.MustParsePrefix("198.18.0.0/15")

func RandomIpBlockPrefix() string {
	return fmt.Sprintf("198.%d.%d.0/24", 18+rand.IntN(2), rand.IntN(256))
}

func IsIpBlockTestPrefix(prefix string) bool {
	p, err := netip.ParsePrefix(prefix)
	if err != nil {
		return false
	}
	return IpBlockTestPool.Contains(p.Addr()) && p.Bits() >= IpBlockTestPool.Bits()
}
