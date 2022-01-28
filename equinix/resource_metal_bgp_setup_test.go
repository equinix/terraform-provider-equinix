package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalBGPSetup_Basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalBGPSetupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalBGPSetupConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "id",
						"equinix_metal_bgp_session.test4", "device_id"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "id",
						"equinix_metal_bgp_session.test6", "device_id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_bgp_session.test4", "default_route", "true"),
					resource.TestCheckResourceAttr(
						"equinix_metal_bgp_session.test6", "default_route", "true"),
					resource.TestCheckResourceAttr(
						"equinix_metal_bgp_session.test4", "address_family", "ipv4"),
					resource.TestCheckResourceAttr(
						"equinix_metal_bgp_session.test6", "address_family", "ipv6"),
					// there will be 2 BGP neighbors, for IPv4 and IPv6
					resource.TestCheckResourceAttr(
						"data.equinix_metal_device_bgp_neighbors.test", "bgp_neighbors.#", "2"),
				),
			},
			{
				ResourceName:      "equinix_metal_bgp_session.test4",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMetalBGPSetupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_bgp_session" {
			continue
		}
		if _, _, err := client.BGPSessions.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("BGPSession still exists")
		}
	}

	return nil
}

func testAccCheckMetalBGPSetupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-bgp_session-%s"
	bgp_config {
		deployment_type = "local"
		md5 = "C179c28c41a85b"
		asn = 65000
	}
}

resource "equinix_metal_device" "test" {
    hostname         = "tfacc-test-bgp-sesh"
    plan             = "t1.small.x86"
    facilities       = ["ewr1"]
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = "${metal_project.test.id}"
}

resource "equinix_metal_bgp_session" "test4" {
	device_id = "${metal_device.test.id}"
	address_family = "ipv4"
	default_route = true
}

resource "equinix_metal_bgp_session" "test6" {
	device_id = "${metal_device.test.id}"
	address_family = "ipv6"
	default_route = true
}

data "equinix_metal_device_bgp_neighbors" "test" {
  device_id  = metal_bgp_session.test4.device_id
}
`, name)
}
