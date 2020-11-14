package metal

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalSpotMarketRequest_Basic(t *testing.T) {
	var key packngo.SpotMarketRequest
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSpotMarketRequestConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("metal_spot_market_request.request", &key),
					resource.TestCheckResourceAttr("metal_spot_market_request.request", "devices_max", "1"),
					resource.TestCheckResourceAttr("metal_spot_market_request.request", "devices_min", "1"),
				),
			},
		},
	})
}

func testAccCheckMetalSpotMarketRequestDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_spot_market_request" {
			continue
		}
		if _, _, err := client.SpotMarketRequests.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Spot market request key still exists")
		}
	}

	return nil
}

func testAccCheckMetalSpotMarketRequestExists(n string, key *packngo.SpotMarketRequest) resource.TestCheckFunc {
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

func testAccCheckMetalSpotMarketRequestConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

data "metal_spot_market_price" "test" {
  facility = "ewr1"
  plan     = "baremetal_0"
}


resource "metal_spot_market_request" "request" {
  project_id       = "${metal_project.test.id}"
  max_bid_price    = data.metal_spot_market_price.test.price * 1.2
  facilities       = ["ewr1"]
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_18_04"
    plan             = "baremetal_0"
  }
}`, name)
}
