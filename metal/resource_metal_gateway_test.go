package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccMetalGatewayConfig_PrivateIPv4() string {
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
`)
}

func TestAccMetalGateway_PrivateIPv4(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalGatewayDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalGatewayConfig_PrivateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_gateway.test", "project_id",
						"metal_project.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_gateway.test", "private_ipv4_subnet_size", "8"),
				),
			},
		},
	})
}

func testAccMetalGatewayConfig_ExistingReservation() string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-gateway-test"
}

resource "metal_vlan" "test" {
    description = "test VLAN in SV"
    metro       = "sv"
    project_id  = metal_project.test.id
}

resource "metal_reserved_ip_block" "test" {
    project_id = metal_project.test.id
    metro      = "sv"
    quantity   = 2
}

resource "metal_gateway" "test" {
    project_id        = metal_project.test.id
    vlan_id           = metal_vlan.test.id
    ip_reservation_id = metal_reserved_ip_block.test.id
}
`)
}

func TestAccMetalGateway_ExistingReservation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalGatewayDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalGatewayConfig_PrivateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_gateway.test", "project_id",
						"metal_project.test", "id"),
					resource.TestCheckResourceAttrPair(
						"metal_gateway.test", "ip_reservation_id",
						"metal_reserved_ip_block.test", "id"),
				),
			},
		},
	})
}

func testAccCheckMetalGatewayDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_gateway" {
			continue
		}
		if _, _, err := client.MetalGateways.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Gateway still exists")
		}
	}

	return nil
}

func TestAccMetalGateway_importBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalGatewayDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalGatewayConfig_PrivateIPv4(),
			},
			{
				ResourceName:      "metal_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
