package acceptance

import (
	"strings"
	"testing"
	"time"
)

func TestDeviceTerminationTime(t *testing.T) {
	expectedTime := time.Now().UTC().Add(60 * time.Minute)
	result := DeviceTerminationTime()

	parsedTime, err := time.Parse(time.RFC3339, result)
	if err != nil {
		t.Fatalf("Failed to parse time: %v", err)
	}
	drift := parsedTime.Sub(expectedTime)

	if drift < 0 {
		drift = -drift // get absolute value
	}

	if drift > time.Second {
		t.Errorf("Expected time close to %v, but got %v", expectedTime, parsedTime)
	}
}

func TestConfAccMetalDeviceBase(t *testing.T) {
	plans := []string{"c3.small.x86", "m3.large.x86"}
	metros := []string{"ny5", "sv15"}
	os := []string{"ubuntu_22_04"}

	result := ConfAccMetalDeviceBase(plans, metros, os)

	// Check that all plans are included
	for _, plan := range plans {
		if !strings.Contains(result, plan) {
			t.Errorf("Expected plan %q to be in result", plan)
		}
	}

	// Check that all metros are included
	for _, metro := range metros {
		if !strings.Contains(result, metro) {
			t.Errorf("Expected metro %q to be in result", metro)
		}
	}

	// Check that the OS is included
	if !strings.Contains(result, os[0]) {
		t.Errorf("Expected OS %q to be in result", os[0])
	}

	// Check for correct quoting
	if !strings.Contains(result, `"c3.small.x86"`) {
		t.Errorf("Expected quoted plan name in result")
	}
}
