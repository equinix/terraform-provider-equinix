package equinix

import "github.com/equinix/ecx-go/v2"

// Deprecated: use slices package instead
func stringsFound(source []string, target []string) bool {
	for i := range source {
		if !isStringInSlice(source[i], target) {
			return false
		}
	}
	return true
}

// Deprecated: use slices package instead
func atLeastOneStringFound(source []string, target []string) bool {
	for i := range source {
		if isStringInSlice(source[i], target) {
			return true
		}
	}
	return false
}

// Deprecated: use slices package instead
func isStringInSlice(needle string, hay []string) bool {
	for i := range hay {
		if needle == hay[i] {
			return true
		}
	}
	return false
}

// Deprecated
func isEmpty(v interface{}) bool {
	switch v := v.(type) {
	case int:
		return v == 0
	case *int:
		return ecx.IntValue(v) == 0
	case string:
		return v == ""
	case *string:
		return ecx.StringValue(v) == ""
	case nil:
		return true
	default:
		return false
	}
}
