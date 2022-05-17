package spotmarketrequest

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/tfacc"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func TestAccMetalSpotMarketRequest_basic(t *testing.T) {
	var key packngo.SpotMarketRequest
	projSuffix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { tfacc.PreCheck(t) },
		Providers:    tfacc.AccProviders,
		CheckDestroy: testAccMetalSpotMarketRequestCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalSpotMarketRequestConfig_basic(projSuffix),
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
	client := tfacc.AccProvider.Meta().(*config.Config).MetalClient

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

		client := tfacc.AccProvider.Meta().(*config.Config).MetalClient

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

func testAccMetalSpotMarketRequestConfig_basic(projSuffix string) string {
	return tfacc.NPrintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

data "equinix_metal_spot_market_price" "test" {
  facility = local.facility
  plan     = local.plan
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = equinix_metal_spot_market_request.request.id
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = data.equinix_metal_spot_market_price.test.price * 1.2
  facilities       = [data.equinix_metal_spot_market_price.test.facility]
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }

  lifecycle {
    ignore_changes = [
      instance_parameters,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}

func testAccCheckMetalSpotMarketRequestConfig_import(projSuffix string) string {
	return tfacc.NPrintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

data "equinix_metal_spot_market_price" "test" {
  facility = local.facility
  plan     = local.plan
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = data.equinix_metal_spot_market_price.test.price * 1.2
  facilities       = [data.equinix_metal_spot_market_price.test.facility]
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }

  lifecycle {
    ignore_changes = [
      instance_parameters,
      facilities,
    ]
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}

func TestAccMetalSpotMarketRequest_Import(t *testing.T) {
	projSuffix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { tfacc.PreCheck(t) },
		Providers:    tfacc.AccProviders,
		CheckDestroy: testAccMetalSpotMarketRequestCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSpotMarketRequestConfig_import(projSuffix),
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
