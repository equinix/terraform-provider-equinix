package route_aggregation_rule_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFabricRouteAggregationRuleDataSourcesConfig(name, description string) string {
	return fmt.Sprintf(`

		
		resource "equinix_fabric_route_aggregation_rule" "new-rar" {
  			route_aggregation_id = "8f8a2ddb-25f8-416e-ad0a-202a9d2af9e1"
  			name = "%[1]s"
  			description = "%[2]s"
  			prefix = "192.166.0.0/24"
		}

		data "equinix_fabric_route_aggregation_rule" "data_rar" {
			route_aggregation_id = "8f8a2ddb-25f8-416e-ad0a-202a9d2af9e1"
  			route_aggregation_rule_id = equinix_fabric_route_aggregation_rule.new-rar.id
		}

		
		data "equinix_fabric_route_aggregation_rules" "data_rars" {
			depends_on = [equinix_fabric_route_aggregation_rule.new-rar,]
  			route_aggregation_id = "8f8a2ddb-25f8-416e-ad0a-202a9d2af9e1"
  			pagination = {
    			limit = 2
    			offset = 1
  			}
		}
	`, name, description)
}

func TestAccFabricRouteAggregationRuleDataSources_PFCR(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregation_rules.data_rars", "pagination.limit", "2"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregation_rules.data_rars", "pagination.offset", "1"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
