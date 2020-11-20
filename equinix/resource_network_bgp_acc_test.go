package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	spEnvVar      = "TF_ACC_ECX_L2_SP_NAME"
	authkeyEnvVar = "TF_ACC_ECX_L2_AUTHKEY"
)

func TestAccNeBGP(t *testing.T) {
	t.Parallel()
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	spName, _ := schema.EnvDefaultFunc(spEnvVar, "AWS Direct Connect")()
	authKey, _ := schema.EnvDefaultFunc(authkeyEnvVar, "123456789012")()
	contextBasic := map[string]interface{}{
		"device_resourceName":    "test",
		"device_name":            fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device_throughput":      500,
		"device_throughput_unit": "Mbps",
		"device_metro_code":      metro.(string),
		"device_type_code":       "CSR1000V",
		"device_package_code":    "SEC",
		"device_notifications":   []string{"marry@equinix.com", "john@equinix.com"},
		"device_hostname":        fmt.Sprintf("tf-%s", randString(6)),
		"device_term_length":     1,
		"device_version":         "16.09.05",
		"device_core_count":      2,
		"conn_resourceName":      "test",
		"conn_name":              fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"conn_profile_name":      spName.(string),
		"conn_speed":             50,
		"conn_speed_unit":        "MB",
		"conn_notifications":     []string{"marry@equinix.com", "john@equinix.com"},
		"conn_authorization_key": authKey.(string),
		"bgp_resourceName":       "test",
		"bgp_local_ip_address":   "1.1.1.1/30",
		"bgp_local_asn":          12345,
		"bgp_remote_ip_address":  "1.1.1.2",
		"bgp_remote_asn":         22211,
	}
	contextUpdate := copyMap(contextBasic)
	contextUpdate["bgp_authentication_key"] = randString(10)
	resourceName := fmt.Sprintf("equinix_network_bgp.%s", contextBasic["bgp_resourceName"].(string))
	var bgpConfig ne.BGPConfiguration
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNeBGP(contextBasic),
				Check: resource.ComposeTestCheckFunc(
					testAccNeBGPExists(resourceName, &bgpConfig),
					testAccNeBGPAttributes(&bgpConfig, contextBasic),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "provisioning_status", ne.BGPProvisioningStatusProvisioned),
				),
			},
			{
				Config: testAccNeBGPUpdate(contextUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccNeBGPExists(resourceName, &bgpConfig),
					testAccNeBGPAttributes(&bgpConfig, contextUpdate),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "provisioning_status", ne.BGPProvisioningStatusPendingUpdate),
				),
			},
		},
	})
}

func testAccNeBGP(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{device_metro_code}"
  status     = "Active"
}

resource "equinix_network_device" "%{device_resourceName}" {
  name            = "%{device_name}"
  throughput      = %{device_throughput}
  throughput_unit = "%{device_throughput_unit}"
  metro_code      = data.equinix_network_account.test.metro_code
  type_code       = "%{device_type_code}"
  package_code    = "%{device_package_code}"
  notifications   = %{device_notifications}
  hostname        = "%{device_hostname}"
  term_length     = %{device_term_length}
  account_number  = data.equinix_network_account.test.number
  version         = "%{device_version}"
  core_count      = %{device_core_count}
}

data "equinix_ecx_l2_sellerprofile" "test" {
  name = "%{conn_profile_name}"
}

locals {
  sp_first_metro        = tolist(data.equinix_ecx_l2_sellerprofile.test.metro)[0]
  sp_first_metro_region = keys(local.sp_first_metro.regions)[0]
}

resource "equinix_ecx_l2_connection" "%{conn_resourceName}" {
  name              = "%{conn_name}"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.test.uuid
  speed             = %{conn_speed}
  speed_unit        = "%{conn_speed_unit}"
  notifications     = %{conn_notifications}
  device_uuid       = equinix_network_device.%{device_resourceName}.uuid
  seller_region     = local.sp_first_metro_region
  seller_metro_code = local.sp_first_metro.code
  authorization_key = "%{conn_authorization_key}"
}

resource "equinix_network_bgp" "%{bgp_resourceName}" {
  connection_id      = equinix_ecx_l2_connection.%{conn_resourceName}.id
  local_ip_address   = "%{bgp_local_ip_address}"
  local_asn          = %{bgp_local_asn}
  remote_ip_address  = "%{bgp_remote_ip_address}"
  remote_asn         = %{bgp_remote_asn}
}
`, ctx)
}

func testAccNeBGPUpdate(ctx map[string]interface{}) string {
	return nprintf(`
data "equinix_network_account" "test" {
  metro_code = "%{device_metro_code}"
  status     = "Active"
}

resource "equinix_network_device" "%{device_resourceName}" {
  name            = "%{device_name}"
  throughput      = %{device_throughput}
  throughput_unit = "%{device_throughput_unit}"
  metro_code      = data.equinix_network_account.test.metro_code
  type_code       = "%{device_type_code}"
  package_code    = "%{device_package_code}"
  notifications   = %{device_notifications}
  hostname        = "%{device_hostname}"
  term_length     = %{device_term_length}
  account_number  = data.equinix_network_account.test.number
  version         = "%{device_version}"
  core_count      = %{device_core_count}
}

data "equinix_ecx_l2_sellerprofile" "test" {
  name = "%{conn_profile_name}"
}

locals {
  sp_first_metro        = tolist(data.equinix_ecx_l2_sellerprofile.test.metro)[0]
  sp_first_metro_region = keys(local.sp_first_metro.regions)[0]
}

resource "equinix_ecx_l2_connection" "%{conn_resourceName}" {
  name              = "%{conn_name}"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.test.uuid
  speed             = %{conn_speed}
  speed_unit        = "%{conn_speed_unit}"
  notifications     = %{conn_notifications}
  device_uuid       = equinix_network_device.%{device_resourceName}.uuid
  seller_region     = local.sp_first_metro_region
  seller_metro_code = local.sp_first_metro.code
  authorization_key = "%{conn_authorization_key}"
}

resource "equinix_network_bgp" "%{bgp_resourceName}" {
  connection_id      = equinix_ecx_l2_connection.%{conn_resourceName}.id
  local_ip_address   = "%{bgp_local_ip_address}"
  local_asn          = %{bgp_local_asn}
  remote_ip_address  = "%{bgp_remote_ip_address}"
  remote_asn         = %{bgp_remote_asn}
  authentication_key = "%{bgp_authentication_key}"
}
`, ctx)
}

func testAccNeBGPExists(resourceName string, bgpConfig *ne.BGPConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*Config).ne
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource has no ID attribute set")
		}
		resp, err := client.GetBGPConfiguration(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error when fetching BGP configuration '%s': %s", rs.Primary.ID, err)
		}
		*bgpConfig = *resp
		return nil
	}
}

func testAccNeBGPAttributes(bgpConfig *ne.BGPConfiguration, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["bgp_local_ip_address"]; ok && bgpConfig.LocalIPAddress != v.(string) {
			return fmt.Errorf("local_ip_address does not match %v - %v", bgpConfig.LocalIPAddress, v)
		}
		if v, ok := ctx["bgp_local_asn"]; ok && bgpConfig.LocalASN != v.(int) {
			return fmt.Errorf("local_asn does not match %v - %v", bgpConfig.LocalASN, v)
		}
		if v, ok := ctx["bgp_remote_ip_address"]; ok && bgpConfig.RemoteIPAddress != v.(string) {
			return fmt.Errorf("remote_ip_address does not match %v - %v", bgpConfig.RemoteIPAddress, v)
		}
		if v, ok := ctx["bgp_remote_asn"]; ok && bgpConfig.RemoteASN != v.(int) {
			return fmt.Errorf("remote_asn does not match %v - %v", bgpConfig.RemoteASN, v)
		}
		if v, ok := ctx["bgp_authentication_key"]; ok && bgpConfig.AuthenticationKey != v.(string) {
			return fmt.Errorf("authentication_key does not match %v - %v", bgpConfig.AuthenticationKey, v)
		}
		if v, ok := ctx["conn_resourceName"]; ok {
			connResourceName := "equinix_ecx_l2_connection." + v.(string)
			rs, ok := s.RootModule().Resources[connResourceName]
			if !ok {
				return fmt.Errorf("related connection resource not found: %s", connResourceName)
			}
			if bgpConfig.ConnectionUUID != rs.Primary.ID {
				return fmt.Errorf("connection_id does not match %v - %v", bgpConfig.ConnectionUUID, rs.Primary.ID)
			}
		}
		if v, ok := ctx["device_resourceName"]; ok {
			deviceResourceName := "equinix_network_device." + v.(string)
			rs, ok := s.RootModule().Resources[deviceResourceName]
			if !ok {
				return fmt.Errorf("related device resource not found: %s", deviceResourceName)
			}
			if bgpConfig.DeviceUUID != rs.Primary.ID {
				return fmt.Errorf("device_id does not match %v - %v", bgpConfig.DeviceUUID, rs.Primary.ID)
			}
		}
		return nil
	}
}
