package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/packethost/packngo"
)

func TestAccResourceMetalOrganizationMember_owner(t *testing.T) {
	rInt := acctest.RandInt()
	org := &packngo.Organization{}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		// TODO: CheckDestroy: testAccMetalOrganizationMemberCheckDestroyed,
		CheckDestroy: nil,
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
			/*
				{
					ResourceName: "equinix_metal_organization_member.owner",
					Config:       testAccResourceMetalOrganizationMember_basic(rInt) + testAccResourceMetalOrganizationMember_owner(),
					ExpectError:  regexp.MustCompile("User is already a member of the Organization"),
				},
			*/
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		// TODO: CheckDestroy: testAccMetalOrganizationMemberCheckDestroyed,
		CheckDestroy: nil,
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
	return fmt.Sprintf(`
	resource "equinix_metal_organization_member" "owner" {
		invitee = "/* TODO: Add org owner email or token owner email here */"
		roles = ["owner"]
		projects_ids = []
		organization_id = equinix_metal_organization.test.id
	}
	`)
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
