package route_filter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var resourceConfig = `
resource "equinix_fabric_route_filter" "test" {
  name = var.name
  project {
    project_id = "291639000636552"
  }
  type        = "BGP_IPv4_PREFIX_FILTER"
  description = var.description
}

variable "name" {
  type = string
}

variable "description" {
  type = string
}
`

func TestAccFabricRouteFilterPolicy_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				ConfigVariables: config.Variables{
					"name":        config.StringVariable("RF_Policy_X_PFCR"),
					"description": config.StringVariable("Route Filter Policy for X Purpose"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_route_filter.test", map[string]knownvalue.Check{
						"id":                      knownvalue.NotNull(),
						"name":                    knownvalue.StringExact("RF_Policy_X_PFCR"),
						"description":             knownvalue.StringExact("Route Filter Policy for X Purpose"),
						"type":                    knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER"),
						"state":                   knownvalue.StringExact(string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)),
						"not_matched_rule_action": knownvalue.StringExact("DENY"),
						"rules_count":             knownvalue.Int32Exact(0),
					}),
				},
			},
		},
	})

}

func CheckRouteFilterDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_route_filter" {
			continue
		}

		err := route_filter.WaitForDeletion(ctx, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for route filter deletion. ID: %s, Err: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
