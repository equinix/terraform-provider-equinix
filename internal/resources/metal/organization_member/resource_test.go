package organizationmember_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/packethost/packngo"
)

func TestAccResourceMetalOrganizationMember_owner(t *testing.T) {
	rInt := acctest.RandInt()
	org := &packngo.Organization{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMetalOrganizationMember_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalOrganizationExists("equinix_metal_organization.test", org),
				),
			},
			{
				ResourceName: "equinix_metal_organization_member.owner",
				Config:       testAccResourceMetalOrganizationMember_basic(rInt) + testAccResourceMetalOrganizationMember_owner(),
				ImportStateIdFunc: resource.ImportStateIdFunc(func(s *terraform.State) (string, error) {
					return fmt.Sprintf("%s:%s", org.PrimaryOwner.Email, org.ID), nil
				}),
				ImportState: true,
			},
			{
				Config: testAccResourceMetalOrganizationMember_basic(rInt),
			},
		},
	})
}

func TestAccResourceMetalOrganizationMember_basic(t *testing.T) {
	rInt := acctest.RandInt()
	org := &packngo.Organization{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMetalOrganizationMember_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalOrganizationExists("equinix_metal_organization.test", org),
				),
			},
			{
				ResourceName: "equinix_metal_organization_member.member",
				Config:       testAccResourceMetalOrganizationMember_basic(rInt) + testAccResourceMetalOrganizationMember_member(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_organization_member.member", "state",
						"invited"),
				),
				ImportStateVerify: true,
			},
			{
				Config:  testAccResourceMetalOrganizationMember_basic(rInt),
				Destroy: true,
			},
		},
	})
}

func testAccResourceMetalOrganizationMember_basic(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_organization" "test" {
	name = "tfacc-resource-org-member-%d"
	description = "tfacc-resource-org-member-desc"
	address {
		address = "tfacc org street"
		city = "london"
		zip_code = "12345"
		country = "GB"
	}
}

resource "equinix_metal_project" "test" {
	organization_id = equinix_metal_organization.test.id
	name = "tfacc-resource-project-%d"
}
`, r, r)
}

func testAccResourceMetalOrganizationMember_owner() string {
	return `
	resource "equinix_metal_organization_member" "owner" {
		invitee = "/* TODO: Add org owner email or token owner email here */"
		roles = ["owner"]
		projects_ids = []
		organization_id = equinix_metal_organization.test.id
	}
	`
}

func testAccResourceMetalOrganizationMember_member() string {
	return `
resource "equinix_metal_organization_member" "member" {
    invitee = "tfacc.testing.member@equinixmetal.com"
	roles = ["limited_collaborator"]
    projects_ids = [equinix_metal_project.test.id]
    organization_id = equinix_metal_organization.test.id
	message = "This invitation was sent by the github.com/equinix/terraform-provider-equinix acceptance tests to test equinix_metal_organization_member resources."
}
`
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
