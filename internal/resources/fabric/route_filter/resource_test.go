package route_filter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricRouteFilterPolicy_PFCR(t *testing.T) {
	routeFilterName, routeFilterUpdatedName := "RF_Policy_PFCR", "RF_Filter_PFCR"
	routeFilterDescription, routeFilterUpdatedDescription := "Route Filter Policy for X Purpose", "Route Filter Policy for Y Purpose"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteFilterPolicyConfig(routeFilterName, routeFilterDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_route_filter.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "name", routeFilterName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "state", string(fabricv4.ROUTEFILTERSTATE_PROVISIONED)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "not_matched_rule_action", "DENY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "rules_count", "0"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "description", routeFilterDescription),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccFabricRouteFilterPolicyConfig(routeFilterUpdatedName, routeFilterUpdatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_route_filter.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "name", routeFilterUpdatedName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "state", string(fabricv4.ROUTEFILTERSTATE_REPROVISIONING)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "not_matched_rule_action", "DENY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "rules_count", "0"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_filter.test", "description", routeFilterUpdatedDescription),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricRouteFilterPolicyConfig(policyName, policyDescription string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_filter" "test" {
			name = "%s"
			project {
				project_id = "291639000636552"
			}
			type = "BGP_IPv4_PREFIX_FILTER"
			description = "%s"
		}
	`, policyName, policyDescription)
}

func CheckRouteFilterDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_route_filter" {
			continue
		}

		err := route_filter.WaitForDeletion(rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
