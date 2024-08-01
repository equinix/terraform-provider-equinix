package vrfbgpdynamicneighbor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetalVrfBgpDynamicNeighbor_basic(t *testing.T) {
	rs := acctest.RandString(10)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVrfBgpDynamicNeighborConfig(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_vrf_bgp_dynamic_neighbor.test", "gateway_id",
						"equinix_metal_gateway.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf_bgp_dynamic_neighbor.test", "range",
						"2001:d78:0:0:4000::/66"),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf_bgp_dynamic_neighbor.test", "asn",
						"56789"),
				),
			},
		},
	})
}

func testAccMetalVrfBgpDynamicNeighborConfig(projSuffix string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-vrf-bgp-neighbor-test-%s"
}

resource "equinix_metal_vlan" "test" {
    description = "tfacc-vlan VLAN in SV"
    metro       = "sv"
    project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_vrf" "test" {
  description = "tfacc-vrf VRF in SV"
  name        = "tfacc-vrf-%s"
  metro       = "sv"
  local_asn   = "65000"
  ip_ranges   = ["2001:d78::/59"]
  bgp_dynamic_neighbors_enabled = true

  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = equinix_metal_project.test.id
  type       = "vrf"
  vrf_id     = equinix_metal_vrf.test.id
  network    = "2001:d78::"
  metro      = "sv"
  cidr       = 64
}

resource "equinix_metal_gateway" "test" {
    project_id        = equinix_metal_project.test.id
    vlan_id           = equinix_metal_vlan.test.id
    ip_reservation_id = equinix_metal_reserved_ip_block.test.id
}

resource "equinix_metal_vrf_bgp_dynamic_neighbor" "test" {
    gateway_id = equinix_metal_gateway.test.id
	range      = "2001:d78:0:0:4000::/66"
	asn        = "56789"
}
`, projSuffix, projSuffix)
}

func testAccCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vrf_bgp_dynamic_gateway" {
			continue
		}
		if _, _, err := client.VRFsApi.BgpDynamicNeighborsIdGet(context.Background(), rs.Primary.ID).Execute(); err == nil {
			return fmt.Errorf("Metal VRF BGP dynamic neighbor still exists")
		}
	}

	return nil
}
