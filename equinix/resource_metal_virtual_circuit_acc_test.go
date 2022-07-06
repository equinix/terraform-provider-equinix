package equinix

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("equinix_metal_virtual_circuit", &resource.Sweeper{
		Name:         "equinix_metal_virtual_circuit",
		Dependencies: []string{},
		F:            testSweepVirtualCircuits,
	})
}

func testSweepVirtualCircuits(region string) error {
	log.Printf("[DEBUG] Sweeping VirtualCircuits")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping VirtualCircuits: %s", err)
	}
	metal := config.NewMetalClient()
	orgList, _, err := metal.Organizations.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting organization list for sweeping VirtualCircuits: %s", err)
	}
	vcs := map[string]*packngo.VirtualCircuit{}
	for _, org := range orgList {
		conns, _, err := metal.Connections.OrganizationList(org.ID, &packngo.GetOptions{Includes: []string{"ports"}})
		if err != nil {
			return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting connections list for sweeping VirtualCircuits: %s", err)
		}
		for _, conn := range conns {
			for _, port := range conn.Ports {
				for _, vc := range port.VirtualCircuits {
					if isSweepableTestResource(vc.Name) {
						vcs[vc.ID] = &vc
					}
				}
			}
		}
	}
	for _, vc := range vcs {
		log.Printf("[INFO][SWEEPER_LOG] Deleting VirtualCircuit: %s", vc.Name)
		_, err := metal.VirtualCircuits.Delete(vc.ID)
		if err != nil {
			return fmt.Errorf("[INFO][SWEEPER_LOG] Error deleting VirtualCircuit: %s", err)
		}
	}

	return nil
}

func testAccMetalVirtualCircuitCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_virtual_circuit" {
			continue
		}
		if _, _, err := client.VirtualCircuits.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal VirtualCircuit still exists")
		}
	}

	return nil
}


func testAccMetalConnectionConfig_vc(randint int) string {
	// Dedicated connection in DA metro
	testConnection := os.Getenv(metalDedicatedConnIDEnvVar)

	return fmt.Sprintf(`
        locals {
            conn_id = "%s"
        }

        data "equinix_metal_connection" test {
            connection_id = local.conn_id
        }

        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%[2]d"
        }

        resource "equinix_metal_vlan" "test" {
            project_id = equinix_metal_project.test.id
            metro      = data.equinix_metal_connection.test.metro
			description = "tfacc-vlan test"
        }

        resource "equinix_metal_virtual_circuit" "test" {
            name = "tfacc-vc-%[2]d"
            description = "tfacc-vc-%[2]d"
            connection_id = data.equinix_metal_connection.test.connection_id
            project_id = equinix_metal_project.test.id
            port_id = data.equinix_metal_connection.test.ports[0].id
            vlan_id = equinix_metal_vlan.test.id
            nni_vlan = %[2]d
        }
        `,
		testConnection, randint)
}

func testAccMetalConnectionConfig_vcds(randint int) string {
	return testAccMetalConnectionConfig_vc(randint) + `
	data "equinix_metal_virtual_circuit" "test" {
		virtual_circuit_id = equinix_metal_virtual_circuit.test.id
	}
	`
}

func TestAccMetalVirtualCircuit_dedicated(t *testing.T) {
	ri := acctest.RandIntRange(1024, 1093)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVirtualCircuitCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_vc(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "vlan_id",
						"equinix_metal_vlan.test", "id",
					),
				),
			},
			{
				ResourceName:            "equinix_metal_virtual_circuit.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"connection_id"},
			},
			{
				Config: testAccMetalConnectionConfig_vcds(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "id",
						"data.equinix_metal_virtual_circuit.test", "virtual_circuit_id",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "speed",
						"data.equinix_metal_virtual_circuit.test", "speed",
					),

					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "port_id",
						"data.equinix_metal_virtual_circuit.test", "port_id",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "vlan_id",
						"data.equinix_metal_virtual_circuit.test", "vlan_id",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "nni_vlan",
						"data.equinix_metal_virtual_circuit.test", "nni_vlan",
					),
				),
			},
		},
	})
}
