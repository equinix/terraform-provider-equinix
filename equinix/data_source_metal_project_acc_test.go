package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccDataSourceMetalProject_basic(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalProject_basic(rn),
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

func testAccDataSourceMetalProject_basic(r string) string {
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
