package packet

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("packet_device", &resource.Sweeper{
		Name:         "packet_device",
		F:            testSweepDevices,
		Dependencies: []string{"packet_volume"},
	})
}

func testSweepDevices(region string) error {
	log.Printf("[DEBUG] Sweeping devices")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping devices: %s", err)
	}
	client := meta.(*packngo.Client)

	ps, _, err := client.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("Error getting project list for sweepeing devices: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if strings.HasPrefix(p.Name, "tfacc-") {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := client.Devices.List(pid, nil)
		if err != nil {
			return fmt.Errorf("Error listing devices to sweep: %s", err)
		}
		for _, d := range ds {
			dids = append(dids, d.ID)
		}
	}

	for _, did := range dids {
		log.Printf("Removing device %s", did)
		_, err := client.Devices.Delete(did, true)
		if err != nil {
			return fmt.Errorf("Error deleting device %s", err)
		}
	}
	return nil
}

// Regexp vars for use with resource.ExpectError
var matchErrMustBeProvided = regexp.MustCompile(".* must be provided when .*")
var matchErrShouldNotBeAnIPXE = regexp.MustCompile(`.*"user_data" should not be an iPXE.*`)

func TestAccPacketDevice_FacilityList(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "packet_device.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDeviceConfig_facility_list(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &device),
				),
			},
		},
	})
}

func TestAccPacketDevice_NetworkPortsOrder(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "packet_device.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDeviceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &device),
					testAccCheckPacketDeviceNetworkOrder(r),
					testAccCheckPacketDevicePortsOrder(r),
				),
			},
		},
	})
}

func TestAccPacketDevice_Basic(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "packet_device.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDeviceConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &device),
					testAccCheckPacketDeviceNetwork(r),
					testAccCheckPacketDeviceAttributes(&device),
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
		},
	})
}

func TestAccPacketDevice_Update(t *testing.T) {
	var d1, d2, d3, d4 packngo.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "packet_device.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDeviceConfig_varname(rInt, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &d1),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt)),
				),
			},
			{
				Config: testAccCheckPacketDeviceConfig_varname(rInt+1, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &d2),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt+1)),
					testAccCheckPacketSameDevice(t, &d1, &d2),
				),
			},
			{
				Config: testAccCheckPacketDeviceConfig_varname(rInt+2, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &d3),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt+2)),
					resource.TestCheckResourceAttr(r, "description", fmt.Sprintf("test-desc-%d", rInt+2)),
					resource.TestCheckResourceAttr(r, "tags.0", fmt.Sprintf("%d", rInt+2)),
					testAccCheckPacketSameDevice(t, &d2, &d3),
				),
			},
			{
				Config: testAccCheckPacketDeviceConfig_no_description(rInt+3, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &d4),
					resource.TestCheckResourceAttr(r, "hostname", fmt.Sprintf("tfacc-test-device-%d", rInt+3)),
					resource.TestCheckResourceAttr(r, "tags.0", fmt.Sprintf("%d", rInt+3)),
					testAccCheckPacketSameDevice(t, &d3, &d4),
				),
			},
		},
	})
}

func TestAccPacketDevice_IPXEScriptUrl(t *testing.T) {
	var device, d2 packngo.Device
	rs := acctest.RandString(10)
	r := "packet_device.test_ipxe_script_url"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDeviceConfig_ipxe_script_url(rs, "https://boot.netboot.xyz", "true"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &device),
					testAccCheckPacketDeviceNetwork(r),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", "https://boot.netboot.xyz"),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "true"),
				),
			},
			{
				Config: testAccCheckPacketDeviceConfig_ipxe_script_url(rs, "https://new.netboot.xyz", "false"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &d2),
					testAccCheckPacketDeviceNetwork(r),
					resource.TestCheckResourceAttr(
						r, "ipxe_script_url", "https://new.netboot.xyz"),
					resource.TestCheckResourceAttr(
						r, "always_pxe", "false"),
					testAccCheckPacketSameDevice(t, &device, &d2),
				),
			},
		},
	})
}

func TestAccPacketDevice_IPXEConflictingFields(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "packet_device.test_ipxe_conflict"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckPacketDeviceConfig_ipxe_conflict, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &device),
				),
				ExpectError: matchErrShouldNotBeAnIPXE,
			},
		},
	})
}

func TestAccPacketDevice_IPXEConfigMissing(t *testing.T) {
	var device packngo.Device
	rs := acctest.RandString(10)
	r := "packet_device.test_ipxe_config_missing"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckPacketDeviceConfig_ipxe_missing, rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketDeviceExists(r, &device),
				),
				ExpectError: matchErrMustBeProvided,
			},
		},
	})
}

func testAccCheckPacketDeviceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_device" {
			continue
		}
		if _, _, err := client.Devices.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Device still exists")
		}
	}
	return nil
}

func testAccCheckPacketDeviceAttributes(device *packngo.Device) resource.TestCheckFunc {
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

func testAccCheckPacketDeviceExists(n string, device *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

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

func testAccCheckPacketSameDevice(t *testing.T, before, after *packngo.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected device to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func testAccCheckPacketDevicePortsOrder(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.Attributes["ports.0.name"] != "bond0" {
			return fmt.Errorf("first port should be bond0")
		}
		if rs.Primary.Attributes["ports.1.name"] != "eth0" {
			return fmt.Errorf("second port should be eth0")
		}
		if rs.Primary.Attributes["ports.2.name"] != "eth1" {
			return fmt.Errorf("third port should be eth1")
		}
		return nil
	}
}

func testAccCheckPacketDeviceNetworkOrder(n string) resource.TestCheckFunc {
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

func testAccCheckPacketDeviceNetwork(n string) resource.TestCheckFunc {
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

func TestAccPacketDevice_importBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketDeviceConfig_basic(rs),
			},
			{
				ResourceName:      "packet_device.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketDeviceConfig_no_description(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-device-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-test-device-%d"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  tags             = ["%d"]
}
`, projSuffix, rInt, rInt)
}

func testAccCheckPacketDeviceConfig_varname(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-device-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-test-device-%d"
  description      = "test-desc-%d"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  tags             = ["%d"]
}
`, projSuffix, rInt, rInt, rInt)
}

func testAccCheckPacketDeviceConfig_varname_pxe(rInt int, projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-device-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-test-device-%d"
  description      = "test-desc-%d"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  tags             = ["%d"]
  always_pxe       = true
  ipxe_script_url  = "http://matchbox.foo.wtf:8080/boot.ipxe"
}
`, projSuffix, rInt, rInt, rInt)
}

func testAccCheckPacketDeviceConfig_basic(projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-device-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-test-device"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}`, projSuffix)
}

func testAccCheckPacketDeviceConfig_facility_list(projSuffix string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
  name = "tfacc-device-%s"
}

resource "packet_device" "test"  {

  hostname         = "tfacc-device-test-ipxe-script-url"
  plan             = "t1.small.x86"
  facilities       = ["sjc1", "any"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}`, projSuffix)
}

func testAccCheckPacketDeviceConfig_ipxe_script_url(projSuffix, url, pxe string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
  name = "tfacc-device-%s"
}

resource "packet_device" "test_ipxe_script_url"  {

  hostname         = "tfacc-device-test-ipxe-script-url"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "custom_ipxe"
  user_data        = "#!/bin/sh\ntouch /tmp/test"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  ipxe_script_url  = "%s"
  always_pxe       = "%s"
}`, projSuffix, url, pxe)
}

var testAccCheckPacketDeviceConfig_ipxe_conflict = `
resource "packet_project" "test" {
  name = "tfacc-device-%s"
}

resource "packet_device" "test_ipxe_conflict" {
  hostname         = "tfacc-device-test-ipxe-conflict"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "custom_ipxe"
  user_data        = "#!ipxe\nset conflict ipxe_script_url"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  ipxe_script_url  = "https://boot.netboot.xyz"
  always_pxe       = true
}`

var testAccCheckPacketDeviceConfig_ipxe_missing = `
resource "packet_project" "test" {
  name = "tfacc-device-%s"
}

resource "packet_device" "test_ipxe_missing" {
  hostname         = "tfacc-device-test-ipxe-missing"
  plan             = "t1.small.x86"
  facilities       = ["sjc1"]
  operating_system = "custom_ipxe"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  always_pxe       = true
}`
