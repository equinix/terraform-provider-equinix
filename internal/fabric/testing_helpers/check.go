package testinghelpers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func ExpectKnownAttributes(address string, checks map[string]knownvalue.Check) expectKnownAttributes {
	return expectKnownAttributes{
		address: address,
		checks:  checks,
	}
}

var _ statecheck.StateCheck = expectKnownAttributes{}

type expectKnownAttributes struct {
	address string
	checks  map[string]knownvalue.Check
}

func (e expectKnownAttributes) CheckState(ctx context.Context, req statecheck.CheckStateRequest, res *statecheck.CheckStateResponse) {
	for path, check := range e.checks {
		statecheck.ExpectKnownValue(e.address, tfjsonpath.New(path), check).CheckState(ctx, req, res)
	}
}
