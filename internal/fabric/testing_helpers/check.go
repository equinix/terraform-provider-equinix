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

func ExpectKnownAttributesAt(address string, path tfjsonpath.Path, checks map[string]knownvalue.Check) expectKnownAttributes {
	return expectKnownAttributes{
		address:  address,
		checks:   checks,
		basePath: &path,
	}
}

var _ statecheck.StateCheck = expectKnownAttributes{}

type expectKnownAttributes struct {
	address  string
	checks   map[string]knownvalue.Check
	basePath *tfjsonpath.Path
}

func (e expectKnownAttributes) CheckState(ctx context.Context, req statecheck.CheckStateRequest, res *statecheck.CheckStateResponse) {
	for key, check := range e.checks {
		path := func() tfjsonpath.Path {
			if e.basePath == nil {
				return tfjsonpath.New(key)
			}

			return e.basePath.AtMapKey(key)
		}()

		statecheck.ExpectKnownValue(e.address, path, check).CheckState(ctx, req, res)
	}
}
