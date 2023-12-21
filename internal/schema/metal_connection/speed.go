package metal_connection

import (
	"fmt"
	"strings"
)

var (
	Mega          uint64 = 1000 * 1000
	Giga          uint64 = 1000 * Mega
	AllowedSpeeds        = []struct {
		Int uint64
		Str string
	}{
		{50 * Mega, "50Mbps"},
		{200 * Mega, "200Mbps"},
		{500 * Mega, "500Mbps"},
		{1 * Giga, "1Gbps"},
		{2 * Giga, "2Gbps"},
		{5 * Giga, "5Gbps"},
		{10 * Giga, "10Gbps"},
	}
)

func SpeedStrToUint(speed string) (uint64, error) {
	allowedStrings := []string{}
	for _, allowedSpeed := range AllowedSpeeds {
		if allowedSpeed.Str == speed {
			return allowedSpeed.Int, nil
		}
		allowedStrings = append(allowedStrings, allowedSpeed.Str)
	}
	return 0, fmt.Errorf("invalid speed string: %s. Allowed strings: %s", speed, strings.Join(allowedStrings, ", "))
}

func SpeedUintToStr(speed uint64) (string, error) {
	allowedUints := []uint64{}
	for _, allowedSpeed := range AllowedSpeeds {
		if speed == allowedSpeed.Int {
			return allowedSpeed.Str, nil
		}
		allowedUints = append(allowedUints, allowedSpeed.Int)
	}
	return "", fmt.Errorf("%d is not allowed speed value. Allowed values: %v", speed, allowedUints)
}
