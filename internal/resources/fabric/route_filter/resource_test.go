package route_filter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("name"), knownvalue.StringExact("RF_Policy_X_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("description"), knownvalue.StringExact("Route Filter Policy for X Purpose")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("type"), knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("state"), knownvalue.StringExact(string(fabricv4.ROUTEFILTERSTATE_PROVISIONED))),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("not_matched_rule_action"), knownvalue.StringExact("DENY")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("rules_count"), knownvalue.Int32Exact(0)),
				},
				ExpectNonEmptyPlan: false,
			},
			{
				Config: resourceConfig,
				ConfigVariables: config.Variables{
					"name":        config.StringVariable("RF_Policy_Y_PFCR"),
					"description": config.StringVariable("Route Filter Policy for Y Purpose"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("name"), knownvalue.StringExact("RF_Policy_Y_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("description"), knownvalue.StringExact("Route Filter Policy for Y Purpose")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("type"), knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("state"), knownvalue.StringExact(string(fabricv4.ROUTEFILTERSTATE_REPROVISIONING))),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("not_matched_rule_action"), knownvalue.StringExact("DENY")),
					statecheck.ExpectKnownValue("equinix_fabric_route_filter.test", tfjsonpath.New("rules_count"), knownvalue.Int32Exact(0)),
				},
				ExpectNonEmptyPlan: false,
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
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
