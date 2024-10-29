package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetalProjectAPIKey_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalProjectAPIKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectAPIKeyConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"equinix_metal_project_api_key.test", "token"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project_api_key.test", "project_id",
						"equinix_metal_project.test", "id"),
				),
			},
		},
	})
}

func testAccMetalProjectAPIKeyConfig_basic() string {
	return `

resource "equinix_metal_project" "test" {
    name = "tfacc-project-key-test"
}

resource "equinix_metal_project_api_key" "test" {
    project_id  = equinix_metal_project.test.id
    description = "tfacc-project-key"
    read_only   = true
}`
}

func testAccMetalProjectAPIKeyCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).Metal
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_project_api_key" {
			continue
		}
		if _, err := client.APIKeys.ProjectGet(rs.Primary.ID, rs.Primary.Attributes["project_id"], nil); err == nil {
			return fmt.Errorf("Metal ProjectAPI key still exists")
		}
	}
	return nil
}
