package route_filter_rule_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricRouteFilterRule_DataSources_PFCR(t *testing.T) {
	routeFilterRuleName, routeFilterRuleDescription := "RF_DS_Rule_PFCR", "Route Filter Rule"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteFilterRuleDataSourcesConfig(routeFilterRuleName, routeFilterRuleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_filter_rule.rf_rule", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_filter_rule.rf_rule", "route_filter_id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rule.rf_rule", "name", routeFilterRuleName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rule.rf_rule", "prefix", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rule.rf_rule", "prefix_match", "exact"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rule.rf_rule", "description", routeFilterRuleDescription),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_filter_rules.rf_rules", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_filter_rules.rf_rules", "route_filter_id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rules.rf_rules", "data.0.name", routeFilterRuleName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rules.rf_rules", "data.0.prefix", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rules.rf_rules", "data.0.prefix_match", "exact"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter_rules.rf_rules", "data.0.description", routeFilterRuleDescription),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricRouteFilterRuleDataSourcesConfig(policyName, description string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_filter" "test" {
			name = "rf_ds_test_PFCR"
			project {
				project_id = "291639000636552"
			}
			type = "BGP_IPv4_PREFIX_FILTER"
			description = "Route Filter Policy for X Purpose"
		}

		resource "equinix_fabric_route_filter_rule" "test" {
			route_filter_id = equinix_fabric_route_filter.test.id
			name = "%s"
			prefix = "192.168.0.0/24"
			prefix_match = "exact"
			description = "%s"
		}

		data "equinix_fabric_route_filter_rule" "rf_rule" {
			route_filter_id = equinix_fabric_route_filter.test.id
			uuid = equinix_fabric_route_filter_rule.test.id
		}
		
		data "equinix_fabric_route_filter_rules" "rf_rules" {
			depends_on = [ equinix_fabric_route_filter_rule.test ]
			route_filter_id = equinix_fabric_route_filter.test.id
		}
	`, policyName, description)
}
