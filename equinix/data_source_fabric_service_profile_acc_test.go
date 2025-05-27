package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFabricReadServiceProfileConfig(spName string, portUUID string, portType string, portMetroCode string) string {
	return fmt.Sprintf(`

resource "equinix_fabric_service_profile" "test" {
  name = "%s"
  description = "Generic Read SP"
  self_profile = false
  type = "L2_PROFILE"
  notifications {
      emails = ["opsuser100@equinix.com"]
      type = "BANDWIDTH_ALERT"
  }
  tags = ["VoIP", "Saas"]
  visibility = "PRIVATE"
  allowed_emails = ["panthersfcr@test.com"]
  ports {
      uuid = "%s"
      type = "%s"
      location {
        metro_code = "%s"
      }
      cross_connect_id = ""
      seller_region = ""
      seller_region_description = ""
  }
  access_point_type_configs {
      type = "COLO"
      connection_redundancy_required = false
      allow_bandwidth_auto_approval = false
      allow_remote_connections = false
      connection_label = "test"
      enable_auto_generate_service_key = false
      bandwidth_alert_threshold=  10
      allow_custom_bandwidth = true
      api_config {
        api_available = false
        equinix_managed_vlan = true
        bandwidth_from_api = false
        integration_id = "test"
        equinix_managed_port = true
      }
      authentication_key{
        required = false
        label = "Service Key"
        description = "XYZ"
      }
      supported_bandwidths = [100,500]
  }
  marketing_info {
    promotion = false
  }
}

data "equinix_fabric_service_profile" "test" {
		uuid = equinix_fabric_service_profile.test.uuid
}

data "equinix_fabric_service_profiles" "test" {
	and_filters = true
	filter {
		property = "/name"
		operator = "="
		values = [equinix_fabric_service_profile.test.name]
	}
}`, spName, portUUID, portType, portMetroCode)
}

func TestAccFabricServiceProfileDataSources_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)

	var portUUID, portMetroCode, portType string
	if len(ports) > 0 {
		port := ports["pfcr"]["dot1q"][0]
		portUUID = port.GetUuid()
		portMetroCodeLocation := port.GetLocation()
		portMetroCode = portMetroCodeLocation.GetMetroCode()
		portType = string(port.GetType())
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkServiceProfileDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfileConfig("SP_DataSource_PFCR", portUUID, portType, portMetroCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "name", "SP_DataSource_PFCR"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_service_profile.test", "uuid"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "description", "Generic Read SP"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "visibility", "PRIVATE"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "access_point_type_configs.#", "1"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_profile.test", "metros.0.code", portMetroCode),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.display_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "self_profile"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.name", "SP_DataSource_PFCR"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_service_profiles.test", "data.0.uuid"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.description", "Generic Read SP"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.visibility", "PRIVATE"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.#", "1"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_profiles.test", "data.0.metros.0.code", portMetroCode),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.display_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.self_profile"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
