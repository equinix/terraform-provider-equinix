package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMetalConnection_shared(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_shared(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "metro", "sv"),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalConnectionConfig_shared(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "equinix_metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = equinix_metal_project.test.organization_id
            project_id      = equinix_metal_project.test.id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "shared"
        }`,
		randstr, randstr)
}

func testAccMetalConnectionCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_connection" {
			continue
		}
		if _, _, err := client.Connections.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Connection still exists")
		}
	}

	return nil
}

func TestAccMetalConnection_dedicated(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_dedicated(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "type", "dedicated"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "redundancy", "redundant"),
					resource.TestCheckResourceAttr("equinix_metal_connection.test", "metro", "sv"),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalConnectionConfig_dedicated(rs) + testDataSourceMetalConnectionConfig_dedicated(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "metro", "sv"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "mode", "standard"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "type", "dedicated"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "redundancy", "redundant"),
					resource.TestCheckResourceAttr("data.equinix_metal_connection.test", "metro", "sv"),
				),
			},
		},
	})
}

func testAccMetalConnectionConfig_dedicated(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }
        
        // No project ID. We only use the project resource to get org_id
        resource "equinix_metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = equinix_metal_project.test.organization_id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "dedicated"
			tags            = ["tfacc"]
			mode            = "standard"
        }`,
		randstr, randstr)
}

func testDataSourceMetalConnectionConfig_dedicated() string {
	return `
		data "equinix_metal_connection" "test" {
            connection_id = equinix_metal_connection.test.id
        }`
}

func TestAccMetalConnection_tunnel(t *testing.T) {

	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalConnectionCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_tunnel(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_connection.test", "mode", "tunnel"),
				),
			},
			{
				ResourceName:      "equinix_metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalConnectionConfig_tunnel(randstr string) string {
	return fmt.Sprintf(`
        resource "equinix_metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "equinix_metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = equinix_metal_project.test.organization_id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "dedicated"
            mode            = "tunnel"
        }`,
		randstr, randstr)
}
