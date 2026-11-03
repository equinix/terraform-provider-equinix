package internetaccess

import (
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
)

func TestIsTestService(t *testing.T) {
	cases := map[string]bool{
		"eia_tf_acc_PFCR": true,
		"eia_tf_acc_PNFV": true,
		"eia_tf_acc_PPDS": true,
		"production-eia":  false,
		"PFCR_eia":        false,
		"":                false,
	}
	for name, want := range cases {
		svc := fabricv4.InternetAccessService{Name: name}
		if got := isTestService(&svc); got != want {
			t.Errorf("isTestService(%q) = %v, want %v", name, got, want)
		}
	}
}

func TestSweepableStates(t *testing.T) {
	for _, state := range sweepableStates() {
		if state == string(fabricv4.INTERNETACCESSSERVICESTATE_DEPROVISIONED) || state == string(fabricv4.INTERNETACCESSSERVICESTATE_FAILED) {
			t.Errorf("services in state %s must not be swept", state)
		}
	}
}

func TestCountKept(t *testing.T) {
	services := []fabricv4.InternetAccessService{{Name: "a_PFCR"}, {Name: "mine"}, {Name: "b_PFCR"}, {Name: "yours"}}
	if got := countKept(services); got != 2 {
		t.Errorf("countKept = %d, want 2", got)
	}
}
