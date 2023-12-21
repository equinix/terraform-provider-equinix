package metal_connection

import "testing"

func TestSpeedConversion(t *testing.T) {
	speedUint, err := SpeedStrToUint("50Mbps")
	if err != nil {
		t.Errorf("Error converting speed string to uint64: %s", err)
	}
	if speedUint != 50*Mega {
		t.Errorf("Speed string conversion failed. Expected: %d, got: %d", 50*Mega, speedUint)
	}

	speedStr, err := SpeedUintToStr(50 * Mega)
	if err != nil {
		t.Errorf("Error converting speed uint to string: %s", err)
	}
	if speedStr != "50Mbps" {
		t.Errorf("Speed uint conversion failed. Expected: %s, got: %s", "50Mbps", speedStr)
	}

	speedUint, err = SpeedStrToUint("100Gbps")
	if err == nil {
		t.Errorf("Expected error converting invalid speed string to uint, got: %d", speedUint)
	}

	speedStr, err = SpeedUintToStr(100 * Giga)
	if err == nil {
		t.Errorf("Expected error converting invalid speed uint to string, got: %s", speedStr)
	}
}
