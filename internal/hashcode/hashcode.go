package hashcode

import (
	"hash/crc32"
)

// String hashes a string to a unique hashcode.
//
// crc32 returns a uint32, but for our use we need
// and non negative integer. Here we cast to an integer
// and invert it if the result is negative.
//
// Originally from https://github.com/hashicorp/terraform-plugin-sdk/blob/main/internal/helper/hashcode/hashcode.go
func String(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
