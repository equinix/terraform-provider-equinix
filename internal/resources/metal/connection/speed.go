package connection

import (
	"fmt"
	"regexp"
	"strconv"
)

var (
	mega        int64 = 1000 * 1000
	giga        int64 = 1000 * mega
	SpeedFormat       = regexp.MustCompile(`^(\d+)((M|G)bps)$`)
)

func validateSpeedStr(speed string) error {
	if ok := SpeedFormat.Match([]byte(speed)); !ok {
		return fmt.Errorf("invalid speed string %v, must match %v", speed, SpeedFormat.String())
	}
	return nil
}

func speedIntToStr(speed int64) (string, error) {
	var base int64
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
