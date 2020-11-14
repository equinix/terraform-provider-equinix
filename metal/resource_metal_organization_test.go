package metal

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
	resource.AddTestSweepers("metal_organization", &resource.Sweeper{
		Name:         "metal_organization",
		Dependencies: []string{"metal_project"},
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
		return fmt.Errorf("Error getting org list for sweeping organizations: %s", err)
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
		CheckDestroy: testAccCheckMetalOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalOrgConfigBasic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalOrgExists("metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"metal_organization.test", "name", fmt.Sprintf("tfacc-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"metal_organization.test", "description", "quux"),
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
		CheckDestroy: testAccCheckMetalOrgDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalOrgConfigBasic(rInt),
			},
			{
				ResourceName:      "metal_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalOrgDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_organization" {
			continue
		}
		if _, _, err := client.Organizations.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Organization still exists")
		}
	}

	return nil
}

func testAccCheckMetalOrgExists(n string, org *packngo.Organization) resource.TestCheckFunc {
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

func testAccCheckMetalOrgConfigBasic(r int) string {
	return fmt.Sprintf(`
resource "metal_organization" "test" {
		name = "tfacc-org-%d"
		description = "quux"
}`, r)
}
