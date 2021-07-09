package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func metalProjectAPIKeyConfig_Basic(description string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-project_api_key-%s"
}

resource "metal_project_api_key" "test-ro" {
    description = "tfacc-project-api-key-ro-test"
    project_id = "${metal_project.test.id}"
}

resource "metal_project_api_key" "test-rw" {
    description = "tfacc-project-api-key-rw-test"
	read_only = false
    project_id = "${metal_project.test.id}"
}
`, description)
}

func TestAccMetalProjectAPIKey_Basic(t *testing.T) {
	rs := acctest.RandString(10)
	var key packngo.APIKey

	cfg := metalProjectAPIKeyConfig_Basic(rs)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectAPIKeyExists("metal_project_api_key.test-ro", &key),
					resource.TestCheckResourceAttr(
						"metal_project_api_key.test", "read_only", "true"),
					resource.TestCheckResourceAttr(
						"metal_project_api_key.test", "description", "tfacc-project-api-key-ro-test"),
				),
			},
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectAPIKeyExists("metal_project_api_key.test-rw", &key),
					resource.TestCheckResourceAttr(
						"metal_project_api_key.test", "read_only", "false"),
					resource.TestCheckResourceAttr(
						"metal_project_api_key.test", "description", "tfacc-project-api-key-rw-test"),
				),
			},
		},
	})
}

func testAccCheckMetalProjectAPIKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_project_api_key" {
			continue
		}

		if _, err := client.APIKeys.ProjectGet(rs.Primary.Attributes["project_id"], rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Project API key still exists")
		}
	}

	return nil
}

func testAccCheckMetalProjectAPIKeyExists(n string, key *packngo.APIKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundAPIKey, err := client.APIKeys.ProjectGet(key.Project.ID, rs.Primary.ID, nil)
		if err != nil {
			return err
		}

		if foundAPIKey.ID != rs.Primary.ID {
			return fmt.Errorf("Project API Key not found: %v - %v", rs.Primary.ID, foundAPIKey.ID)
		}

		*key = *foundAPIKey

		return nil
	}
}
