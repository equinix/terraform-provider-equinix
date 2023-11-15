package equinix

import (
	"fmt"
	"log"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

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
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping organizations: %s", err)
	}
	metal := config.NewMetalClient()
	os, _, err := metal.Organizations.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting org list for sweeping organizations: %s", err)
	}
	oids := []string{}
	for _, o := range os {
		if isSweepableTestResource(o.Name) {
			oids = append(oids, o.ID)
		}
	}
	for _, oid := range oids {
		log.Printf("Removing organization %s", oid)
		_, err := metal.Organizations.Delete(oid)
		if err != nil {
			return fmt.Errorf("Error deleting organization %s", err)
		}
	}
	return nil
}

func TestAccMetalOrganization_create(t *testing.T) {
	var org, org2 packngo.Organization
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalOrganizationConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalOrganizationExists("equinix_metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "name", fmt.Sprintf("tfacc-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "address.0.city", "London"),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "address.0.state", ""),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "address.0.zip_code", "12345"),
				),
			},
			{
				Config: testAccMetalOrganizationConfig_basicUpdate(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalOrganizationExists("equinix_metal_organization.test", &org2),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "name", fmt.Sprintf("tfacc-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "description", "baz"),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "address.0.city", "Madrid"),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "address.0.state", "Madrid"),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "twitter", "@Equinix"),
					testAccMetalSameOrganization(t, &org, &org2),
				),
			},
		},
	})
}

func testAccMetalSameOrganization(t *testing.T, before, after *packngo.Organization) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected organization to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func TestAccMetalOrganization_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalOrganizationCheckDestroyed,
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
	client := testAccProvider.Meta().(*config.Config).Metal

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

		client := testAccProvider.Meta().(*config.Config).Metal

		foundOrg, _, err := client.Organizations.Get(rs.Primary.ID, &packngo.GetOptions{Includes: []string{"address", "primary_owner"}})
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
	address {
		address = "tfacc org street"
		city = "London"
		zip_code = "12345"
		country = "GB"
	}
}`, r)
}

func testAccMetalOrganizationConfig_basicUpdate(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_organization" "test" {
	name = "tfacc-org-%d"
	description = "baz"
	address {
		address = "tfacc org street"
		city = "Madrid"
		zip_code = "28108"
		country = "ES"
		state   = "Madrid"
	}
	twitter = "@Equinix"
}`, r)
}
