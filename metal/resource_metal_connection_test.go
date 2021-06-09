package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckMetalConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_connection" {
			continue
		}
		if _, _, err := client.Connections.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Connection still exists")
		}
	}

	return nil
}

func testAccMetalConnectionConfig_Shared(randstr string) string {
	return fmt.Sprintf(`
        resource "metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }

        resource "metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = metal_project.test.organization_id
            project_id      = metal_project.test.id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "shared"
        }`,
		randstr, randstr)
}

func TestAccMetalConnection_Shared(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Shared(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_connection.test", "metro", "sv"),
				),
			},
			{
				ResourceName:      "metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalConnectionConfig_Dedicated(randstr string) string {
	return fmt.Sprintf(`
        resource "metal_project" "test" {
            name = "tfacc-conn-pro-%s"
        }
        
        // No project ID. We only use the project resource to get org_id
        resource "metal_connection" "test" {
            name            = "tfacc-conn-%s"
            organization_id = metal_project.test.organization_id
            metro           = "sv"
            redundancy      = "redundant"
            type            = "dedicated"
        }`,
		randstr, randstr)
}

func TestAccMetalConnection_Dedicated(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalConnectionConfig_Dedicated(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_connection.test", "metro", "sv"),
				),
			},
			{
				ResourceName:      "metal_connection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
