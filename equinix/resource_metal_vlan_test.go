package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalVlanConfig_metro(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}
`, projSuffix, metro, desc)
}

func TestAccMetalVlan_Metro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_metro(rs, metro, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "metro", metro),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "facility", ""),
				),
			},
		},
	})
}

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
					testAccCheckMetalVlanExists("equinix_metal_vlan.foovlan", &vlan),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "description", "testvlan"),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "facility", fac),
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
		if rs.Type != "equinix_metal_vlan" {
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
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = "${equinix_metal_project.foobar.id}"
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
				ResourceName:      "equinix_metal_vlan.foovlan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
