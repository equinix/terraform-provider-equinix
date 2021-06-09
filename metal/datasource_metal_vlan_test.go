package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalDatasourceVlanConfig_ByVxlanFacility(projSuffix, fac, desc string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    facility = "%s"
    description = "%s"
}

data "metal_vlan" "dsvlan" {
    facility = metal_vlan.foovlan.facility
    project_id = metal_vlan.foovlan.project_id
    vxlan = metal_vlan.foovlan.vxlan
}
`, projSuffix, fac, desc)
}

func TestAccMetalDatasourceVlan_ByVxlanFacility(t *testing.T) {
	rs := acctest.RandString(10)
	fac := "sv15"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByVxlanFacility(rs, fac, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "id",
						"data.metal_vlan.dsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalDatasourceVlanConfig_ByVxlanMetro(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "metal_vlan" "dsvlan" {
    metro = metal_vlan.foovlan.metro
    project_id = metal_vlan.foovlan.project_id
    vxlan = metal_vlan.foovlan.vxlan
}
`, projSuffix, metro, desc)
}

func TestAccMetalDatasourceVlan_ByVxlanMetro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByVxlanMetro(rs, metro, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "id",
						"data.metal_vlan.dsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalDatasourceVlanDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("DatasourceVlan still exists")
		}
	}

	return nil
}

func testAccCheckMetalDatasourceVlanConfig_ByVlanId(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "metal_vlan" "dsvlan" {
    vlan_id = metal_vlan.foovlan.id
}
`, projSuffix, metro, desc)
}

func TestAccMetalDatasourceVlan_ByVlanId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByVlanId(rs, metro, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "project_id",
						"data.metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalDatasourceVlanConfig_ByProjectId(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "metal_vlan" "foovlan" {
    project_id = metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "metal_vlan" "dsvlan" {
    project_id = metal_vlan.foovlan.project_id
}
`, projSuffix, metro, desc)
}

func TestAccMetalDatasourceVlan_ByProjectId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalDatasourceVlanDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalDatasourceVlanConfig_ByProjectId(rs, metro, "testvlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "vxlan",
						"data.metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"metal_vlan.foovlan", "project_id",
						"data.metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}
