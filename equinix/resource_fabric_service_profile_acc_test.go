package equinix_test

import (
	"context"
	"fmt"
	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricCreateServiceProfile_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)

	var portUuidDot1Q, portMetroCodeDot1Q, portTypeDot1Q string
	var portUuidQinq, portMetroCodeQinq, portTypeQinq string
	if len(ports) > 0 {
		portDot1Q := ports["pfcr"]["dot1q"][0]
		portQinq := ports["pfcr"]["qinq"][0]
		portUuidDot1Q = portDot1Q.GetUuid()
		portMetroCodeDot1QLocation := portDot1Q.GetLocation()
		portMetroCodeDot1Q = portMetroCodeDot1QLocation.GetMetroCode()
		portTypeDot1Q = string(portDot1Q.GetType())
		portUuidQinq = portQinq.GetUuid()
		portMetroCodeQinqLocation := portQinq.GetLocation()
		portMetroCodeQinq = portMetroCodeQinqLocation.GetMetroCode()
		portTypeQinq = string(portQinq.GetType())
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkServiceProfileDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateServiceProfileConfig(portUuidDot1Q, portTypeDot1Q, portMetroCodeDot1Q),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", "SP_ResourceCreation_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "type", "L2_PROFILE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "visibility", "PRIVATE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "tags.#", "2"),

					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "visibility"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "metros.0.code", portMetroCodeDot1Q),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "metros.0.name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "metros.0.display_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "self_profile"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricCreateServiceProfileConfig(portUuidQinq, portTypeQinq, portMetroCodeQinq),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", "SP_ResourceCreation_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "type", "L2_PROFILE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "visibility", "PRIVATE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "tags.#", "2"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "visibility"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "metros.0.code", portMetroCodeQinq),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "metros.0.name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "metros.0.display_name"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "self_profile"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricCreateServiceProfileConfig(portUUID string, portType string, portMetroCode string) string {
	return fmt.Sprintf(`resource "equinix_fabric_service_profile" "test" {
  name = "SP_ResourceCreation_PFCR"
  description = "Generic SP"
  type = "L2_PROFILE"
  notifications {
      emails = ["opsuser100@equinix.com"]
      type = "BANDWIDTH_ALERT"
  }
  tags = ["VoIP", "Saas"]
  visibility = "PRIVATE"
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
      supported_bandwidths = [500]
  }
  marketing_info {
    promotion = false
  }
}`, portUUID, portType, portMetroCode)
}

func checkServiceProfileDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_service_profile" {
			continue
		}
		err := equinix.WaitAndCheckServiceProfileDeleted(rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion: %v", err)
		}
	}
	return nil
}
