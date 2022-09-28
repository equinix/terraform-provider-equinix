package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMetalOrganizationMember_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// TODO: CheckDestroy: testAccMetalOrganizationMemberCheckDestroyed,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMetalOrganizationMember_basic(rInt) + testAccResourceMetalOrganizationMember_member(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_organization_member.member", "state",
						"invited"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_organization_member.member", "id",
						"equinix_metal_organization_member.member", "invitee",
					),
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
	name = "tfacc-datasource-org-%d"
	description = "tfacc-datasource-org-desc"
	address {
		address = "tfacc org street"
		city = "london"
		zip_code = "12345"
		country = "GB"
	}
}

resource "equinix_metal_project" "test" {
	organization_id = equinix_metal_organization.test.id
	name = "tfacc-datasource-project-%d"
}
`, r, r)
}

func testAccResourceMetalOrganizationMember_member() string {
	return `
resource "equinix_metal_organization_member" "member" {
    invitee = "tfacc.testing.member@equinixmetal.com"
    roles = ["limited_collaborator"]
    projects_ids = [equinix_metal_project.test.id]
    organization_id = equinix_metal_organization.test.id
}
`
}
