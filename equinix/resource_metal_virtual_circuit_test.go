package equinix

// I am not sure what to do with the test code. It's useful, but it won't run
// either:
// * unless the Connections are automatically approved.
// * unless we specify an existing Connection and Project for testing.
//
// I can remove this file from the PR if it looks too bad here.

/*


import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalVirtualCircuitDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_virtual_circuit" {
			continue
		}
		if _, _, err := client.VirtualCircuits.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("VirtualCircuit still exists")
		}
	}

	return nil
}

func tconf(randstr string, randint int) string {
	return fmt.Sprintf(`
        locals {
                project_id = "52000fb2-ee46-4673-93a8-de2c2bdba33b"
                conn_id = "73f12f29-3e19-43a0-8e90-ae81580db1e9"
        }

        data "equinix_metal_connection" test {
            connection_id = local.conn_id
        }

        resource "equinix_metal_vlan" "test" {
            project_id = local.project_id
            metro      = data.equinix_metal_connection.test.metro
        }

        resource "equinix_metal_virtual_circuit" "test" {
            connection_id = local.conn_id
            project_id = local.project_id
            port_id = data.equinix_metal_connection.test.ports[0].id
            vlan_id = metal_vlan.test.id
            nni_vlan = %d
        }


        `,
		randint)
}

func testAccMetalVirtualCircuitConfig_Dedicated(randstr string, randint int) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        // No project ID. We only use the project resource to get org_id
        resource "equinix_metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = metal_project.test.organization_id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "dedicated"
        }

        resource "equinix_metal_vlan" "test" {
            project_id = metal_project.test.id
            metro      = "sv"
        }

        resource "equinix_metal_virtual_circuit" "test" {
            connection_id = metal_connection.test.id
            project_id = metal_project.test.id
            port_id = metal_connection.test.ports[0].id
            vlan_id = metal_vlan.test.id
            nni_vlan = %d
        }


        `,
		randstr, randstr, randint)
}

func TestAccMetalVirtualCircuit_Dedicated(t *testing.T) {

	rs := acctest.RandString(10)
	ri := acctest.RandIntRange(1024, 1093)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVirtualCircuitDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccMetalVirtualCircuitConfig_Dedicated(rs, ri),
				Config: tconf(rs, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test", "vlan_id",
						"equinix_metal_vlan.test", "id",
					),
				),
			},
				{
					ResourceName:      "equinix_metal_virtual_circuit.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
		},
	})
}
*/
