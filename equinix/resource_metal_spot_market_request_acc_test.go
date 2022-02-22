package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalSpotMarketRequest_basic(t *testing.T) {
	var key packngo.SpotMarketRequest
	context := map[string]interface{}{
		"name_suffix": acctest.RandString(10),
		"facility":    "dc13",
		"plan":        "c3.small.x86",
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalSpotMarketRequestCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalSpotMarketRequestConfig_basic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("equinix_metal_spot_market_request.request", &key),
					resource.TestCheckResourceAttr("equinix_metal_spot_market_request.request", "devices_max", "1"),
					resource.TestCheckResourceAttr("equinix_metal_spot_market_request.request", "devices_min", "1"),
					resource.TestCheckResourceAttr("data.equinix_metal_spot_market_request.dreq", "device_ids.#", "1"),
				),
			},
		},
	})
}

func testAccMetalSpotMarketRequestCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_spot_market_request" {
			continue
		}
		if _, _, err := client.SpotMarketRequests.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Spot market request key still exists")
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

		client := testAccProvider.Meta().(*Config).metal

		foundKey, _, err := client.SpotMarketRequests.Get(rs.Primary.ID, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
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

func testAccMetalSpotMarketRequestConfig_basic(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%{name_suffix}"
}

data "equinix_metal_spot_market_price" "test" {
  facility = "%{facility}"
  plan     = "%{plan}"
}

data "equinix_metal_spot_market_request" "dreq" {
	request_id = equinix_metal_spot_market_request.request.id
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = data.equinix_metal_spot_market_price.test.price * 1.2
  facilities       = ["%{facility}"]
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_18_04"
    plan             = "%{plan}"
  }
}`, ctx)
}

func testAccCheckMetalSpotMarketRequestConfig_import(ctx map[string]interface{}) string {
	return nprintf(`
resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%{name_suffix}"
}

data "equinix_metal_spot_market_price" "test" {
  facility = "%{facility}"
  plan     = "%{plan}"
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = data.equinix_metal_spot_market_price.test.price * 1.2
  facilities       = ["%{facility}"]
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = "ubuntu_20_04"
    plan             = "%{plan}"
  }
}`, ctx)
}

func TestAccMetalSpotMarketRequest_Import(t *testing.T) {
	context := map[string]interface{}{
		"name_suffix": acctest.RandString(10),
		"facility":    "dc13",
		"plan":        "c3.small.x86",
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSpotMarketRequestConfig_import(context),
			},
			{
				ResourceName:            "equinix_metal_spot_market_request.request",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_parameters", "wait_for_devices"},
			},
		},
	})
}
