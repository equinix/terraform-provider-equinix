package route_filter_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricRouteFilterPolicy_DataSources_PFCR(t *testing.T) {
	routeFilterName := "RF_DS_Policy_PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteFilterPolicyDataSourcesConfig(routeFilterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_filter.rf_policy", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter.rf_policy", "name", routeFilterName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter.rf_policy", "type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter.rf_policy", "state", "PROVISIONED"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter.rf_policy", "not_matched_rule_action", "DENY"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter.rf_policy", "rules_count", "0"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filter.rf_policy", "description", "Route Filter Policy for X Purpose"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_filters.rf_policies", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filters.rf_policies", "data.0.name", routeFilterName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filters.rf_policies", "data.0.type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filters.rf_policies", "data.0.state", "PROVISIONED"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filters.rf_policies", "data.0.not_matched_rule_action", "DENY"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filters.rf_policies", "data.0.rules_count", "0"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_filters.rf_policies", "data.0.description", "Route Filter Policy for X Purpose"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricRouteFilterPolicyDataSourcesConfig(policyName string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_filter" "test" {
			name = "%s"
			project {
				project_id = "291639000636552"
			}
			type = "BGP_IPv4_PREFIX_FILTER"
			description = "Route Filter Policy for X Purpose"
		}

		data "equinix_fabric_route_filter" "rf_policy" {
			uuid = equinix_fabric_route_filter.test.id
		}
		
		data "equinix_fabric_route_filters" "rf_policies" {
			filter {
				property = "/type"
				operator = "="
				values 	 = ["BGP_IPv4_PREFIX_FILTER"]
			}
			filter {
				property = "/state"
				operator = "="
				values   = ["PROVISIONED"]
			}
			filter {
				property = "/project/projectId"
				operator = "="
				values = ["291639000636552"]
			}
			filter {
				property = "/name"
				operator = "="
				values = [equinix_fabric_route_filter.test.name]
			}
			pagination {
				offset = 0
				limit = 5
				total = 25
			}
			sort {
				direction = "ASC"
				property = "/name"
			}
		}
	`, policyName)
}
