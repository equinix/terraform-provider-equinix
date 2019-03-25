package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketSpotMarketRequest_Basic(t *testing.T) {
	var key packngo.SpotMarketRequest

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketSpotMarketRequestConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketSpotMarketRequestExists("packet_spot_market_request.request", &key),
					resource.TestCheckResourceAttr("packet_spot_market_request.request", "devices_max", "1"),
					resource.TestCheckResourceAttr("packet_spot_market_request.request", "devices_min", "1"),
					resource.TestCheckResourceAttr("packet_spot_market_request.request", "max_bid_price", "0.03"),
				),
			},
		},
	})
}

func testAccCheckPacketSpotMarketRequestDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_spot_market_request" {
			continue
		}
		if _, _, err := client.SpotMarketRequests.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Spot market request key still exists")
		}
	}

	return nil
}

func testAccCheckPacketSpotMarketRequestExists(n string, key *packngo.SpotMarketRequest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundKey, _, err := client.SpotMarketRequests.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundKey.ID != rs.Primary.ID {
			return fmt.Errorf("Spot market request not found: %v - %v", rs.Primary.ID, foundKey)
		}

		*key = *foundKey

		return nil
	}
}

var testAccCheckPacketSpotMarketRequestConfig_basic = `
	  resource "packet_project" "test" {
		name = "TerraformTestProject-SMR"
	  }
	  
	  resource "packet_spot_market_request" "request" {
		project_id       = "${packet_project.test.id}"
		max_bid_price    = 0.03
		facilities       = ["ewr1"]
		devices_min      = 1
		devices_max      = 1
		wait_for_devices = true

		instance_parameters {
		  hostname         = "testspot"
		  billing_cycle    = "hourly"
		  operating_system = "coreos_stable"
		  plan             = "t1.small.x86"
		}
	  }`
