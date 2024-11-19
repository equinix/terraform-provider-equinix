package organization_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalOrganization_create(t *testing.T) {
	var org, org2 packngo.Organization
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalOrganizationCheckDestroyed,
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
				PreConfig: testAccMetalWaitForOrganization,
				Config:    testAccMetalOrganizationConfig_basicUpdate(rInt),
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

func TestAccMetalOrganization_importBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalOrganizationCheckDestroyed,
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

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
func TestAccMetalOrganization_upgradeFromVersion(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		CheckDestroy: testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.29.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
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
				ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
				Config:                   testAccMetalOrganizationConfig_basic(rInt),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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

func testAccMetalWaitForOrganization() {
	// Some aspect of organization creation takes a while
	// to propagate; updating an organization too soon after
	// create causes test failures and probably doesn't
	// reflect real-world usage.
	time.Sleep(5 * time.Minute)
}

func testAccMetalOrganizationCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

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

		client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

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
