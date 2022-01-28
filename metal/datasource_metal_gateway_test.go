package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccDataSourceMetalGatewayConfig_PrivateIPv4() string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-gateway-test"
}

resource "metal_vlan" "test" {
    description = "test VLAN in SV"
    metro       = "sv"
    project_id  = metal_project.test.id
}

resource "metal_gateway" "test" {
    project_id               = metal_project.test.id
    vlan_id                  = metal_vlan.test.id
    private_ipv4_subnet_size = 8
}

data "metal_gateway" "test" {
    gateway_id = metal_gateway.test.id
}
`)
}

func TestAccDataSourceMetalGateway_PrivateIPv4(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalGatewayDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalGatewayConfig_PrivateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.metal_gateway.test", "project_id",
						"metal_project.test", "id"),
					resource.TestCheckResourceAttr(
						"data.metal_gateway.test", "private_ipv4_subnet_size", "8"),
				),
			},
		},
	})
}
