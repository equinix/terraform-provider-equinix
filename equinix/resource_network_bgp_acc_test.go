package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	awsSpUuid = "TF_ACC_NETWORK_FABRIC_SERVICE_UUID"
)

func TestAccNetworkBGP_CSR1000V_Single_AWS(t *testing.T) {
	metro, _ := schema.EnvDefaultFunc(networkDeviceMetroEnvVar, "SY")()
	spName, _ := schema.EnvDefaultFunc(awsSpEnvVar, "AWS Direct Connect")()
	spUuid, _ := schema.EnvDefaultFunc(awsSpUuid, "")()
	authKey, _ := schema.EnvDefaultFunc(awsAuthKeyEnvVar, "123456789012")()
	accountName, _ := schema.EnvDefaultFunc(networkDeviceAccountNameEnvVar, "")()
	projectId, _ := schema.EnvDefaultFunc(networkDeviceProjectId, "")()
	context := map[string]interface{}{
		"device-resourceName":           "test",
		"device-account_name":           accountName.(string),
		"device-self_managed":           false,
		"device-byol":                   false,
		"device-name":                   fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"device-throughput":             500,
		"device-throughput_unit":        "Mbps",
		"device-metro_code":             metro.(string),
		"device-project_id":             projectId.(string),
		"device-type_code":              "CSR1000V",
		"device-package_code":           "SEC",
		"device-notifications":          []string{"marry@equinix.com", "john@equinix.com"},
		"device-hostname":               fmt.Sprintf("tf-%s", acctest.RandString(6)),
		"device-term_length":            1,
		"device-version":                "16.09.05",
		"device-core_count":             2,
		"user-resourceName":             "test",
		"user-username":                 fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"user-password":                 acctest.RandString(10),
		"fabric-service-profile-uuid":   spUuid.(string),
		"connection-resourceName":       "test",
		"connection-name":               fmt.Sprintf("%s-%s", tstResourcePrefix, acctest.RandString(6)),
		"connection-profile_name":       spName.(string),
		"connection-bandwidth":          50,
		"connection-notifications_type": "ALL",
		"connection-seller_metro_code":  "SV",
		"connection-type":               "EVPL_VC",
		"connection-seller_region":      "us-west-1",
		"connection-authorization_key":  authKey.(string),
		"bgp-resourceName":              "test",
		"bgp-local_ip_address":          "1.1.1.1/30",
		"bgp-local_asn":                 12345,
		"bgp-remote_ip_address":         "1.1.1.2",
		"bgp-remote_asn":                22211,
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
				Config: newTestAccConfig(context).withDevice().withSSHUser().withVDConnection().withBGP().build(),
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
				Config: newTestAccConfig(contextWithChanges).withDevice().withSSHUser().withVDConnection().withBGP().build(),
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
  connection_id      = equinix_fabric_connection.%{connection-resourceName}.id
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
  connection_id      = equinix_fabric_connection.%{connection-resourceName}.id
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
			connResourceName := "equinix_fabric_connection." + v.(string)
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

func testAccVDFabricL2Connection(ctx map[string]interface{}) string {
	var config string
	if _, ok := ctx["zside-service_token"]; !ok {
		if _, ok := ctx["connection-profile_uuid"]; !ok {
			config += nprintf(`
data "equinix_fabric_service_profile" "pri" {
uuid = "%{fabric-service-profile-uuid}"
}`, ctx)
		}
	}
	if _, ok := ctx["connection-secondary_profile_name"]; ok {
		config += nprintf(`
data "equinix_fabric_service_profile" "sec" {
uuid = "%{fabric-service-profile-uuid}"
}`, ctx)
	}

	config += nprintf(`
resource "equinix_fabric_connection" "%{connection-resourceName}" {
  name                  = "%{connection-name}"
  type = "EVPL_VC"
  bandwidth                 = %{connection-bandwidth}
notifications {
    type   = "ALL"
    emails = %{device-notifications}
  }
a_side {
    access_point {
      type = "VD"
      virtual_device {
        type = "EDGE"
        uuid = equinix_network_device.%{device-resourceName}.id
      }
      interface {
        type = "CLOUD"
        id = 7
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      authentication_key = "%{connection-authorization_key}"
      seller_region = "%{connection-seller_region}"
      profile {
        type = "L2_PROFILE"
        uuid = "%{fabric-service-profile-uuid}"
      }
      location {
        metro_code = "%{connection-seller_metro_code}"
      }
    }
  }`, ctx)

	if _, ok := ctx["service_token"]; ok {
		config += nprintf(`
  service_token         = "%{service_token}"`, ctx)
	}
	if _, ok := ctx["zside-service_token"]; ok {
		config += nprintf(`
  zside_service_token   = "%{zside-service_token}"`, ctx)
	}
	if _, ok := ctx["zside-port_uuid"]; ok {
		config += nprintf(`
  zside_port_uuid       = "%{zside-port_uuid}"`, ctx)
	}
	if _, ok := ctx["connection-purchase_order_number"]; ok {
		config += nprintf(`
  purchase_order_number = "%{connection-purchase_order_number}"`, ctx)
	}

	if _, ok := ctx["port-uuid"]; ok {
		config += nprintf(`
  port_uuid             = "%{port-uuid}"`, ctx)
	} else if _, ok := ctx["port-resourceName"]; ok {
		config += nprintf(`
  port_uuid             = data.equinix_ecx_port.%{port-resourceName}.id`, ctx)
	}
	if _, ok := ctx["connection-vlan_stag"]; ok {
		config += nprintf(`
  vlan_stag             = %{connection-vlan_stag}`, ctx)
	}
	if _, ok := ctx["connection-vlan_ctag"]; ok {
		config += nprintf(`
  vlan_ctag             = %{connection-vlan_ctag}`, ctx)
	}
	if _, ok := ctx["connection-named_tag"]; ok {
		config += nprintf(`
  named_tag             = "%{connection-named_tag}"`, ctx)
	}
	if _, ok := ctx["connection-device_interface_id"]; ok {
		config += nprintf(`
  device_interface_id   = %{connection-device_interface_id}`, ctx)
	}
	if _, ok := ctx["connection-secondary_name"]; ok {
		config += nprintf(`
  secondary_connection {
    name                = "%{connection-secondary_name}"`, ctx)
		if _, ok := ctx["connection-secondary_profile_name"]; ok {
			config += nprintf(`
    profile_uuid        = data.equinix_fabric_sellerprofile.sec.id`, ctx)
		}
		if _, ok := ctx["secondary-port_uuid"]; ok {
			config += nprintf(`
	port_uuid             = "%{secondary-port_uuid}"`, ctx)
		} else if _, ok := ctx["port-secondary_resourceName"]; ok {
			config += nprintf(`
    port_uuid           = data.equinix_ecx_port.%{port-secondary_resourceName}.id`, ctx)
		}
		if _, ok := ctx["device-secondary_name"]; ok {
			config += nprintf(`
    device_uuid         = equinix_network_device.%{device-resourceName}.redundant_id`, ctx)
		}
		if _, ok := ctx["connection-secondary_vlan_stag"]; ok {
			config += nprintf(`
    vlan_stag           = %{connection-secondary_vlan_stag}`, ctx)
		}
		if _, ok := ctx["connection-secondary_vlan_ctag"]; ok {
			config += nprintf(`
    vlan_ctag           = %{connection-secondary_vlan_ctag}`, ctx)
		}
		if _, ok := ctx["connection-secondary_device_interface_id"]; ok {
			config += nprintf(`
    device_interface_id = %{connection-secondary_device_interface_id}`, ctx)
		}
		if _, ok := ctx["connection-secondary_speed"]; ok {
			config += nprintf(`
    speed               = %{connection-secondary_speed}`, ctx)
		}
		if _, ok := ctx["connection-secondary_speed_unit"]; ok {
			config += nprintf(`
    speed_unit          = "%{connection-secondary_speed_unit}"`, ctx)
		}
		if _, ok := ctx["connection-secondary_seller_metro_code"]; ok {
			config += nprintf(`
    seller_metro_code   = "%{connection-secondary_seller_metro_code}"`, ctx)
		}
		if _, ok := ctx["connection-secondary_seller_region"]; ok {
			config += nprintf(`
    seller_region       = "%{connection-secondary_seller_region}"`, ctx)
		}
		if _, ok := ctx["connection-secondary_authorization_key"]; ok {
			config += nprintf(`
    authorization_key   = "%{connection-secondary_authorization_key}"`, ctx)
		}
		if _, ok := ctx["secondary-service_token"]; ok {
			config += nprintf(`
    service_token       = "%{secondary-service_token}"`, ctx)
		}
		config += `
 	}`
	}
	config += `
}`
	return config
}

func (t *testAccConfig) withVDConnection() *testAccConfig {
	t.config += testAccVDFabricL2Connection(t.ctx)
	return t
}
