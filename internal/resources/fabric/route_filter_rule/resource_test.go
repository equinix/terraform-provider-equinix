package route_filter_rule_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter_rule"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricRouteFilterRule_PFCR(t *testing.T) {
	routeFilterRuleName, routeFilterRuleUpdatedName := "RF_Rule_PFCR", "RF_RuleB_PFCR"
	routeFilterRulePrefix, routeFilterRulePrefixUpdated := "192.168.0.0/24", "192.172.0.0/24"
	routeFilterRuleDescription := "Route Filter Rule for X Purpose"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteFilterRuleConfig(routeFilterRuleName, routeFilterRulePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_route_filter_rule.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "name", routeFilterRuleName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "prefix", routeFilterRulePrefix),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "prefix_match", "exact"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "description", routeFilterRuleDescription),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccFabricRouteFilterRuleConfig(routeFilterRuleUpdatedName, routeFilterRulePrefixUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_route_filter_rule.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "name", routeFilterRuleUpdatedName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "prefix", routeFilterRulePrefixUpdated),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "prefix_match", "exact"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter_rule.test", "description", routeFilterRuleDescription),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricRouteFilterRuleConfig(policyName, policyPrefix string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_filter" "test" {
			name = "rf_test_PFCR"
			project {
				project_id = "291639000636552"
			}
			type = "BGP_IPv4_PREFIX_FILTER"
			description = "Route Filter Policy for X Purpose"
		}

		resource "equinix_fabric_route_filter_rule" "test" {
			route_filter_id = equinix_fabric_route_filter.test.id
			name = "%s"
			prefix = "%s"
			prefix_match = "exact"
			description = "Route Filter Rule for X Purpose"
		}
	`, policyName, policyPrefix)
}

func CheckRouteFilterRuleDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_route_filter_rule" {
			continue
		}

		routeFilterId := rs.Primary.Attributes["route_filter_id"]

		err := route_filter_rule.WaitForDeletion(routeFilterId, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
