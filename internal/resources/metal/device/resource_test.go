package device_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/resources/metal/device"
	"github.com/google/uuid"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Regexp vars for use with resource.ExpectError
var (
	matchErrMustBeProvided     = regexp.MustCompile(".* must be provided when .*")
	matchErrShouldNotBeAnIPXE  = regexp.MustCompile(`.*"user_data" should not be an iPXE.*`)
	matchErrDeviceReadyTimeout = regexp.MustCompile(".* timeout while waiting for state to become 'active, failed'.*")
	matchErrDeviceLocked       = regexp.MustCompile(".*Cannot delete a locked item.*")
)

func TestMetalDevice_readErrorHandling(t *testing.T) {
	type args struct {
		newResource bool
		handler     func(w http.ResponseWriter, r *http.Request)
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
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusForbidden)
				},
			},
			wantErr: false,
		},
		{
			name: "notFoundAfterProvision",
			args: args{
				newResource: false,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusNotFound)
				},
			},
			wantErr: false,
		},
		{
			name: "forbiddenWaitForActiveDeviceProvision",
			args: args{
				newResource: true,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusForbidden)
				},
			},
			wantErr: true,
		},
		{
			name: "notFoundProvision",
			args: args{
				newResource: true,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusNotFound)
				},
			},
			wantErr: true,
		},
		{
			name: "errorProvision",
			args: args{
				newResource: true,
				handler: func(w http.ResponseWriter, r *http.Request) {
					w.Header().Add("Content-Type", "application/json")
					w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
					w.WriteHeader(http.StatusBadRequest)
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			d := new(schema.ResourceData)
			if tt.args.newResource {
				d.MarkNewResource()
			} else {
				d.SetId(uuid.New().String())
			}

			mockAPI := httptest.NewServer(http.HandlerFunc(tt.args.handler))
			meta := &config.Config{
				BaseURL: mockAPI.URL,
				Token:   "fakeTokenForMock",
			}

			err := meta.Load(ctx)
			if err != nil {
				log.Printf("failed to load provider config during test: %v", err)
			}

			if err := device.Read(ctx, d, meta); (err != nil) != tt.wantErr {
				t.Errorf("device.Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			mockAPI.Close()
		})
	}
}

func TestAccMetalDevice_facilityList(t *testing.T) {
	var device metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
	var device metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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

func TestAccMetalDevice_update(t *testing.T) {
	var d1, d2, d3, d4, d5 metalv1.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
	var device, d2 metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test_ipxe_script_url"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
	var device metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test_ipxe_conflict"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMetalDeviceConfig_ipxe_conflict, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
				),
				ExpectError: matchErrShouldNotBeAnIPXE,
			},
		},
	})
}

func TestAccMetalDevice_IPXEConfigMissing(t *testing.T) {
	var device metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test_ipxe_config_missing"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMetalDeviceConfig_ipxe_missing, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), rs),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &device),
				),
				ExpectError: matchErrMustBeProvided,
			},
		},
	})
}

func TestAccMetalDevice_allowUserdataChanges(t *testing.T) {
	var d1, d2 metalv1.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "equinix_metal_device.test"

	userdata1 := fmt.Sprintf("#!/usr/bin/env sh\necho 'Allow userdata changes %d'\n", rInt)
	userdata2 := fmt.Sprintf("#!/usr/bin/env sh\necho 'Allow userdata changes %d'\n", rInt+1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
	var d1, d2 metalv1.Device
	rs := acctest.RandString(10)
	rInt := acctest.RandInt()
	r := "equinix_metal_device.test"

	customdata1 := fmt.Sprintf(`{"message": "Allow customdata changes %d"}`, rInt)
	customdata2 := fmt.Sprintf(`{"message": "Allow customdata changes %d"}`, rInt+1)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalDeviceConfig_allowAttributeChanges(rInt, rs, "", "", "project_id"),
				ExpectError: regexp.MustCompile(`Error: behavior.allow_changes was given project_id, but only supports \[.+\]`),
			},
		},
	})
}

func testAccMetalDeviceCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

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

func testAccMetalDeviceAttributes(device *metalv1.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if device.GetHostname() != "tfacc-test-device" {
			return fmt.Errorf("Bad name: %s", device.GetHostname())
		}
		if device.GetState() != metalv1.DEVICESTATE_ACTIVE {
			return fmt.Errorf("Device should be 'active', not '%s'", device.GetState())
		}

		return nil
	}
}

func testAccMetalDeviceExists(n string, device *metalv1.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

		foundDevice, _, err := client.DevicesApi.FindDeviceById(context.TODO(), rs.Primary.ID).Execute()
		if err != nil {
			return err
		}
		if foundDevice.GetId() != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundDevice)
		}

		*device = *foundDevice

		return nil
	}
}

func testAccMetalSameDevice(t *testing.T, before, after *metalv1.Device) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.GetId() != after.GetId() {
			t.Fatalf("Expected device to be the same, but it was recreated: %s -> %s", before.GetId(), after.GetId())
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
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
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
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, rInt, rInt, acceptance.TestDeviceTerminationTime())
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
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, rInt, rInt, acceptance.TestDeviceTerminationTime())
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
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, rInt, rInt, userdata, customdata, acceptance.TestDeviceTerminationTime(), attributeName)
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
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, rInt, rInt, rInt, acceptance.TestDeviceTerminationTime())
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
}`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix)
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
}`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, acceptance.TestDeviceTerminationTime())
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
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, projSuffix, userSSHKey, projSSHKey, projSSHKey)
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
}`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, acceptance.TestDeviceTerminationTime())
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
}`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, url, pxe, acceptance.TestDeviceTerminationTime())
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

func testAccMetalDeviceConfig_timeout(projSuffix, createTimeout, updateTimeout, deleteTimeout string) string {
	if createTimeout == "" {
		createTimeout = "20m"
	}
	if updateTimeout == "" {
		updateTimeout = "20m"
	}
	if deleteTimeout == "" {
		deleteTimeout = "20m"
	}

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

  timeouts {
	create = "%s"
	update = "%s"
    delete = "%s"
  }
}
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, acceptance.TestDeviceTerminationTime(), createTimeout, updateTimeout, deleteTimeout)
}

func testAccMetalDeviceConfig_reinstall_timeout(projSuffix, updateTimeout string) string {
	if updateTimeout == "" {
		updateTimeout = "20m"
	}

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
  user_data = "#!/usr/bin/env sh\necho Reinstall\n"
  termination_time = "%s"

  reinstall {
	  enabled = true
	  deprovision_fast = true
  }

  timeouts {
	update = "%s"
  }
}
`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, acceptance.TestDeviceTerminationTime(), updateTimeout)
}

func testAccWaitForMetalDeviceActive(project, deviceHostName string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		defaultTimeout := 20 * time.Minute

		rs, ok := state.RootModule().Resources[project]
		if !ok {
			return "", fmt.Errorf("Project Not found in the state: %s", project)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("No Record ID is set")
		}

		meta := acceptance.TestAccProvider.Meta()
		rd := new(schema.ResourceData)
		client := meta.(*config.Config).NewMetalClientForTesting()
		resp, _, err := client.DevicesApi.FindProjectDevices(context.TODO(), rs.Primary.ID).Search(deviceHostName).Execute()
		if err != nil {
			return "", fmt.Errorf("error while fetching devices for project [%s], error: %w", rs.Primary.ID, err)
		}
		devices := resp.Devices
		if len(devices) == 0 {
			return "", fmt.Errorf("Not able to find devices in project [%s]", rs.Primary.ID)
		}
		if len(devices) > 1 {
			return "", fmt.Errorf("Found more than one device with the hostname in project [%s]", rs.Primary.ID)
		}

		rd.SetId(devices[0].GetId())
		return devices[0].GetId(), device.WaitForActiveDevice(context.Background(), rd, acceptance.TestAccProvider.Meta(), defaultTimeout)
	}
}

func TestAccMetalDeviceCreate_timeout(t *testing.T) {
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"
	hostname := "tfacc-test-device"
	project := "equinix_metal_project.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalDeviceConfig_timeout(rs, "10s", "", ""),
				ExpectError: matchErrDeviceReadyTimeout,
			},
			{
				/**
				Step 1 errors out, state doesnt have device, need to import that in the state before deleting
				*/
				ResourceName:       r,
				ImportState:        true,
				ImportStateIdFunc:  testAccWaitForMetalDeviceActive(project, hostname),
				ImportStatePersist: true,
			},
			{
				Config:  testAccMetalDeviceConfig_timeout(rs, "", "", ""),
				Destroy: true,
			},
		},
	})
}

func TestAccMetalDeviceUpdate_timeout(t *testing.T) {
	var d1 metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_timeout(rs, "", "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d1),
				),
			},
			{
				Config:      testAccMetalDeviceConfig_reinstall_timeout(rs, "10s"),
				ExpectError: matchErrDeviceReadyTimeout,
			},
		},
	})
}

func TestAccMetalDevice_LockingAndUnlocking(t *testing.T) {
	var d1 metalv1.Device
	rs := acctest.RandString(10)
	r := "equinix_metal_device.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalDeviceCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalDeviceConfig_lockable(rs, true),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d1),
				),
			},
			{
				Config:      testAccMetalDeviceConfig_lockable(rs, true),
				Destroy:     true,
				ExpectError: matchErrDeviceLocked,
			},
			{
				Config: testAccMetalDeviceConfig_lockable(rs, false),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalDeviceExists(r, &d1),
				),
			},
			{
				Config:  testAccMetalDeviceConfig_lockable(rs, false),
				Destroy: true,
			},
		},
	})
}

func testAccMetalDeviceConfig_lockable(projSuffix string, locked bool) string {
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
  locked           = %v
  termination_time = "%s"
}`, acceptance.ConfAccMetalDevice_base(acceptance.Preferable_plans, acceptance.Preferable_metros, acceptance.Preferable_os), projSuffix, locked, acceptance.TestDeviceTerminationTime())
}
