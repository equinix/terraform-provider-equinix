package validation

import "github.com/equinix/ecx-go/v2"

func IsEmpty(v interface{}) bool {
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
