package route_filter_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricRouteFilterPolicy_DataSources_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteFilterPolicyDataSourcesConfig("RF_Policy_PFCR"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_route_filter.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "state", "PROVISIONED"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "not_matched_rules_action", "0"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "rules_count", "0"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "description", "Route Filter Policy for X Purpose"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricRouteFilterPolicyDataSourcesConfig(policyName string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_filter" "test" {
			name = "%s",
			project {
				projectId = "291639000636552"
			},
			type = "BGP_IPv4_PREFIX_FILTER",
			description = "Route Filter Policy for X Purpose",
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
				values = ["%s"]
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
