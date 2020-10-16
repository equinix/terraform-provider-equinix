package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNeBGP(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{
		"device_resourceName":    "test",
		"device_name":            fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device_throughput":      500,
		"device_throughput_unit": "Mbps",
		"device_metro_code":      "SV",
		"device_type_code":       "CSR1000V",
		"device_package_code":    "SEC",
		"device_notifications":   []string{"marry@equinix.com", "john@equinix.com"},
		"device_hostname":        fmt.Sprintf("tf-%s", randString(6)),
		"device_term_length":     1,
		"device_version":         "16.09.05",
		"device_core_count":      2,
		"conn_resourceName":      "test",
		"conn_name":              fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"conn_profile_name":      "AWS Direct Connect",
		"conn_speed":             50,
		"conn_speed_unit":        "MB",
		"conn_notifications":     []string{"marry@equinix.com", "john@equinix.com"},
		"conn_seller_region":     "us-west-2",
		"conn_seller_metro_code": "SV",
		"conn_authorization_key": "123456789012",
		"bgp_resourceName":       "test",
		"bgp_local_ip_address":   "1.1.1.1/30",
		"bgp_local_asn":          12345,
		"bgp_remote_ip_address":  "1.1.1.2",
		"bgp_remote_asn":         22211,
		"bgp_authentication_key": "secret",
	}
	resourceName := fmt.Sprintf("equinix_network_bgp.%s", context["bgp_resourceName"].(string))
	var bgpConfig ne.BGPConfiguration
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNeBGP(context),
				Check: resource.ComposeTestCheckFunc(
					testAccNeBGPExists(resourceName, &bgpConfig),
					testAccNeBGPAttributes(bgpConfig, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "device_id"),
					resource.TestCheckResourceAttrSet(resourceName, "provisioning_status"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
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
  metro_code      = data.equinix_ne_account.test.metro_code
  type_code       = "%{device_type_code}"
  package_code    = "%{device_package_code}"
  notifications   = %{device_notifications}
  hostname        = "%{device_hostname}"
  term_length     = %{device_term_length}
  account_number  = data.equinix_ne_account.test.number
  version         = "%{device_version}"
  core_count      = %{device_core_count}
}

data "equinix_ecx_l2_sellerprofile" "test" {
  name = "%{conn_profile_name}"
}

resource "equinix_ecx_l2_connection" "%{conn_resourceName}" {
  name              = "%{conn_name}"
  profile_uuid      = data.equinix_ecx_l2_sellerprofile.test.uuid
  speed             = %{conn_speed}
  speed_unit        = "%{conn_speed_unit}"
  notifications     = %{conn_notifications}
  device_uuid       = equinix_network_device.%{device_resourceName}.uuid
  seller_region     = "%{conn_seller_region}"
  seller_metro_code = "%{conn_seller_metro_code}"
  authorization_key = "%{conn_authorization_key}"
}

resource "equinix_network_bgp" "%{bgp_resourceName}" {
  connection_id      = equinix_ecx_l2_connection."%{conn_resourceName}".uuid
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

func testAccNeBGPAttributes(bgpConfig ne.BGPConfiguration, ctx map[string]interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if v, ok := ctx["bgp_local_ip_address"]; ok && bgpConfig.LocalIPAddress != v.(string) {
			return fmt.Errorf("LocalIPAddress does not match %v - %v", bgpConfig.LocalIPAddress, v)
		}
		if v, ok := ctx["bgp_local_asn"]; ok && bgpConfig.LocalASN != v.(int) {
			return fmt.Errorf("LocalASN does not match %v - %v", bgpConfig.LocalASN, v)
		}
		if v, ok := ctx["bgp_remote_ip_address"]; ok && bgpConfig.RemoteIPAddress != v.(string) {
			return fmt.Errorf("RemoteIPAddress does not match %v - %v", bgpConfig.RemoteIPAddress, v)
		}
		if v, ok := ctx["bgp_remote_asn"]; ok && bgpConfig.RemoteASN != v.(int) {
			return fmt.Errorf("RemoteASN does not match %v - %v", bgpConfig.RemoteASN, v)
		}
		if v, ok := ctx["bgp_authentication_key"]; ok && bgpConfig.AuthenticationKey != v.(string) {
			return fmt.Errorf("AuthenticationKey does not match %v - %v", bgpConfig.AuthenticationKey, v)
		}
		if v, ok := ctx["conn_resourceName"]; ok {
			connResourceName := "equinix_ecx_l2_connection." + v.(string)
			rs, ok := s.RootModule().Resources[connResourceName]
			if !ok {
				return fmt.Errorf("related connection resource not found: %s", connResourceName)
			}
			if bgpConfig.ConnectionUUID != rs.Primary.ID {
				return fmt.Errorf("ConnectionUUID does not match %v - %v", bgpConfig.ConnectionUUID, rs.Primary.ID)
			}
		}
		return nil
	}
}
