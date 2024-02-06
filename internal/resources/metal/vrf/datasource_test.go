package vrf_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccDataSourceMetalVrfDataSource_byID(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	datasourceName := "data.equinix_metal_vrf.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:         acceptance.TestExternalProviders,
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVrfDataSourceConfig_byID(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", datasourceName),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_vrf.foobar", "local_asn"),
				),
				// Why was follwing flag set? The plan is applied and then it's empty.
				// It's causing errors in acceptance tests. Was this because of some API bug?
				// ExpectNonEmptyPlan: true,
			},
			{
				Config:      testAccDataSourceMetalVrfDataSourceConfig_byID(rInt),
				ExpectError: regexp.MustCompile("was not found"),
			},
			{
				// Exit the tests with an empty state and a valid config
				// following the previous error config. This is needed for the
				// destroy step to succeed.
				Config: `/* this config intentionally left blank */`,
			},
		},
	})
}

func testAccDataSourceMetalVrfDataSourceConfig_byID(r int) string {
	testMetro := "da"

	config := fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-vrfs-%d"
}

resource "equinix_metal_vrf" "test" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	project_id = "${equinix_metal_project.test.id}"
}

data "equinix_metal_vrf" "foobar" {
	vrf_id = equinix_metal_vrf.test.id
}`, r, r, testMetro)

	return config
}

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccDataSourceMetalVrf_upgradeFromVersion(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	datasourceName := "data.equinix_metal_vrf.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		PreventPostDestroyRefresh: true,
		ExternalProviders:         acceptance.TestExternalProviders,
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		CheckDestroy:              testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.24.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: testAccDataSourceMetalVrfDataSourceConfig_byID(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", datasourceName),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_vrf.foobar", "local_asn"),
				),
			},
			{
				ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
				Config:                   testAccDataSourceMetalVrfDataSourceConfig_byID(rInt),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
