package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceMetalVlan_byVxlanFacility(t *testing.T) {
	rs := acctest.RandString(10)
	fac := "sv15"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byVxlanFacility(rs, fac, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "id",
						"data.equinix_metal_vlan.dsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byVxlanFacility(projSuffix, fac, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    facility = "%s"
    description = "%s"
}

data "equinix_metal_vlan" "dsvlan" {
    facility = equinix_metal_vlan.foovlan.facility
    project_id = equinix_metal_vlan.foovlan.project_id
    vxlan = equinix_metal_vlan.foovlan.vxlan
}
`, projSuffix, fac, desc)
}

func TestAccDataSourceMetalVlan_byVxlanMetro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byVxlanMetro(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "id",
						"data.equinix_metal_vlan.dsvlan", "id",
					),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.barvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_vlan.bardsvlan", "vxlan", "6",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.barvlan", "id",
						"data.equinix_metal_vlan.bardsvlan", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byVxlanMetro(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
    metro = equinix_metal_vlan.foovlan.metro
    project_id = equinix_metal_vlan.foovlan.project_id
    vxlan = equinix_metal_vlan.foovlan.vxlan
}

resource "equinix_metal_vlan" "barvlan" {
    project_id = equinix_metal_project.foobar.id
    metro = equinix_metal_vlan.foovlan.metro
    description = "%s"
    vxlan = 6
}

data "equinix_metal_vlan" "bardsvlan" {
    metro = equinix_metal_vlan.barvlan.metro
    project_id = equinix_metal_vlan.barvlan.project_id
    vxlan = equinix_metal_vlan.barvlan.vxlan
}
`, projSuffix, metro, desc, desc)
}

func TestAccDataSourceMetalVlan_byVlanId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byVlanId(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "project_id",
						"data.equinix_metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byVlanId(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
    vlan_id = equinix_metal_vlan.foovlan.id
}
`, projSuffix, metro, desc)
}

func TestAccDataSourceMetalVlan_byProjectId(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalDatasourceVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVlanConfig_byProjectId(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "vxlan",
						"data.equinix_metal_vlan.dsvlan", "vxlan",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vlan.foovlan", "project_id",
						"data.equinix_metal_vlan.dsvlan", "project_id",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalVlanConfig_byProjectId(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
    project_id = equinix_metal_vlan.foovlan.project_id
}
`, projSuffix, metro, desc)
}

func testAccMetalDatasourceVlanCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Data source VLAN still exists")
		}
	}

	return nil
}
