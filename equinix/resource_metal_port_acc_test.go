package equinix

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func confAccMetalPort_base(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port-test-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-metal-port-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%s"

  lifecycle {
	ignore_changes = [
	  plan,
	  metro
	]
  }
}

locals {
  bond0_id = [for p in equinix_metal_device.test.ports: p.id if p.name == "bond0"][0]
  eth1_id  = [for p in equinix_metal_device.test.ports: p.id if p.name == "eth1"][0]
  eth0_id  = [for p in equinix_metal_device.test.ports: p.id if p.name == "eth0"][0]
}

`, confAccMetalDevice_base(preferable_plans, preferable_metros), name, testDeviceTerminationTime())
}

func confAccMetalPort_L3(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id    = local.bond0_id
  bonded     = true
  depends_on = [
	equinix_metal_port.eth1,
  ]
}

resource "equinix_metal_port" "eth1" {
  port_id = local.eth1_id
  bonded  = true
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_L2Bonded(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2  = true
  bonded  = true
  reset_on_delete = true
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_L2Individual(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2  = true
  bonded  = false
  reset_on_delete = true
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_HybridUnbonded(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2  = false
  bonded  = true
  depends_on = [
	equinix_metal_port.eth1,
  ]
}

resource "equinix_metal_port" "eth1" {
  port_id = local.eth1_id
  bonded  = false
  reset_on_delete = true
}

`, confAccMetalPort_base(name))
}

func confAccMetalPort_HybridBonded(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id  = local.bond0_id
  layer2   = false
  bonded   = true
  vlan_ids = [equinix_metal_vlan.test.id]
  reset_on_delete = true
}

resource "equinix_metal_vlan" "test" {
  description = "tfacc-vlan test"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}
`, confAccMetalPort_base(name))
}

func confAccMetalPort_HybridBondedVxlan(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id   = local.bond0_id
  layer2    = false
  bonded    = true
  vxlan_ids = [equinix_metal_vlan.test1.vxlan, equinix_metal_vlan.test2.vxlan]
  reset_on_delete = true
}

resource "equinix_metal_vlan" "test1" {
  description = "tfacc-vlan test1"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
  vxlan       = 1001
}

resource "equinix_metal_vlan" "test2" {
  description = "tfacc-vlan test2"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
  vxlan       = 1002
}
`, confAccMetalPort_base(name))
}

func TestAccMetalPort_hybridBondedVxlan(t *testing.T) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalPortDestroyed,
		Steps: []resource.TestStep{
			{
				Config: confAccMetalPort_HybridBondedVxlan(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "vxlan_ids.#", "2"),
					resource.TestMatchResourceAttr("equinix_metal_port.bond0", "vxlan_ids.0",
						regexp.MustCompile("1001|1002")),
					resource.TestMatchResourceAttr("equinix_metal_port.bond0", "vxlan_ids.1",
						regexp.MustCompile("1001|1002")),
				),
			},
			{
				// Remove equinix_metal_port resources to trigger reset_on_delete
				Config: confAccMetalPort_base(rs),
			},
			{
				Config: confAccMetalPort_L3(rs),
			},
		},
	})
}

func testAccMetalPortTemplate(t *testing.T, conf func(string) string, expectedType string) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalPortDestroyed,
		Steps: []resource.TestStep{
			{
				Config: conf(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "name", "bond0"),
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "type", "NetworkBondPort"),
					resource.TestCheckResourceAttrSet("equinix_metal_port.bond0", "bonded"),
					resource.TestCheckResourceAttrSet("equinix_metal_port.bond0", "disbond_supported"),
					resource.TestCheckResourceAttrSet("equinix_metal_port.bond0", "port_id"),
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "network_type", expectedType),
				),
			},
			{
				ResourceName:            "equinix_metal_port.bond0",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"reset_on_delete"},
			},
			{
				// Remove equinix_metal_port resources to trigger reset_on_delete
				Config: confAccMetalPort_base(rs),
			},
			{
				Config: confAccMetalPort_L3(rs),
			},
			{
				ResourceName: "equinix_metal_port.bond0",
				ImportState:  true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "network_type", "layer3"),
				),
			},
		},
	})
}

func TestAccMetalPort_L2Bonded(t *testing.T) {
	testAccMetalPortTemplate(t, confAccMetalPort_L2Bonded, "layer2-bonded")
}

func TestAccMetalPort_L2Individual(t *testing.T) {
	testAccMetalPortTemplate(t, confAccMetalPort_L2Individual, "layer2-individual")
}

func TestAccMetalPort_hybridUnbonded(t *testing.T) {
	testAccMetalPortTemplate(t, confAccMetalPort_HybridUnbonded, "hybrid")
}

func TestAccMetalPort_hybridBonded(t *testing.T) {
	testAccMetalPortTemplate(t, confAccMetalPort_HybridBonded, "hybrid-bonded")
}

func testAccMetalPortDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	port_ids := []string{}

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "equinix_metal_port" {
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
