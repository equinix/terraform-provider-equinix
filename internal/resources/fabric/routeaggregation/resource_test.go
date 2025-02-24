package routeaggregation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccFabricRouteAggregationConfig(name string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_aggregation" "test" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "%s"
		  description = "Test Route Aggregation"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}
	`, name)
}

func TestAccFabricRouteAggregation_PFCR(t *testing.T) {
	routeAggregationName := "stream_PFCR"
	upRouteAggregationName := "stream_up_PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteAggregationConfig(routeAggregationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_aggregation.test", "name", routeAggregationName),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "state"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "project.project_id"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "name", routeAggregationName),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "type", "BGP_IPv4_PREFIX_AGGREGATION"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "description", "Test Route Aggregation"),
				),
			},
			{
				Config: testAccFabricRouteAggregationConfig(upRouteAggregationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_aggregation.test", "name", upRouteAggregationName),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "state"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "project.project_id"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "name", upRouteAggregationName),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "type", "BGP_IPv4_PREFIX_AGGREGATION"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "description", "Test Route Aggregation"),
				),
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
