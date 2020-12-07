package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalVlan_Basic(t *testing.T) {
	var vlan packngo.VirtualNetwork
	rs := acctest.RandString(10)
	fac := "ewr1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_var(rs, fac, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalVlanExists("metal_vlan.foovlan", &vlan),
					resource.TestCheckResourceAttr(
						"metal_vlan.foovlan", "description", "testvlan"),
					resource.TestCheckResourceAttr(
						"metal_vlan.foovlan", "facility", fac),
				),
			},
		},
	})
}

func testAccCheckMetalVlanExists(n string, vlan *packngo.VirtualNetwork) resource.TestCheckFunc {
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

func testAccCheckMetalVlanDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Vlan still exists")
		}
	}

	return nil
}

func testAccCheckMetalVlanConfig_var(projSuffix, facility, desc string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "metal_vlan" "foovlan" {
    project_id = "${metal_project.foobar.id}"
    facility = "%s"
    description = "%s"
}
`, projSuffix, facility, desc)
}

func TestAccMetalVlan_importBasic(t *testing.T) {
	rs := acctest.RandString(10)
	fac := "ewr1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_var(rs, fac, "testvlan"),
			},
			{
				ResourceName:      "metal_vlan.foovlan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
