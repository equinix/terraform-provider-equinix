package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkBGP_CSR1000V_Single_AWS(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SV")()
	spName, _ := schema.EnvDefaultFunc(awsSpEnvVar, "AWS Direct Connect")()
	context := map[string]interface{}{
		"device-resourceName":          "test",
		"device-self_managed":          true,
		"device-byol":                  true,
		"device-name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"device-throughput":            500,
		"device-throughput_unit":       "Mbps",
		"device-metro_code":            metro.(string),
		"device-type_code":             "CSR1000V",
		"device-package_code":          "SEC",
		"device-notifications":         []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":              fmt.Sprintf("tf-%s", randString(6)),
		"device-term_length":           1,
		"device-version":               "16.09.05",
		"device-core_count":            2,
		"sshkey-resourceName":          "test",
		"sshkey-name":                  fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"sshkey-public_key":            "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQCXdzXBHaVpKpdO0udnB+4JOgUq7APO2rPXfrevvlZrps98AtlwXXVWZ5duRH5NFNfU4G9HCSiAPsebgjY0fG85tcShpXfHfACLt0tBW8XhfLQP2T6S50FQ1brBdURMDCMsD7duOXqvc0dlbs2/KcswHvuUmqVzob3bz7n1bQ48wIHsPg4ARqYhy5LN3OkllJH/6GEfqi8lKZx01/P/gmJMORcJujuOyXRB+F2iXBVYdhjML3Qg4+tEekBcVZOxUbERRZ0pvQ52Y6wUhn2VsjljixyqeOdmD0m6DayDQgSWms6bKPpBqN7zhXXk4qe8bXT4tQQba65b2CQ2A91jw2KgM/YZNmjyUJ+Rf1cQosJf9twqbAZDZ6rAEmj9zzvQ5vD/CGuzxdVMkePLlUK4VGjPu7cVzhXrnq4318WqZ5/lNiCST8NQ0fssChN8ANUzr/p/wwv3faFMVNmjxXTZMsbMFT/fbb2MVVuqNFN65drntlg6/xEao8gZROuRYiakBx8= user@host",
		"connection-resourceName":      "test",
		"connection-name":              fmt.Sprintf("%s-%s", tstResourcePrefix, randString(6)),
		"connection-profile_name":      spName.(string),
		"connection-speed":             50,
		"connection-speed_unit":        "MB",
		"connection-notifications":     []string{"marry@equinix.com", "john@equinix.com"},
		"connection-seller_metro_code": "SV",
		"connection-seller_region":     "us-west-1",
		"connection-authorization_key": "123456789012",
		"bgp-resourceName":             "test",
		"bgp-local_ip_address":         "1.1.1.1/30",
		"bgp-local_asn":                12345,
		"bgp-remote_ip_address":        "1.1.1.2",
		"bgp-remote_asn":               22211,
	}
	contextWithChanges := copyMap(context)
	contextWithChanges["bgp-authentication_key"] = randString(10)
	resourceName := fmt.Sprintf("equinix_network_bgp.%s", context["bgp-resourceName"].(string))
	var bgpConfig ne.BGPConfiguration
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: newTestAccConfig(context).withDevice().withSSHKey().withConnection().withBGP().build(),
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
				Config: newTestAccConfig(contextWithChanges).withDevice().withSSHKey().withConnection().withBGP().build(),
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
	config += nprintf(`
resource "equinix_network_bgp" "%{bgp-resourceName}" {
  connection_id      = equinix_ecx_l2_connection.%{connection-resourceName}.id
  local_ip_address   = "%{bgp-local_ip_address}"
  local_asn          = %{bgp-local_asn}
  remote_ip_address  = "%{bgp-remote_ip_address}"
  remote_asn         = %{bgp-remote_asn}`, ctx)
	if _, ok := ctx["bgp-authentication_key"]; ok {
		config += nprintf(`
  authentication_key = "%{bgp-authentication_key}"`, ctx)
	}
	config += `
}`
	if _, ok := ctx["connection-secondary_name"]; ok {
		config += nprintf(`
resource "equinix_network_bgp" "%{bgp-secondary_resourceName}" {
  connection_id      = equinix_ecx_l2_connection.%{connection-resourceName}.id
  local_ip_address   = "%{bgp-secondary_local_ip_address}"
  local_asn          = %{bgp-secondary_local_asn}
  remote_ip_address  = "%{bgp-secondary_remote_ip_address}"
  remote_asn         = %{bgp-secondary_remote_asn}`, ctx)
		if _, ok := ctx["bgp-secondary_authentication_key"]; ok {
			config += nprintf(`
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
