package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketBGPSession_Basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketBGPSessionDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketBGPSessionConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "id",
						"packet_bgp_session.test", "device_id"),
					resource.TestCheckResourceAttr(
						"packet_bgp_session.test", "default_route", "true"),
				),
			},
			resource.TestStep{
				ResourceName:      "packet_bgp_session.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketBGPSessionDestroy(s *terraform.State) error {
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

func testAccCheckPacketBGPSessionConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "%s"
	bgp_config {
		deployment_type = "local"
		md5 = "C179c28c41a85b"
		asn = 65000
	}
}

resource "packet_device" "test" {
    hostname         = "terraform-test-bgp-sesh"
    plan             = "t1.small.x86"
    facilities       = ["ewr1"]
    operating_system = "ubuntu_16_04"
    billing_cycle    = "hourly"
    project_id       = "${packet_project.test.id}"
}

resource "packet_bgp_session" "test" {
	device_id = "${packet_device.test.id}"
	address_family = "ipv4"
	default_route = true
}`, name)
}
