package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccOrgCreate(t *testing.T) {
	var org packngo.Organization

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketOrgDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketOrgConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketOrgExists("packet_organization.test", &org),
					testAccCheckPacketOrgAttributes(&org),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "name", "foobar"),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "description", "quux"),
				),
			},
		},
	})
}

func testAccCheckPacketOrgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_organization" {
			continue
		}
		if _, _, err := client.Organizations.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Organization still exists")
		}
	}

	return nil
}

func testAccCheckPacketOrgAttributes(org *packngo.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if org.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", org.Name)
		}
		return nil
	}
}

func testAccCheckPacketOrgExists(n string, org *packngo.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundOrg, _, err := client.Organizations.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundOrg.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundOrg)
		}

		*org = *foundOrg

		return nil
	}
}

var testAccCheckPacketOrgConfigBasic = fmt.Sprintf(`
resource "packet_organization" "test" {
		name = "foobar"
		description = "quux"
}`)
