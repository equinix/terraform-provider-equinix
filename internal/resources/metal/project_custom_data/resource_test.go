package project_custom_data_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetalProjectCustomData_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "equinix_metal_project_custom_data.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalProjectCustomDataCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectCustomDataConfig(rInt, `{"owner":"platform","env":"test"}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "custom_data", `{"env":"test","owner":"platform"}`),
					resource.TestCheckResourceAttrPair(resourceName, "project_id", "equinix_metal_project.test", "id"),
				),
			},
			{
				Config: testAccMetalProjectCustomDataConfig(rInt, `{"owner":"platform","env":"prod"}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "custom_data", `{"env":"prod","owner":"platform"}`),
					resource.TestCheckResourceAttrPair(resourceName, "project_id", "equinix_metal_project.test", "id"),
				),
			},
			{
				Config: testAccMetalProjectCustomDataProjectOnlyConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectCustomDataProjectCleared("equinix_metal_project.test"),
				),
			},
		},
	})
}

func testAccMetalProjectCustomDataProjectCleared(projectResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

		rs, ok := s.RootModule().Resources[projectResource]
		if !ok {
			return fmt.Errorf("project resource %q not found in state", projectResource)
		}

		project, _, err := client.ProjectsApi.FindProjectById(context.Background(), rs.Primary.ID).Execute()
		if err != nil {
			return fmt.Errorf("could not read project %q: %w", rs.Primary.ID, err)
		}

		if len(project.GetCustomdata()) != 0 {
			return fmt.Errorf("project custom data still exists")
		}

		return nil
	}
}

func testAccMetalProjectCustomDataCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_project_custom_data" {
			continue
		}

		project, _, err := client.ProjectsApi.FindProjectById(context.Background(), rs.Primary.Attributes["project_id"]).Execute()
		if err != nil {
			continue
		}

		if len(project.GetCustomdata()) != 0 {
			return fmt.Errorf("project custom data still exists")
		}
	}

	return nil
}

func testAccMetalProjectCustomDataConfig(r int, customData string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-project-custom-data-%d"
}

resource "equinix_metal_project_custom_data" "test" {
    project_id  = equinix_metal_project.test.id
    custom_data = <<JSON
%s
JSON
}
`, r, customData)
}

func testAccMetalProjectCustomDataProjectOnlyConfig(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-project-custom-data-%d"
}
`, r)
}
