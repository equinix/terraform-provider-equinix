package route_filter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/route_filter"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricRouteFilterPolicy_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteFilterPolicyConfig("RF_Policy_PFCR"),
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

func testAccFabricRouteFilterPolicyConfig(policyName string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_filter" "test" {
			name = "%s",
			project {
				projectId = "291639000636552"
			},
			type = "BGP_IPv4_PREFIX_FILTER",
			description = "Route Filter Policy for X Purpose",
		}
	`, policyName)
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
