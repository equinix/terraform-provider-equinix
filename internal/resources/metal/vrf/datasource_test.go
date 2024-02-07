package vrf_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceMetalVrfDataSource_byID(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	datasourceKey := "data.equinix_metal_vrf.test"
	name := "tfacc-vrf-" + strconv.Itoa(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		PreventPostDestroyRefresh: true,
		ExternalProviders:         acceptance.TestExternalProviders,
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		CheckDestroy:              testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalVrfDataSourceConfig_byID(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						datasourceKey, "name", name),
					resource.TestCheckResourceAttrSet(
						datasourceKey, "vrf_id"),
				),
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
	local_asn = "65000"
	ip_ranges = ["192.168.100.0/25"]
	project_id = equinix_metal_project.test.id
}

data "equinix_metal_vrf" "test" {
	vrf_id = equinix_metal_vrf.test.id
}`, r, r, testMetro)

	return config
}
