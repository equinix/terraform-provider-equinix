package equinix

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/packethost/packngo"
)

var (
	matchErrOverbid = regexp.MustCompile(".* exceeds the maximum bid price .*")
	matchErrTimeout = regexp.MustCompile(".* timeout while waiting for state to become 'done'.*")
)

func TestAccMetalSpotMarketRequest_basic(t *testing.T) {
	projSuffix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalSpotMarketRequestCheckDestroyed,
		ErrorCheck:               skipIfOverbidOrTimedOut(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMetalSpotMarketRequestConfig_basic(projSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSpotMarketRequestExists("equinix_metal_spot_market_request.request"),
					resource.TestCheckResourceAttr("equinix_metal_spot_market_request.request", "devices_max", "1"),
					resource.TestCheckResourceAttr("equinix_metal_spot_market_request.request", "devices_min", "1"),
					resource.TestCheckResourceAttr("data.equinix_metal_spot_market_request.dreq", "device_ids.#", "1"),
				),
			},
		},
	})
}

func testAccMetalSpotMarketRequestCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).Metal

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

func testAccCheckMetalSpotMarketRequestExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*config.Config).Metal

		foundKey, _, err := client.SpotMarketRequests.Get(rs.Primary.ID, &packngo.GetOptions{Includes: []string{"project", "devices", "facilities", "metro"}})
		if err != nil {
			return err
		}
		if foundKey.ID != rs.Primary.ID {
			return fmt.Errorf("Spot market request not found: %v - %v", rs.Primary.ID, foundKey)
		}

		return nil
	}
}

func testAccMetalSpotMarketRequestConfig_basic(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

data "equinix_metal_spot_market_price" "test" {
  metro = local.metro
  plan  = local.plan
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = equinix_metal_spot_market_request.request.id
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = format("%%.2f", data.equinix_metal_spot_market_price.test.price)
  metro            = data.equinix_metal_spot_market_price.test.metro
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}

func testAccCheckMetalSpotMarketRequestConfig_import(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

data "equinix_metal_spot_market_price" "test" {
  metro = local.metro
  plan  = local.plan
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = format("%%.2f", data.equinix_metal_spot_market_price.test.price)
  metro            = data.equinix_metal_spot_market_price.test.metro
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}

func TestAccMetalSpotMarketRequest_Import(t *testing.T) {
	projSuffix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalSpotMarketRequestCheckDestroyed,
		ErrorCheck:               skipIfOverbidOrTimedOut(t),
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

// In CI, we frequently see errors that the spot market bid price is higher
// than the maximum allowed price.  However, reducing our bid price even by
// 0.01 causes the tests to consistently hit the 30-minute `wait_for_devices`
// timeout because our bid isn't high enough to actually win a device.  Spot
// prices fluctuate by design, so it's impossible to automatically, consistently
// find a bid that will win devices before the timeout _and_ stay below the max
// allowed bid.  This function serves to smooth out acceptance test workflow
// failures by skipping the test after the fact if our bid was too high or too
// low, since the resource is behaving as expected for those scenarios.
func skipIfOverbidOrTimedOut(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		if matchErrOverbid.MatchString(err.Error()) {
			t.Skipf("price was higher than max allowed bid; skipping")
		}
		if matchErrTimeout.MatchString(err.Error()) {
			t.Skipf("timed out waiting for devices (bid was probably too low); skipping")
		}

		return err
	}
}

func testAccMetalSpotMarketRequestConfig_timeout(projSuffix, createTimeout string) string {
	if createTimeout == "" {
		createTimeout = "30m"
	}

	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-spot_market_request-%s"
}

data "equinix_metal_spot_market_price" "test" {
  metro = local.metro
  plan  = local.plan
}

data "equinix_metal_spot_market_request" "dreq" {
  request_id = equinix_metal_spot_market_request.request.id
  timeouts {
    create = "%s"
  }
}

resource "equinix_metal_spot_market_request" "request" {
  project_id       = equinix_metal_project.test.id
  max_bid_price    = format("%%.2f", data.equinix_metal_spot_market_price.test.price)
  metro            = data.equinix_metal_spot_market_price.test.metro
  devices_min      = 1
  devices_max      = 1
  wait_for_devices = true

  instance_parameters {
    hostname         = "tfacc-testspot"
    billing_cycle    = "hourly"
    operating_system = local.os
    plan             = local.plan
  }

  timeouts {
    create = "%s"
  }
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, createTimeout, createTimeout)
}

func TestAccMetalSpotMarketRequestCreate_WithTimeout(t *testing.T) {
	projSuffix := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalSpotMarketRequestCheckDestroyed,
		ErrorCheck: func(err error) error {
			if matchErrOverbid.MatchString(err.Error()) {
				t.Skipf("price was higher than max allowed bid; skipping")
			}
			return err
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalSpotMarketRequestConfig_timeout(projSuffix, "5s"),
				ExpectError: matchErrTimeout,
			},
		},
	})
}
