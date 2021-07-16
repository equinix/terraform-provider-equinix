package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccMetalProjectAPIKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_project_api_key" {
			continue
		}
		if _, err := client.APIKeys.ProjectGet(rs.Primary.ID, rs.Primary.Attributes["project_id"], nil); err == nil {
			return fmt.Errorf("ProjectAPI key still exists")
		}
	}
	return nil
}

func testAccMetalProjectAPIKeyConfig_Basic() string {
	return fmt.Sprintf(`

resource "metal_project" "test" {
    name = "tfacc-project-key-test"
}

resource "metal_project_api_key" "test" {
    project_id  = metal_project.test.id
    description = "tfacc-project-key"
    read_only   = true
}`)
}

func TestAccMetalProjectAPIKey_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectAPIKeyConfig_Basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"metal_project_api_key.test", "token"),
					resource.TestCheckResourceAttrPair(
						"metal_project_api_key.test", "project_id",
						"metal_project.test", "id"),
				),
			},
		},
	})
}
