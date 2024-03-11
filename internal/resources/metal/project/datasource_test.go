package project_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccDataSourceMetalProject_byId(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalProject_byId(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%s", rn)),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project.foobar", "id",
						"data.equinix_metal_project.test", "id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project.foobar", "organization_id",
						"data.equinix_metal_project.test", "organization_id"),
				),
			},
		},
	})
}

func testAccDataSourceMetalProject_byId(r string) string {
	return fmt.Sprintf(`
terraform {
	provider_meta "equinix" {
		module_name = "test"
	}
}

resource "equinix_metal_project" "foobar" {
	name = "tfacc-project-%s"
	bgp_config {
		deployment_type = "local"
		md5 = "2SFsdfsg43"
		asn = 65000
	}
}

data equinix_metal_project "test" {
	project_id = equinix_metal_project.foobar.id
}
`, r)
}

func TestAccDataSourceMetalProject_byName(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalProject_byName(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%s", rn)),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project.foobar", "id",
						"data.equinix_metal_project.test", "id"),
				),
			},
		},
	})
}

func testAccDataSourceMetalProject_byName(r string) string {
	return fmt.Sprintf(`
terraform {
	provider_meta "equinix" {
		module_name = "test"
	}
}

resource "equinix_metal_project" "foobar" {
	name = "tfacc-project-%s"
	bgp_config {
		deployment_type = "local"
		md5 = "2SFsdfsg43"
		asn = 65000
	}
}

data equinix_metal_project "test" {
	name = equinix_metal_project.foobar.name
}
`, r)
}
