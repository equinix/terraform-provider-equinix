package connection

import (
	"testing"
)

func TestSpeedConversion(t *testing.T) {
	speedStr, err := speedIntToStr(50 * mega)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "50Mbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "50Mbps", speedStr)
	}

	speedStr, err = speedIntToStr(100 * giga)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "100Gbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "100Gbps", speedStr)
	}

	speedStr, err = speedIntToStr(100*giga + 2)
	if err == nil {
		t.Errorf("Expected error converting invalid speed uint to string, got: %s", speedStr)
	}
}
