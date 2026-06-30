package routeaggregation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var resourceConfig = `
resource "equinix_fabric_route_aggregation" "test" {
  type        = "BGP_IPv4_PREFIX_AGGREGATION"
  name        = var.name
  description = "Test Route Aggregation"
  project = {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
}

variable "name" {
  type = string
}
`

func TestAccFabricRouteAggregation_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				ConfigVariables: tfconfig.Variables{
					"name": tfconfig.StringVariable("route_agg_PFCR"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("state"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("href"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("name"), knownvalue.StringExact("route_agg_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("description"), knownvalue.StringExact("Test Route Aggregation")),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("type"), knownvalue.StringExact("BGP_IPv4_PREFIX_AGGREGATION")),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("project"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
						}),
					),
				},
			},
			{
				Config: resourceConfig,
				ConfigVariables: tfconfig.Variables{
					"name": tfconfig.StringVariable("route_agg_up_PFCR"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("state"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("href"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("name"), knownvalue.StringExact("route_agg_up_PFCR")),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("description"), knownvalue.StringExact("Test Route Aggregation")),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("type"), knownvalue.StringExact("BGP_IPv4_PREFIX_AGGREGATION")),
					statecheck.ExpectKnownValue("equinix_fabric_route_aggregation.test", tfjsonpath.New("project"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
						}),
					),
				},
			},
		},
	})

}

func CheckRouteAggregationDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_route_aggregation" {
			continue
		}

		if routeAggregation, _, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, rs.Primary.ID).Execute(); err == nil {
			if routeAggregation.GetState() == fabricv4.ROUTEAGGREGATIONSTATE_PROVISIONED {
				return fmt.Errorf("fabric route aggregation %s still exists and is %s",
					rs.Primary.ID, routeAggregation.GetState())
			}
		}
	}
	return nil
}
