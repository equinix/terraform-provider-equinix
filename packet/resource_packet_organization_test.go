package packet

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("packet_organization", &resource.Sweeper{
		Name:         "packet_organization",
		Dependencies: []string{"packet_project"},
		F:            testSweepOrganizations,
	})
}

func testSweepOrganizations(region string) error {
	log.Printf("[DEBUG] Sweeping organizations")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping organizations: %s", err)
	}
	client := meta.(*packngo.Client)

	os, _, err := client.Organizations.List(nil)
	if err != nil {
		return fmt.Errorf("Error getting org list for sweepeing organizations: %s", err)
	}
	oids := []string{}
	for _, o := range os {
		if strings.HasPrefix(o.Name, "tfacc-") {
			oids = append(oids, o.ID)
		}
	}
	for _, oid := range oids {
		log.Printf("Removing organization %s", oid)
		_, err := client.Organizations.Delete(oid)
		if err != nil {
			return fmt.Errorf("Error deleting organization %s", err)
		}
	}
	return nil
}

func TestAccOrgCreate(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketOrgConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketOrgExists("packet_organization.test", &org),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "name", fmt.Sprintf("tfacc-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"packet_organization.test", "description", "quux"),
				),
			},
		},
	})
}

func TestAccOrg_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketOrgConfigBasic(rInt),
			},
			{
				ResourceName:      "packet_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
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

func testAccCheckPacketOrgConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "packet_organization" "test" {
		name = "tfacc-org-%d"
		description = "quux"
}`, r)
}
