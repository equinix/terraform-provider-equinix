package equinix

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("equinix_metal_organization", &resource.Sweeper{
		Name:         "equinix_metal_organization",
		Dependencies: []string{"equinix_metal_project"},
		F:            testSweepOrganizations,
	})
}

func testSweepOrganizations(region string) error {
	log.Printf("[DEBUG] Sweeping organizations")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping organizations: %s", err)
	}
	client := meta.Client()

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

func TestAccMetalOrganizationCreate(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalOrganizationConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalOrganizationExists("equinix_metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "name", fmt.Sprintf("tfacc-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "description", "quux"),
				),
			},
		},
	})
}

func TestAccMetalAccOrganization_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalOrganizationConfig_basic(rInt),
			},
			{
				ResourceName:      "equinix_metal_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalOrganizationCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_organization" {
			continue
		}
		if _, _, err := client.Organizations.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Organization still exists")
		}
	}

	return nil
}

func testAccMetalOrganizationExists(n string, org *packngo.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Config).Client()

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

func testAccMetalOrganizationConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_organization" "test" {
		name = "tfacc-org-%d"
		description = "quux"
}`, r)
}
