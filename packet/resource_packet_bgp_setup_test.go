package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketBGPSetup_Basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketBGPSetupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketBGPSetupConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "id",
						"packet_bgp_session.test4", "device_id"),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "id",
						"packet_bgp_session.test6", "device_id"),
					resource.TestCheckResourceAttr(
						"packet_bgp_session.test4", "default_route", "true"),
					resource.TestCheckResourceAttr(
						"packet_bgp_session.test6", "default_route", "true"),
					resource.TestCheckResourceAttr(
						"packet_bgp_session.test4", "address_family", "ipv4"),
					resource.TestCheckResourceAttr(
						"packet_bgp_session.test6", "address_family", "ipv6"),
					// there will be 2 BGP neighbors, for IPv4 and IPv6
					resource.TestCheckResourceAttr(
						"data.packet_device_bgp_neighbors.test", "bgp_neighbors.#", "2"),
				),
			},
			{
				ResourceName:      "packet_bgp_session.test4",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketBGPSetupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_bgp_session" {
			continue
		}
		if _, _, err := client.BGPSessions.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("BGPSession still exists")
		}
	}

	return nil
}

func testAccCheckPacketBGPSetupConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-bgp_session-%s"
	bgp_config {
		deployment_type = "local"
		md5 = "C179c28c41a85b"
		asn = 65000
	}
}

resource "packet_device" "test" {
    hostname         = "tfacc-test-bgp-sesh"
    plan             = "t1.small.x86"
    facilities       = ["ewr1"]
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = "${packet_project.test.id}"
}

resource "packet_bgp_session" "test4" {
	device_id = "${packet_device.test.id}"
	address_family = "ipv4"
	default_route = true
}

resource "packet_bgp_session" "test6" {
	device_id = "${packet_device.test.id}"
	address_family = "ipv6"
	default_route = true
}

data "packet_device_bgp_neighbors" "test" {
  device_id  = packet_bgp_session.test4.device_id
}
`, name)
}
