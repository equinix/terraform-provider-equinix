package asn

import (
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

func TestIsTestAsn(t *testing.T) {
	asn := func(asnType string, number int64) *fabricv4.Asn {
		a := &fabricv4.Asn{AdditionalProperties: map[string]interface{}{typeKey: asnType}}
		a.SetNumber(number)
		return a
	}
	tests := map[string]struct {
		asn  *fabricv4.Asn
		want bool
	}{
		"private asn":         {asn("PRIVATE_ASN", 64529), true},
		"customer public asn": {asn("PUBLIC_ASN", 432), false},
		"unknown type":        {asn("", 64529), false},
		"no type at all":      {&fabricv4.Asn{}, false},
	}
	for name, tc := range tests {
		if got := isTestAsn(tc.asn); got != tc.want {
			t.Errorf("%s: isTestAsn = %v, want %v", name, got, tc.want)
		}
	}
}
