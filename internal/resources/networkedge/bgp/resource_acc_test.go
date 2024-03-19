package bgp

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/nprintf"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkBGP_CSR1000V_Single_AWS(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	spName, _ := schema.EnvDefaultFunc(awsSpEnvVar, "AWS Direct Connect")()
	authKey, _ := schema.EnvDefaultFunc(awsAuthKeyEnvVar, "123456789012")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	context := map[string]interface{}{
		"device-resourceName":          "test",
		"device-account_name":          accountName.(string),
		"device-self_managed":          false,
		"device-byol":                  false,
		"device-name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-throughput":            500,
		"device-throughput_unit":       "Mbps",
		"device-metro_code":            metro.(string),
		"device-type_code":             "CSR1000V",
		"device-package_code":          "SEC",
		"device-notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":              fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":           1,
		"device-version":               "16.09.05",
		"device-core_count":            2,
		"user-resourceName":            "test",
		"user-username":                fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"user-password":                acctest.RandString(10),
		"connection-resourceName":      "test",
		"connection-name":              fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"connection-profile_name":      spName.(string),
		"connection-speed":             50,
		"connection-speed_unit":        "MB",
		"connection-notifications":     []string{"marry@equinix.com", "john@equinix.com"},
		"connection-seller_metro_code": "SV",
		"connection-seller_region":     "us-west-1",
		"connection-authorization_key": authKey.(string),
		"bgp-resourceName":             "test",
		"bgp-local_ip_address":         "1.1.1.1/30",
		"bgp-local_asn":                12345,
		"bgp-remote_ip_address":        "1.1.1.2",
		"bgp-remote_asn":               22211,
	}
	contextWithChanges := copyMap(context)
	contextWithChanges["bgp-authentication_key"] = acctest.RandString(10)
	resourceName := fmt.Sprintf("equinix_network_bgp.%s", context["bgp-resourceName"].(string))
	var bgpConfig ne.BGPConfiguration
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withSSHUser().withConnection().withBGP().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeBGPExists(resourceName, &bgpConfig),
					testAccNeBGPAttributes(&bgpConfig, context),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "provisioning_status", ne.BGPProvisioningStatusProvisioned),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: newTestAccConfig(contextWithChanges).withDevice().withSSHUser().withConnection().withBGP().build(),
				Check: resource.ComposeTestCheckFunc(
					testAccNeBGPExists(resourceName, &bgpConfig),
					testAccNeBGPAttributes(&bgpConfig, contextWithChanges),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttr(resourceName, "provisioning_status", ne.BGPProvisioningStatusPendingUpdate),
				),
			},
		},
	})
}

func (t *testAccConfig) withBGP() *testAccConfig {
	t.config += testAccNetworkBGP(t.ctx)
	return t
}

func testAccNetworkBGP(ctx map[string]interface{}) string {
	var config string
	config += nprintf.NPrintf(`
resource "equinix_network_bgp" "%{bgp-resourceName}" {
  connection_id      = equinix_ecx_l2_connection.%{connection-resourceName}.id
  local_ip_address   = "%{bgp-local_ip_address}"
  local_asn          = %{bgp-local_asn}
  remote_ip_address  = "%{bgp-remote_ip_address}"
  remote_asn         = %{bgp-remote_asn}`, ctx)
	if _, ok := ctx["bgp-authentication_key"]; ok {
		config += nprintf.NPrintf(`
  authentication_key = "%{bgp-authentication_key}"`, ctx)
	}
	config += `
}`
	if _, ok := ctx["connection-secondary_name"]; ok {
		config += nprintf.NPrintf(`
resource "equinix_network_bgp" "%{bgp-secondary_resourceName}" {
  connection_id      = equinix_ecx_l2_connection.%{connection-resourceName}.id
  local_ip_address   = "%{bgp-secondary_local_ip_address}"
  local_asn          = %{bgp-secondary_local_asn}
  remote_ip_address  = "%{bgp-secondary_remote_ip_address}"
  remote_asn         = %{bgp-secondary_remote_asn}`, ctx)
		if _, ok := ctx["bgp-secondary_authentication_key"]; ok {
			config += nprintf.NPrintf(`
  authentication_key = "%{bgp-secondary_authentication_key}"`, ctx)
		}
		config += `
}`
	}
	return config
}

func testAccNeBGPExists(resourceName string, bgpConfig *ne.BGPConfiguration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		client := testAccProvider.Meta().(*config.Config).Ne
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
		if v, ok := ctx["bgp-local_ip_address"]; ok && ne.StringValue(bgpConfig.LocalIPAddress) != v.(string) {
			return fmt.Errorf("local_ip_address does not match %v - %v", ne.StringValue(bgpConfig.LocalIPAddress), v)
		}
		if v, ok := ctx["bgp-local_asn"]; ok && ne.IntValue(bgpConfig.LocalASN) != v.(int) {
			return fmt.Errorf("local_asn does not match %v - %v", ne.IntValue(bgpConfig.LocalASN), v)
		}
		if v, ok := ctx["bgp-remote_ip_address"]; ok && ne.StringValue(bgpConfig.RemoteIPAddress) != v.(string) {
			return fmt.Errorf("remote_ip_address does not match %v - %v", ne.StringValue(bgpConfig.RemoteIPAddress), v)
		}
		if v, ok := ctx["bgp-remote_asn"]; ok && ne.IntValue(bgpConfig.RemoteASN) != v.(int) {
			return fmt.Errorf("remote_asn does not match %v - %v", ne.IntValue(bgpConfig.RemoteASN), v)
		}
		if v, ok := ctx["bgp-authentication_key"]; ok && ne.StringValue(bgpConfig.AuthenticationKey) != v.(string) {
			return fmt.Errorf("authentication_key does not match %v - %v", ne.StringValue(bgpConfig.AuthenticationKey), v)
		}
		if v, ok := ctx["connection-resourceName"]; ok {
			connResourceName := "equinix_ecx_l2_connection." + v.(string)
			rs, ok := s.RootModule().Resources[connResourceName]
			if !ok {
				return fmt.Errorf("related connection resource not found: %s", connResourceName)
			}
			if ne.StringValue(bgpConfig.ConnectionUUID) != rs.Primary.ID {
				return fmt.Errorf("connection_id does not match %v - %v", bgpConfig.ConnectionUUID, rs.Primary.ID)
			}
		}
		if v, ok := ctx["device-resourceName"]; ok {
			deviceResourceName := "equinix_network_device." + v.(string)
			rs, ok := s.RootModule().Resources[deviceResourceName]
			if !ok {
				return fmt.Errorf("related device resource not found: %s", deviceResourceName)
			}
			if ne.StringValue(bgpConfig.DeviceUUID) != rs.Primary.ID {
				return fmt.Errorf("device_id does not match %v - %v", bgpConfig.DeviceUUID, rs.Primary.ID)
			}
		}
		return nil
	}
}
