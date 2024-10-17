package organization_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/packethost/packngo"
)

func TestAccDataSourceOrganizations_basic(t *testing.T) {
	var org packngo.Organization
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalOrganizationCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalOrganizationConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalOrganizationExists("equinix_metal_organization.test", &org),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_organization.test", "description", "quux"),
					resource.TestCheckResourceAttr(
						"data.equinix_metal_organization.test", "name",
						fmt.Sprintf("tfacc-datasource-org-%d", rInt)),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_organization.test", "address.0.address",
						"data.equinix_metal_organization.test", "address.0.address",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_organization.test", "address.0.city",
						"data.equinix_metal_organization.test", "address.0.city",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_organization.test", "address.0.country",
						"data.equinix_metal_organization.test", "address.0.country",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_organization.test", "address.0.zip_code",
						"data.equinix_metal_organization.test", "address.0.zip_code",
					),
				),
			},
		},
	})
}

func testAccDataSourceMetalOrganizationConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_organization" "test" {
	name = "tfacc-datasource-org-%d"
	description = "quux"
	address {
		address = "tfacc org street"
		city = "london"
		zip_code = "12345"
		country = "GB"
	}
}

data "equinix_metal_organization" "test" {
  organization_id = equinix_metal_organization.test.id
}
`, r)
}
