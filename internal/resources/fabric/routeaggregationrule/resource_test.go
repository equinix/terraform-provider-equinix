package routeaggregationrule_test

import (
	"context"
	"fmt"
	"testing"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccFabricRouteAggregationRuleConfig(prefix string) string {
	return fmt.Sprintf(`
		
		resource "equinix_fabric_route_aggregation" "test" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "test-aggregation"
		  description = "Test Route Aggregation"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}

		resource "equinix_fabric_route_aggregation_rule" "test" {
			route_aggregation_id = equinix_fabric_route_aggregation.test.id
			name = "RouteAggregationRulePFCR"
  			description = "Test aggregation rule"
  			prefix = "%s"
		}
	`, prefix)
}

func TestAccFabricRouteAggregationRule_PNFV(t *testing.T) {
	routeAggregationPrefix := "192.169.0.0/24"
	upRouteAggregationPrefix := "192.168.0.0/24"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteAggregationRuleConfig(routeAggregationPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_aggregation_rule.test", "name", "RouteAggregationRulePFCR"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "state"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "prefix", "192.169.0.0/24"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "type", "BGP_IPv4_PREFIX_AGGREGATION_RULE"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "description", "Test aggregation rule"),
				),
			},
			{
				Config: testAccFabricRouteAggregationRuleConfig(upRouteAggregationPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_aggregation_rule.test", "name", "RouteAggregationRulePFCR"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "state"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "prefix", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "type", "BGP_IPv4_PREFIX_AGGREGATION_RULE"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "description", "Test aggregation rule"),
				),
			},
		},
	})

}

func CheckRouteAggregationRuleDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_route_aggregation_rule" {
			continue
		}

		routeAggregationId := rs.Primary.Attributes["route_aggregation_id"]
		routeAggregationRule, _, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationId, rs.Primary.ID).Execute()
		if err != nil {
			if genericError, ok := err.(*fabricv4.GenericOpenAPIError); ok {
				if fabricErrs, ok := genericError.Model().([]fabricv4.Error); ok {
					if equinix_errors.HasErrorCode(fabricErrs, "EQ-3044402") {
						return nil
					}
				}
			}
			if routeAggregationRule.GetState() == fabricv4.ROUTEAGGREGATIONRULESTATE_PROVISIONED {
				return fmt.Errorf("fabric stream %s still exists and is %s",
					rs.Primary.ID, routeAggregationRule.GetState())
			}

			return fmt.Errorf("unexpected API error checking deletion: %v", err)
		}
	}
	return nil
}
