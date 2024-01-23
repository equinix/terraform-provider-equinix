package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetalGateway_privateIPv4(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalGatewayCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalGatewayConfig_privateIPv4(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_gateway.test", "project_id",
						"equinix_metal_project.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_gateway.test", "private_ipv4_subnet_size", "8"),
				),
			},
		},
	})
}

func testAccMetalGatewayConfig_privateIPv4() string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-gateway-test"
}

resource "equinix_metal_vlan" "test" {
    description = "tfacc-vlan VLAN in SV"
    metro       = "sv"
    project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_gateway" "test" {
    project_id               = equinix_metal_project.test.id
    vlan_id                  = equinix_metal_vlan.test.id
    private_ipv4_subnet_size = 8
}
`)
}

func TestAccMetalGateway_existingReservation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalGatewayCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalGatewayConfig_existingReservation(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_gateway.test", "project_id",
						"equinix_metal_project.test", "id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_gateway.test", "ip_reservation_id",
						"equinix_metal_reserved_ip_block.test", "id"),
				),
			},
		},
	})
}

func testAccMetalGatewayConfig_existingReservation() string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-gateway-test"
}

resource "equinix_metal_vlan" "test" {
    description = "tfacc-vlan in SV"
    metro       = "sv"
    project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_reserved_ip_block" "test" {
    project_id = equinix_metal_project.test.id
    metro      = "sv"
    quantity   = 8
}

resource "equinix_metal_gateway" "test" {
    project_id        = equinix_metal_project.test.id
    vlan_id           = equinix_metal_vlan.test.id
    ip_reservation_id = equinix_metal_reserved_ip_block.test.id
}
`)
}

func testAccMetalGatewayCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_gateway" {
			continue
		}
		if _, _, err := client.MetalGateways.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Gateway still exists")
		}
	}

	return nil
}

func TestAccMetalGateway_importBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalGatewayCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalGatewayConfig_privateIPv4(),
			},
			{
				ResourceName:      "equinix_metal_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
