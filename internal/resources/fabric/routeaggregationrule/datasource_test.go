package routeaggregationrule_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFabricRouteAggregationRuleDataSourcesConfig(name, description string) string {
	return fmt.Sprintf(`

		resource "equinix_fabric_route_aggregation" "test" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "test-aggregation"
		  description = "Test Route Aggregation"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}

		resource "equinix_fabric_route_aggregation_rule" "new-rar" {
  			route_aggregation_id = equinix_fabric_route_aggregation.test.id
  			name = "%[1]s"
  			description = "%[2]s"
  			prefix = "192.166.0.0/24"
		}

		data "equinix_fabric_route_aggregation_rule" "data_rar" {
			depends_on = [equinix_fabric_route_aggregation_rule.new-rar]
			route_aggregation_id = equinix_fabric_route_aggregation.test.id
  			route_aggregation_rule_id = equinix_fabric_route_aggregation_rule.new-rar.id
		}

		
		data "equinix_fabric_route_aggregation_rules" "data_rars" {
			depends_on = [equinix_fabric_route_aggregation_rule.new-rar]
  			route_aggregation_id = equinix_fabric_route_aggregation.test.id
  			pagination = {
    			limit = 32
    			offset = 0
  			}
		}
	`, name, description)
}

func TestAccFabricRouteAggregationRuleDataSources_PNFV(t *testing.T) {
	routeAggregationName := "route_agg_rule_PFCR"
	routeAggregatioDescription := "route aggregation rule PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteAggregationRuleDataSourcesConfig(routeAggregationName, routeAggregatioDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation_rule.data_rar", "name", routeAggregationName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation_rule.data_rar", "type", "BGP_IPv4_PREFIX_AGGREGATION_RULE"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation_rule.data_rar", "description", routeAggregatioDescription),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rule.data_rar", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rule.data_rar", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rule.data_rar", "prefix"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rule.data_rar", "prefix"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rule.data_rar", "state"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rule.data_rar", "change_log.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rules.data_rars", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rules.data_rars", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rules.data_rars", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rules.data_rars", "data.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rules.data_rars", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation_rules.data_rars", "data.0.change_log.created_by"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregation_rules.data_rars", "data.#", "1"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregation_rules.data_rars", "pagination.%", "5"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregation_rules.data_rars", "pagination.limit", "32"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregation_rules.data_rars", "pagination.offset", "0"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
