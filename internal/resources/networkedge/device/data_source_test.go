package equinix

import (
	"reflect"
	"testing"

	"github.com/equinix/ne-go"
)

func TestGetDeviceStatusList(t *testing.T) {
	input := ""
	result, err := getNeDeviceStatusList(input)
	if err != nil {
		t.Errorf("got error %v for input: %v", err, input)
	}

	if len(*result) != 0 {
		t.Errorf("bad %v len: %v", *result, len(*result))
	}

	test := func(input string, expected *[]string, expectError bool) {
		result, err := getNeDeviceStatusList(input)
		if err != nil {
			if expectError {
				return // got an expected error, so good
			}
			t.Errorf("got unexpected error: { %v } with input: %v", err, input)
			return
		}

		if expectError {
			t.Errorf("did not receive expected error with input: %v", input)
			return
		}

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("bad %v !=  %v", *result, *expected)
		}
	}

	input = "provisioned" // single item
	expected := &[]string{ne.DeviceStateProvisioned}
	test(input, expected, false)

	input = "provisioning,provisioned" // 2 items
	expected = &[]string{ne.DeviceStateProvisioning, ne.DeviceStateProvisioned}
	test(input, expected, false)

	input = "provisioning, \tprovisioned  " // whitespace
	expected = &[]string{ne.DeviceStateProvisioning, ne.DeviceStateProvisioned}
	test(input, expected, false)

	input = "provisioning, \tproVISioned  " // capitalization
	expected = &[]string{ne.DeviceStateProvisioning, ne.DeviceStateProvisioned}
	test(input, expected, false)

	input = "provisioning, provisioned,somethingInvalid " // error on invalid entries
	test(input, nil, true)
}
