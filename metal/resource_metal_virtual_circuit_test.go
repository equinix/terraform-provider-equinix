package metal

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
		if rs.Type != "metal_connection" {
			continue
		}
		if _, _, err := client.VirtualCircuits.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("VirtualCircuit still exists")
		}
	}

	return nil
}

func testAccMetalVirtualCircuitConfig_Dedicated(randstr string, randint int) string {
	return fmt.Sprintf(`
        resource "metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        // No project ID. We only use the project resource to get org_id
        resource "metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = metal_project.test.organization_id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "dedicated"
        }

        resource "metal_vlan" "test" {
            project_id = metal_project.test.id
            metro      = "sv"
        }

        resource "metal_virtual_circuit" "test" {
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
				Config: testAccMetalVirtualCircuitConfig_Dedicated(rs, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_virtual_circuit.test", "vlan_id",
						"metal_vlan.test", "id",
					),
				),
			},
			{
				ResourceName:      "metal_virtual_circuit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
