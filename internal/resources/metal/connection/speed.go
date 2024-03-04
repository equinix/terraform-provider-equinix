package connection

import (
	"fmt"
	"strings"
)

var (
	mega          uint64 = 1000 * 1000
	giga          uint64 = 1000 * mega
	allowedSpeeds        = []struct {
		Int uint64
		Str string
	}{
		{50 * mega, "50Mbps"},
		{200 * mega, "200Mbps"},
		{500 * mega, "500Mbps"},
		{1 * giga, "1Gbps"},
		{2 * giga, "2Gbps"},
		{5 * giga, "5Gbps"},
		{10 * giga, "10Gbps"},
		{100 * giga, "100Gbps"},
	}
)

func allowedSpeedsString() string {
	allowedStrings := []string{}
	for _, allowedSpeed := range allowedSpeeds {
		allowedStrings = append(allowedStrings, allowedSpeed.Str)
	}
	return strings.Join(allowedStrings, ", ")
}

func speedStrToUint(speed string) (uint64, error) {
	allowedStrings := []string{}
	for _, allowedSpeed := range allowedSpeeds {
		if allowedSpeed.Str == speed {
			return allowedSpeed.Int, nil
		}
		allowedStrings = append(allowedStrings, allowedSpeed.Str)
	}
	return 0, fmt.Errorf("invalid speed string: %s. Allowed strings: %s", speed, strings.Join(allowedStrings, ", "))
}

func speedUintToStr(speed uint64) (string, error) {
	allowedUints := []uint64{}
	for _, allowedSpeed := range allowedSpeeds {
		if speed == allowedSpeed.Int {
			return allowedSpeed.Str, nil
		}
		allowedUints = append(allowedUints, allowedSpeed.Int)
	}
	return "", fmt.Errorf("%d is not allowed speed value. Allowed values: %v", speed, allowedUints)
}
