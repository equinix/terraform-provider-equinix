package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketVlan_Basic(t *testing.T) {
	var vlan packngo.VirtualNetwork
	rs := acctest.RandString(10)
	fac := "ewr1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketVlanConfig_var(rs, fac, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVlanExists("packet_vlan.foovlan", &vlan),
					resource.TestCheckResourceAttr(
						"packet_vlan.foovlan", "description", "testvlan"),
					resource.TestCheckResourceAttr(
						"packet_vlan.foovlan", "facility", fac),
				),
			},
		},
	})
}

func testAccCheckPacketVlanExists(n string, vlan *packngo.VirtualNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundVlan, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundVlan.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundVlan)
		}

		*vlan = *foundVlan

		return nil
	}
}

func testAccCheckPacketVlanDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Vlan still exists")
		}
	}

	return nil
}

func testAccCheckPacketVlanConfig_var(projSuffix, facility, desc string) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "TerraformTestProject-%s"
}

resource "packet_vlan" "foovlan" {
    project_id = "${packet_project.foobar.id}"
    facility = "%s"
    description = "%s"
}
`, projSuffix, facility, desc)
}
