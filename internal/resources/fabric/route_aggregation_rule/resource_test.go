package route_aggregation_rule_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccFabricRouteAggregationRuleConfig(name string) string {
	return fmt.Sprintf(`
		
		resource "equinix_fabric_route_aggregation" "test" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "%s"
		  description = "Test Route Aggregation"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}

		resource "equinix_fabric_route_aggregation_rule" "test" {
			route_aggregation_id = equinix_fabric_route_aggregation.test.id
			name = "%s"
  			description = "Test aggregation rule"
  			prefix = "192.169.0.0/24"
		}
	`, name)
}

func TestAccFabricRouteAggregationRule_PFCR(t *testing.T) {
	routeAggregationRuleName := "RouteAggregationRulePFCR"
	//upRouteAggregationName := "stream_up_PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteAggregationRuleConfig(routeAggregationRuleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_aggregation_rule.test", "name", routeAggregationRuleName),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "state"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation_rule.test", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "name", routeAggregationRuleName),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "type", "BGP_IPv4_PREFIX_AGGREGATION_RULE"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation_rule.test", "description", "Test aggregation rule"),
				),
				ExpectNonEmptyPlan: false,
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
		_, resp, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationId, rs.Primary.ID).Execute()
		if err != nil {
			// Check if the response exists and contains status 400 or 404
			if resp != nil && (resp.StatusCode == 400 || resp.StatusCode == 404) {
				fmt.Printf("Resource %s not found, treating as deleted\n", rs.Primary.ID)
				return nil
			}

			// Handle specific API error messages
			var apiErr *fabricv4.GenericOpenAPIError
			if errors.As(err, &apiErr) {
				errorBody := apiErr.Body()
				var errorResponse map[string]interface{}
				if jsonErr := json.Unmarshal(errorBody, &errorResponse); jsonErr == nil {
					if errorCode, exists := errorResponse["errorCode"]; exists && errorCode == "EQ-3044402" {
						fmt.Printf("Detected EQ-3044402 for resource %s, treating as deleted\n", rs.Primary.ID)
						return nil // Successfully handled the expected deletion case
					}
				}
			}

			return fmt.Errorf("unexpected API error checking deletion: %v", err)
		}

		if routeAggregationRule, _, err := client.RouteAggregationRulesApi.GetRouteAggregationRuleByUuid(ctx, routeAggregationId, rs.Primary.ID).Execute(); err == nil {
			if routeAggregationRule.GetState() == fabricv4.ROUTEAGGREGATIONRULESTATE_PROVISIONED {
				return fmt.Errorf("fabric stream %s still exists and is %s",
					rs.Primary.ID, routeAggregationRule.GetState())
			}
		}
	}
	return nil
}
