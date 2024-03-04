package connection

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	mega        uint64 = 1000 * 1000
	giga        uint64 = 1000 * mega
	SpeedFormat        = regexp.MustCompile(`^(\d+)((M|G)bps)$`)
)

func speedStrToUint(speed string) (uint64, error) {
	parts := SpeedFormat.FindStringSubmatch(speed)
	if parts != nil {
		base, err := strconv.Atoi(parts[1])
		if parts[2] == "Mbps" {
			return uint64(base) * mega, nil
		} else if parts[2] == "Gbps" {
			return uint64(base) * giga, nil
		}
		return 0, err
	}
	return 0, fmt.Errorf("invalid speed string %v, must match %v", speed, SpeedFormat.String())
}

func speedUintToStr(speed uint64) (string, error) {
	var base uint64
	var unit string

	if (speed % giga) == 0 {
		unit = "Gbps"
		base = speed / giga
	} else if (speed % mega) == 0 {
		unit = "Mbps"
		base = speed / mega
	} else {
		return "", fmt.Errorf("unsupported speed value %v", speed)
	}
	return strconv.Itoa(int(base)) + unit, nil
}
