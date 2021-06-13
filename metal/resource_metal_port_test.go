package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func confAccMetalPort_base(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-port-test-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-metal-port-test"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}

locals {
  bond0_id = [for p in metal_device.test.ports: p.id if p.name == "bond0"][0]
  eth1_id = [for p in metal_device.test.ports: p.id if p.name == "eth1"][0]
}

`, name)
}

func confAccMetalPort_L3(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_port" "bond0" {
  port_id = local.bond0_id
  bonded = true
  depends_on = [
    metal_port.eth1,
  ]
}

resource "metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = true
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_L2Bonded(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = true
  bonded = true
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_L2Individual(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = true
  bonded = false
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_HybridUnbonded(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = false
  bonded = true
  depends_on = [
    metal_port.eth1,
  ]
}

resource "metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = false
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_HybridBonded(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_port" "bond0" {
  port_id = local.bond0_id
  layer2 = false
  bonded = true
  vlan_ids = [metal_vlan.test.id]
}

resource "metal_vlan" "test" {
  description = "test"
  metro = "sv"
  project = metal_project.test.id
}

resource "metal_port" "eth1" {
  port_id = local.eth1_id
  bonded = false
}

`, confAccMetalPort_base(name))
}

func metalPortTestTemplate(t *testing.T, conf func(string) string, expectedType string) {
	rs := acctest.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalPortDestroy,
		Steps: []resource.TestStep{
			{
				Config: conf(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("metal_port.bond0", "network_type", expectedType),
				),
			},
			{
				Config: confAccMetalPort_L3(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("metal_port.bond0", "network_type", "layer3"),
				),
			},
		},
	})
}

func TestAccMetalPort_L2Bonded(t *testing.T) {
	metalPortTestTemplate(t, confAccMetalPort_L2Bonded, "layer2-bonded")
}

func TestAccMetalPort_L2Individual(t *testing.T) {
	metalPortTestTemplate(t, confAccMetalPort_L2Individual, "layer2-individual")
}

func TestAccMetalPort_HybridUnbonded(t *testing.T) {
	metalPortTestTemplate(t, confAccMetalPort_HybridUnbonded, "hybrid")
}

func TestAccMetalPort_HybridBonded(t *testing.T) {
	metalPortTestTemplate(t, confAccMetalPort_HybridBonded, "hybrid-bonded")
}

func testAccMetalPortDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	port_ids := []string{}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "metal_port" {
			shouldReset := rs.Primary.Attributes["reset_on_delete"]
			if shouldReset == "true" {
				port_ids = append(port_ids, rs.Primary.ID)
			}
		}
	}
	for _, pid := range port_ids {
		p, _, err := client.Ports.Get(pid, nil)
		if err != nil {
			return fmt.Errorf("Error getting port %s during destroy check", pid)
		}
		err = portProperlyDestroyed(p)
		if err != nil {
			return err
		}
	}
	return nil
}
