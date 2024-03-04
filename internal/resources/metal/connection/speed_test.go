package connection

import (
	"testing"
)

func TestSpeedConversion(t *testing.T) {
	speedUint, err := speedStrToUint("50Mbps")
	if err != nil {
		t.Errorf("Error converting speed string to uint64: %s", err)
	}
	if speedUint != 50*mega {
		t.Errorf("Speed string conversion failed. Expected: %d, got: %d", 50*mega, speedUint)
	}

	speedStr, err := speedUintToStr(50 * mega)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "50Mbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "50Mbps", speedStr)
	}

	speedUint, err = speedStrToUint("100Gbps")
	if err != nil {
		t.Errorf("Error converting speed string to uint64: %s", err)
	}
	if speedUint != 100*giga {
		t.Errorf("Speed string conversion failed. Expected: %d, got: %d", 100*giga, speedUint)
	}

	speedStr, err = speedUintToStr(100 * giga)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "100Gbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "100Gbps", speedStr)
	}

	speedStr, err = speedUintToStr(100*giga + 2)
	if err == nil {
		t.Errorf("Expected error converting invalid speed uint to string, got: %s", speedStr)
	}

	speedUint, err = speedStrToUint("100kWh")
	if err == nil {
		t.Errorf("Expected error converting invalid speed string to uint, got: %d", speedUint)
	}
}
