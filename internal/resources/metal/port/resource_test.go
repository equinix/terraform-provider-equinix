package port_test

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/port"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var (
	matchErrPortReadyTimeout = regexp.MustCompile(".* timeout while waiting for state to become 'completed'.*")
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
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}

locals {
  bond0_id = [for p in equinix_metal_device.test.ports: p.id if p.name == "bond0"][0]
  eth1_id  = [for p in equinix_metal_device.test.ports: p.id if p.name == "eth1"][0]
  eth0_id  = [for p in equinix_metal_device.test.ports: p.id if p.name == "eth0"][0]
}

`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), name, acceptance.TestDeviceTerminationTime())
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

func confAccMetalPort_L2IndividualNativeVlan(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id = local.bond0_id
  layer2  = true
  bonded  = false
  reset_on_delete = true
}

resource "equinix_metal_port" "eth1" {
  port_id = local.eth1_id
  bonded  = false
  reset_on_delete = true
  vlan_ids = [equinix_metal_vlan.test1.id, equinix_metal_vlan.test2.id]
  native_vlan_id = equinix_metal_vlan.test1.id
  depends_on = [
	equinix_metal_port.bond0,
  ]
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

func confAccMetalPort_HybridBonded_timeout(rInt int, name, createTimeout, updateTimeout string) string {
	if createTimeout == "" {
		createTimeout = "20m"
	}
	if updateTimeout == "" {
		updateTimeout = "20m"
	}

	return fmt.Sprintf(`
%s

resource "equinix_metal_port" "bond0" {
  port_id  = local.bond0_id
  layer2   = false
  bonded   = true
  reset_on_delete = true
  vlan_ids = [equinix_metal_vlan.test.id]
  timeouts {
    create = "%s"
	update = "%s"
  }
  depends_on = [equinix_metal_vlan.test]
}

resource "equinix_metal_vlan" "test" {
  description = "tfacc-vlan test-%d"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}
`, confAccMetalPort_base(name), createTimeout, updateTimeout, rInt)
}

func TestAccMetalPort_hybridBondedVxlan(t *testing.T) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalPortDestroyed,
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

func TestAccMetalPort_L2IndividualNativeVlan(t *testing.T) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalPortDestroyed,
		Steps: []resource.TestStep{
			{
				Config: confAccMetalPort_L2IndividualNativeVlan(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.eth1", "vxlan_ids.#", "2"),
					resource.TestMatchResourceAttr("equinix_metal_port.eth1", "vxlan_ids.0",
						regexp.MustCompile("1001|1002")),
					resource.TestMatchResourceAttr("equinix_metal_port.eth1", "vxlan_ids.1",
						regexp.MustCompile("1001|1002")),
					resource.TestCheckResourceAttr("equinix_metal_port.eth1", "bonded", "false"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port.eth1", "native_vlan_id",
						"equinix_metal_vlan.test1", "id"),
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
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalPortDestroyed,
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
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

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
		p, resp, err := client.PortsApi.FindPortById(context.Background(), pid).Execute()
		if err != nil {
			if resp != nil && slices.Contains([]int{http.StatusNotFound, http.StatusForbidden}, resp.StatusCode) {
				continue
			}
			return fmt.Errorf("Error getting port %s during destroy check: %v", pid, err)
		}
		err = port.ProperlyDestroyed(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccWaitForPortActive(deviceName, portName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		rs, ok := state.RootModule().Resources[deviceName]
		if !ok {
			return "", fmt.Errorf("Device Not found in the state: %s", deviceName)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No Record ID is set")
		}

		meta := acceptance.TestAccProvider.Meta()
		client := meta.(*config.Config).NewMetalClientForTesting()
		device, _, err := client.DevicesApi.FindDeviceById(context.Background(), rs.Primary.ID).Include([]string{"ports"}).Execute()
		if err != nil {
			return "", fmt.Errorf("error while fetching device with Id [%s], error: %w", rs.Primary.ID, err)
		}
		if device == nil {
			return "", fmt.Errorf("Not able to find devices with Id [%s]", rs.Primary.ID)
		}
		if len(device.NetworkPorts) == 0 {
			return "", fmt.Errorf("Found no ports for the device with Id [%s]", rs.Primary.ID)
		}

		for _, port := range device.NetworkPorts {
			if port.GetName() == portName {
				return port.GetId(), nil
			}
		}

		return "", fmt.Errorf("No port with name [%s] found", portName)
	}
}

func TestAccMetalPortCreate_hybridBonded_timeout(t *testing.T) {
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	deviceName := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalPortDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      confAccMetalPort_HybridBonded_timeout(rInt, rs, "5s", ""),
				ExpectError: matchErrPortReadyTimeout,
			},
			{
				/**
				Step 1 errors out, state doesnt have port, need to import that in the state before deleting
				*/
				ResourceName:       "equinix_metal_port.bond0",
				ImportState:        true,
				ImportStateIdFunc:  testAccWaitForPortActive(deviceName, "bond0"),
				ImportStatePersist: true,
			},
			{
				ResourceName: "equinix_metal_port.bond0",
				ImportState:  true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "network_type", "layer3"),
				),
			},
			{
				Config: confAccMetalPort_HybridBonded_timeout(rInt, rs, "5s", ""),
			},
			{
				Config:  confAccMetalPort_HybridBonded_timeout(rInt, rs, "5s", ""),
				Destroy: true,
			},
		},
	})
}

func TestAccMetalPortUpdate_hybridBonded_timeout(t *testing.T) {
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalPortDestroyed,
		Steps: []resource.TestStep{
			{
				Config: confAccMetalPort_HybridBonded_timeout(rInt, rs, "", "5s"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "name", "bond0"),
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "type", "NetworkBondPort"),
					resource.TestCheckResourceAttrSet("equinix_metal_port.bond0", "bonded"),
					resource.TestCheckResourceAttrSet("equinix_metal_port.bond0", "disbond_supported"),
					resource.TestCheckResourceAttrSet("equinix_metal_port.bond0", "port_id"),
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "network_type", "hybrid-bonded"),
				),
			},
			{
				Config:      confAccMetalPort_HybridBonded_timeout(rInt+1, rs, "", "5s"),
				ExpectError: matchErrPortReadyTimeout,
			},
			{
				ResourceName: "equinix_metal_port.bond0",
				ImportState:  true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_port.bond0", "network_type", "layer3"),
				),
			},
			{
				Config: confAccMetalPort_HybridBonded_timeout(rInt+1, rs, "", ""),
			},
			{
				Config:  confAccMetalPort_HybridBonded_timeout(rInt+1, rs, "", ""),
				Destroy: true,
			},
		},
	})
}
