package gateway_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccDataSourceMetalGateway_privateIPv4(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalGatewayCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalGatewayConfig_privateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.equinix_metal_gateway.test", "project_id",
						"equinix_metal_project.test", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_gateway.test", "private_ipv4_subnet_size", "8"),
				),
			},
		},
	})
}

func testAccDataSourceMetalGatewayConfig_privateIPv4() string {
	return `
resource "equinix_metal_project" "test" {
    name = "tfacc-gateway-test"
}

resource "equinix_metal_vlan" "test" {
    description = "tfacc-vlan test VLAN in SV"
    metro       = "sv"
    project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_gateway" "test" {
    project_id               = equinix_metal_project.test.id
    vlan_id                  = equinix_metal_vlan.test.id
    private_ipv4_subnet_size = 8
}

data "equinix_metal_gateway" "test" {
    gateway_id = equinix_metal_gateway.test.id
}
`
}

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccDataSourceMetalProjectSSHKey_upgradeFromVersion(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccMetalGatewayCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.27.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: testAccDataSourceMetalGatewayConfig_privateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.equinix_metal_gateway.test", "project_id",
						"equinix_metal_project.test", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_gateway.test", "private_ipv4_subnet_size", "8"),
				),
			},
			{
				ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
				Config:                   testAccDataSourceMetalGatewayConfig_privateIPv4(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
