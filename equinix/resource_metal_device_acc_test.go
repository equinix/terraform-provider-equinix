package equinix

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

// list of plans and metros and os used as filter criteria to find available hardware to run tests
var (
	preferable_plans  = []string{"x1.small.x86", "t1.small.x86", "c2.medium.x86", "c3.small.x86", "c3.medium.x86", "m3.small.x86"}
	preferable_metros = []string{"ch", "ny", "sv", "ty", "am"}
	preferable_os     = []string{"ubuntu_20_04"}
)

func init() {
	resource.AddTestSweepers("equinix_metal_device", &resource.Sweeper{
		Name: "equinix_metal_device",
		F:    testSweepDevices,
	})
}

func testSweepDevices(region string) error {
	log.Printf("[DEBUG] Sweeping devices")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping devices: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweepeing devices: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if isSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.Devices.List(pid, nil)
		if err != nil {
			log.Printf("Error listing devices to sweep: %s", err)
			continue
		}
		for _, d := range ds {
			if isSweepableTestResource(d.Hostname) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing device %s", did)
		_, err := metal.Devices.Delete(did, true)
		if err != nil {
			return fmt.Errorf("Error deleting device %s", err)
		}
	}
	return nil
}

// Regexp vars for use with resource.ExpectError
var (
	matchErrMustBeProvided    = regexp.MustCompile(".* must be provided when .*")
	matchErrShouldNotBeAnIPXE = regexp.MustCompile(`.*"user_data" should not be an iPXE.*`)
)

// This function should be used to find available plans in all test where a metal_device resource is needed.
//
// TODO consider adding a datasource for equinix_metal_operating_system and making the local.os conditional
//
//	https://github.com/equinix/terraform-provider-equinix/pull/220#discussion_r915418418equinix_metal_operating_system
//	https://github.com/equinix/terraform-provider-equinix/discussions/221
func confAccMetalDevice_base(plans, metros, os []string) string {
	return fmt.Sprintf(`
data "equinix_metal_plans" "test" {
    sort {
        attribute = "id"
        direction = "asc"
    }

    filter {
        attribute = "name"
        values    = [%s]
    }
    filter {
        attribute = "available_in_metros"
        values    = [%s]
    }
    filter {
        attribute = "deployment_types"
        values    = ["on_demand", "spot_market"]
    }
}

// Select a metal plan randomly and lock it in
// so that we don't pick a different one for
// every subsequent terraform plan
resource "random_integer" "plan_idx" {
  min = 0
  max = length(data.equinix_metal_plans.test.plans) - 1
}

// Select a metal facility randomly and lock it in
// so that we don't pick a different one for
// every subsequent terraform plan
resource "random_integer" "facility_idx" {
  min = 0
  max = length(local.facilities) - 1
}

// Select a metal metro randomly and lock it in
// so that we don't pick a different one for
// every subsequent terraform plan
resource "random_integer" "metro_idx" {
  min = 0
  max = length(local.metros) - 1
}

locals {
    // Select a random plan
    selected_plan     = data.equinix_metal_plans.test.plans[random_integer.plan_idx.result]
    plan              = local.selected_plan.slug

    // Select a random facility from the facilities in which the selected plan is available, excluding decommed facilities
    facilities             = sort(tolist(setsubtract(local.selected_plan.available_in, ["nrt1", "dfw2", "ewr1", "ams1", "sjc1", "ld7", "sy4", "ny6"])))
    facility               = local.facilities[random_integer.facility_idx.result]

    // Select a random metro from the metros in which the selected plan is available
    metros             = sort(tolist(local.selected_plan.available_in_metros))
    metro              = local.metros[random_integer.metro_idx.result]

    os = [%s][0]
}
`, fmt.Sprintf("\"%s\"", strings.Join(plans[:], `","`)), fmt.Sprintf("\"%s\"", strings.Join(metros[:], `","`)), fmt.Sprintf("\"%s\"", strings.Join(os[:], `","`)))
}

func testDeviceTerminationTime() string {
	return time.Now().UTC().Add(60 * time.Minute).Format(time.RFC3339)
}

func TestAccMetalDevice_facilityList(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_facility_list(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
				),
			},
		},
	})
}

func TestAccMetalDevice_sshConfig(t *testing.T) {
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"
	userSSHKey, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	projSSHKey, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_ssh_key(rs, userSSHKey, projSSHKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttrPair(
						r,
						"ssh_key_ids.*",
						"equinix_metal_ssh_key.test",
						"id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						r,
						"ssh_key_ids.*",
						"equinix_metal_project_ssh_key.test",
						"id",
					),
				),
			},
		},
	})
}

func TestAccMetalDevice_basic(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_minimal(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
					testAccMetalDeviceNetwork(r),
					resource.TestCheckResourceAttrSet(
						r, "hostname"),
					resource.TestCheckResourceAttr(
						r, "billing_cycle", "hourly"),
					resource.TestCheckResourceAttr(
						r, "network_type", "layer3"),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", ""),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "false"),
					resource.TestCheckResourceAttrSet(
						r, "root_password"),
					resource.TestCheckResourceAttrPair(
						r, "deployed_facility", r, "facilities.0"),
				),
			},
			{
				Config: testAccMetalDeviceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
					testAccMetalDeviceNetwork(r),
					testAccMetalDeviceAttributes(&device),
					testAccMetalDeviceNetworkOrder(r),
				),
			},
		},
	})
}

func TestAccMetalDevice_metro(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_metro(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
					testAccMetalDeviceNetwork(r),
					testAccMetalDeviceAttributes(&device),
				),
			},
		},
	})
}

func TestAccMetalDevice_update(t *testing.T) {
	var d1, d2, d3, d4, d5 packngo.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_varname(rInt, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d1),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt)),
				),
			},
			{
				Config: testAccMetalDeviceConfig_varname(rInt+1, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d2),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt+1)),
					testAccMetalSameDevice(t, &d1, &d2),
				),
			},
			{
				Config: testAccMetalDeviceConfig_varname(rInt+2, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d3),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt+2)),
					resource.TestCheckResourceAttr(r, "description", fmt.Sprintf("test-desc-%d", rInt+2)),
					resource.TestCheckResourceAttr(r, "tags.0", fmt.Sprintf("%d", rInt+2)),
					testAccMetalSameDevice(t, &d2, &d3),
				),
			},
			{
				Config: testAccMetalDeviceConfig_no_description(rInt+3, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d4),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt+3)),
					resource.TestCheckResourceAttr(r, "tags.0", fmt.Sprintf("%d", rInt+3)),
					testAccMetalSameDevice(t, &d3, &d4),
				),
			},
			{
				Config: testAccMetalDeviceConfig_reinstall(rInt+4, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d5),
					testAccMetalSameDevice(t, &d4, &d5),
				),
			},
		},
	})
}

func TestAccMetalDevice_IPXEScriptUrl(t *testing.T) {
	var device, d2 packngo.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test_ipxe_script_url"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_ipxe_script_url(rs, "https://boot.netboot.xyz", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
					testAccMetalDeviceNetwork(r),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", "https://boot.netboot.xyz"),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "true"),
				),
			},
			{
				Config: testAccMetalDeviceConfig_ipxe_script_url(rs, "https://new.netboot.xyz", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d2),
					testAccMetalDeviceNetwork(r),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", "https://new.netboot.xyz"),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "false"),
					testAccMetalSameDevice(t, &device, &d2),
				),
			},
		},
	})
}

func TestAccMetalDevice_IPXEConflictingFields(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test_ipxe_conflict"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMetalDeviceConfig_ipxe_conflict, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
				),
				ExpectError: matchErrShouldNotBeAnIPXE,
			},
		},
	})
}

func TestAccMetalDevice_IPXEConfigMissing(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test_ipxe_config_missing"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMetalDeviceConfig_ipxe_missing, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
				),
				ExpectError: matchErrMustBeProvided,
			},
		},
	})
}

func TestAccMetalDevice_allowUserdataChanges(t *testing.T) {
	var d1, d2 packngo.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "equinix_metal_device.test"

	userdata1 := fmt.Sprintf("#!/usr/bin/env sh\necho 'Allow userdata changes %d'\n", rInt)
	userdata2 := fmt.Sprintf("#!/usr/bin/env sh\necho 'Allow userdata changes %d'\n", rInt+1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_allowAttributeChanges(rInt, rs, userdata1, "", "user_data"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d1),
					resource.TestCheckResourceAttr(r, "user_data", userdata1),
				),
			},
			{
				Config: testAccMetalDeviceConfig_allowAttributeChanges(rInt, rs, userdata2, "", "user_data"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d2),
					resource.TestCheckResourceAttr(r, "user_data", userdata2),
					testAccMetalSameDevice(t, &d1, &d2),
				),
			},
		},
	})
}

func TestAccMetalDevice_allowCustomdataChanges(t *testing.T) {
	var d1, d2 packngo.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "equinix_metal_device.test"

	customdata1 := fmt.Sprintf(`{"message": "Allow customdata changes %d"}`, rInt)
	customdata2 := fmt.Sprintf(`{"message": "Allow customdata changes %d"}`, rInt+1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_allowAttributeChanges(rInt, rs, "", customdata1, "custom_data"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d1),
					resource.TestCheckResourceAttr(r, "custom_data", customdata1),
				),
			},
			{
				Config: testAccMetalDeviceConfig_allowAttributeChanges(rInt, rs, "", customdata2, "custom_data"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d2),
					resource.TestCheckResourceAttr(r, "custom_data", customdata2),
					testAccMetalSameDevice(t, &d1, &d2),
				),
			},
		},
	})
}

func TestAccMetalDevice_allowChangesErrorOnUnsupportedAttribute(t *testing.T) {
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalDeviceConfig_allowAttributeChanges(rInt, rs, "", "", "project_id"),
				ExpectError: regexp.MustCompile(`Error: behavior.allow_changes was given project_id, but only supports \[.+\]`),
			},
		},
	})
}

func testAccMetalDeviceCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_device" {
			continue
		}
		if _, _, err := client.Devices.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Device still exists")
		}
	}
	return nil
}

func testAccMetalDeviceAttributes(device *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if device.Hostname != "tfacc-test-device" {
			return fmt.Errorf("Bad name: %s", device.Hostname)
		}
		if device.State != "active" {
			return fmt.Errorf("Device should be 'active', not '%s'", device.State)
		}

		return nil
	}
}

func testAccMetalDeviceExists(n string, device *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Config).metal

		foundDevice, _, err := client.Devices.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundDevice.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundDevice)
		}

		*device = *foundDevice

		return nil
	}
}

func testAccMetalSameDevice(t *testing.T, before, after *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected device to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func testAccMetalDeviceNetworkOrder(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.Attributes["network.0.family"] != "4" {
			return fmt.Errorf("first netowrk should be public IPv4")
		}
		if rs.Primary.Attributes["network.0.public"] != "true" {
			return fmt.Errorf("first netowrk should be public IPv4")
		}
		if rs.Primary.Attributes["network.1.family"] != "6" {
			return fmt.Errorf("second netowrk should be public IPv6")
		}
		if rs.Primary.Attributes["network.2.family"] != "4" {
			return fmt.Errorf("third netowrk should be private IPv4")
		}
		if rs.Primary.Attributes["network.2.public"] == "true" {
			return fmt.Errorf("third netowrk should be private IPv4")
		}
		return nil
	}
}

func testAccMetalDeviceNetwork(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var ip net.IP
		var k, v string
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		k = "access_public_ipv6"
		v = rs.Primary.Attributes[k]
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("\"%s\" is not a valid IP address: %s",
				k, v)
		}

		k = "access_public_ipv4"
		v = rs.Primary.Attributes[k]
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("\"%s\" is not a valid IP address: %s",
				k, v)
		}

		k = "access_private_ipv4"
		v = rs.Primary.Attributes[k]
		ip = net.ParseIP(v)
		if ip == nil {
			return fmt.Errorf("\"%s\" is not a valid IP address: %s",
				k, v)
		}

		return nil
	}
}

func TestAccMetalDevice_importBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: testExternalProviders,
		Providers:         testAccProviders,
		CheckDestroy:      testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_basic(rs),
			},
			{
				ResourceName:            "equinix_metal_device.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"termination_time"}, // Remove when API returns termination_time for on-demand instances
			},
		},
	})
}

func testAccMetalDeviceConfig_no_description(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-%d"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  tags             = ["%d"]
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, rInt, rInt, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_reinstall(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-%d"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  tags             = ["%d"]
  user_data = "#!/usr/bin/env sh\necho Reinstall\n"
  termination_time = "%s"

  reinstall {
	  enabled = true
	  deprovision_fast = true
  }
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, rInt, rInt, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_allowAttributeChanges(rInt int, projSuffix string, userdata string, customdata string, attributeName string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-%d"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  tags             = ["%d"]
  user_data        = %q
  custom_data      = %q
  termination_time = "%s"

  behavior {
    allow_changes = [
      "%s"
    ]
  }
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, rInt, rInt, userdata, customdata, testDeviceTerminationTime(), attributeName)
}

func testAccMetalDeviceConfig_varname(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-%d"
  description      = "test-desc-%d"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  tags             = ["%d"]
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, rInt, rInt, rInt, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_varname_pxe(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device-%d"
  description      = "test-desc-%d"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  tags             = ["%d"]
  always_pxe       = true
  ipxe_script_url  = "http://matchbox.foo.wtf:8080/boot.ipxe"
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, rInt, rInt, rInt, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_metro(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_minimal(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test" {
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  project_id       = "${equinix_metal_project.test.id}"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix)
}

func testAccMetalDeviceConfig_basic(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}


resource "equinix_metal_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%s"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_ssh_key(projSuffix, userSSHKey, projSSHKey string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-device-%s"
}

resource "equinix_metal_ssh_key" "test" {
	name = "tfacc-ssh-key-%s"
	public_key = "%s"
}

resource "equinix_metal_project_ssh_key" "test" {
	project_id = equinix_metal_project.test.id
	name = "tfacc-project-key-%s"
	public_key = "%s"
}

resource "equinix_metal_device" "test" {
	hostname         = "tfacc-test-device"
	plan             = local.plan
	metro            = local.metro
	operating_system = local.os
	billing_cycle    = "hourly"
	project_id       = equinix_metal_project.test.id
	user_ssh_key_ids = [equinix_metal_ssh_key.test.owner_id]
	project_ssh_key_ids = [equinix_metal_project_ssh_key.test.id]
  }
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, projSuffix, userSSHKey, projSSHKey, projSSHKey)
}

func testAccMetalDeviceConfig_facility_list(projSuffix string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test"  {

  hostname         = "tfacc-device-test-ipxe-script-url"
  plan             = local.plan
  facilities       = local.facilities
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  termination_time = "%s"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, testDeviceTerminationTime())
}

func testAccMetalDeviceConfig_ipxe_script_url(projSuffix, url, pxe string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test_ipxe_script_url"  {

  hostname         = "tfacc-device-test-ipxe-script-url"
  plan             = local.plan
  metro            = local.metro
  operating_system = "custom_ipxe"
  user_data        = "#!/bin/sh\ntouch /tmp/test"
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  ipxe_script_url  = "%s"
  always_pxe       = "%s"
  termination_time = "%s"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), projSuffix, url, pxe, testDeviceTerminationTime())
}

var testAccMetalDeviceConfig_ipxe_conflict = `
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test_ipxe_conflict" {
  hostname         = "tfacc-device-test-ipxe-conflict"
  plan             = local.plan
  metro            = local.metro
  operating_system = "custom_ipxe"
  user_data        = "#!ipxe\nset conflict ipxe_script_url"
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  ipxe_script_url  = "https://boot.netboot.xyz"
  always_pxe       = true
}`

var testAccMetalDeviceConfig_ipxe_missing = `
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-device-%s"
}

resource "equinix_metal_device" "test_ipxe_missing" {
  hostname         = "tfacc-device-test-ipxe-missing"
  plan             = local.plan
  metro            = local.metro
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = "${equinix_metal_project.test.id}"
  always_pxe       = true
}`

type mockDeviceService struct {
	GetFn func(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error)
}

func (m *mockDeviceService) Get(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error) {
	return m.GetFn(deviceID, opts)
}

func (m *mockDeviceService) Create(device *packngo.DeviceCreateRequest) (*packngo.Device, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("Create")
}

func (m *mockDeviceService) Delete(string, bool) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("Delete")
}

func (m *mockDeviceService) List(string, *packngo.ListOptions) ([]packngo.Device, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("List")
}

func (m *mockDeviceService) Update(string, *packngo.DeviceUpdateRequest) (*packngo.Device, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("Update")
}

func (m *mockDeviceService) Reboot(string) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("Reboot")
}

func (m *mockDeviceService) Rescue(string) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("Rescue")
}

func (m *mockDeviceService) Reinstall(string, *packngo.DeviceReinstallFields) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("Reinstall")
}

func (m *mockDeviceService) PowerOff(string) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("PowerOff")
}

func (m *mockDeviceService) PowerOn(string) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("PowerOn")
}

func (m *mockDeviceService) Lock(string) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("Lock")
}

func (m *mockDeviceService) Unlock(string) (*packngo.Response, error) {
	return nil, mockFuncNotImplemented("Unlock")
}

func (m *mockDeviceService) ListBGPSessions(string, *packngo.ListOptions) ([]packngo.BGPSession, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("ListBGPSessions")
}

func (m *mockDeviceService) ListBGPNeighbors(string, *packngo.ListOptions) ([]packngo.BGPNeighbor, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("ListBGPNeighbors")
}

func (m *mockDeviceService) ListEvents(string, *packngo.ListOptions) ([]packngo.Event, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("ListEvents")
}

func (m *mockDeviceService) GetBandwidth(string, *packngo.BandwidthOpts) (*packngo.BandwidthIO, *packngo.Response, error) {
	return nil, nil, mockFuncNotImplemented("GetBandwidth")
}

func mockFuncNotImplemented(f string) error {
	return fmt.Errorf("mockDeviceService %s function not yet implemented", f)
}

var _ packngo.DeviceService = (*mockDeviceService)(nil)

func TestAccMetalDevice_readErrorHandling(t *testing.T) {
	type args struct {
		newResource bool
		meta        *Config
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "forbiddenAfterProvision",
			args: args{
				newResource: false,
				meta: &Config{
					metal: &packngo.Client{
						Devices: &mockDeviceService{
							GetFn: func(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error) {
								httpResp := &http.Response{Status: "403 Forbidden", StatusCode: 403}
								return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "notFoundAfterProvision",
			args: args{
				newResource: false,
				meta: &Config{
					metal: &packngo.Client{
						Devices: &mockDeviceService{
							GetFn: func(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error) {
								httpResp := &http.Response{
									Status:     "404 NotFound",
									StatusCode: 404,
									Header:     http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"12345"}},
								}
								return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "forbiddenWaitForActiveDeviceProvision",
			args: args{
				newResource: true,
				meta: &Config{
					metal: &packngo.Client{
						Devices: &mockDeviceService{
							GetFn: func(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error) {
								httpResp := &http.Response{Status: "403 Forbidden", StatusCode: 403}
								return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "notFoundProvision",
			args: args{
				newResource: true,
				meta: &Config{
					metal: &packngo.Client{
						Devices: &mockDeviceService{
							GetFn: func(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error) {
								httpResp := &http.Response{Status: "404 NotFound", StatusCode: 404}
								return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "errorProvision",
			args: args{
				newResource: true,
				meta: &Config{
					metal: &packngo.Client{
						Devices: &mockDeviceService{
							GetFn: func(deviceID string, opts *packngo.GetOptions) (*packngo.Device, *packngo.Response, error) {
								httpResp := &http.Response{Status: "400 BadRequest", StatusCode: 400}
								return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d *schema.ResourceData
			d = new(schema.ResourceData)
			if tt.args.newResource {
				d.MarkNewResource()
			}
			if err := resourceMetalDeviceRead(d, tt.args.meta); (err != nil) != tt.wantErr {
				t.Errorf("resourceMetalDeviceRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
