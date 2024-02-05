package equinix_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
)

func TestAccFabricCreateServiceProfile_SP_FCR_ba(t *testing.T) {
	spName1 := "fabric_tf_acc_test_01"
	spName2 := "fabric_tf_acc_test_02"
	typel2 := "L2_PROFILE"
	portType := "XF_PORT"
	ports := GetFabricEnvPorts(t)
	portTcFirst := ports["pfcr"]["dot1q"][0]
	portTcSecond := ports["pfcr"]["qinq"][0]
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkServiceProfileDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateServiceProfileConfig(spName1, typel2, portTcFirst.Uuid, portType, portTcFirst.Location.MetroCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", fmt.Sprint(spName1)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "type", fmt.Sprint(typel2)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "state", fmt.Sprint("ACTIVE")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "visibility", fmt.Sprint("PRIVATE")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "tags.#", fmt.Sprint("2")),

					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "visibility"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "metros.0.code", portTcFirst.Location.MetroCode),
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
				Config: testAccFabricCreateServiceProfileConfig(spName2, typel2, portTcSecond.Uuid, portType, portTcSecond.Location.MetroCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", fmt.Sprint(spName2)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "type", fmt.Sprint(typel2)),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "state", fmt.Sprint("ACTIVE")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "visibility", fmt.Sprint("PRIVATE")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "tags.#", fmt.Sprint("2")),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "visibility"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "metros.0.code", portTcSecond.Location.MetroCode),
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

func testAccFabricCreateServiceProfileConfig(name string, typel2 string, portUUID string, portType string, portMetroCode string) string {
	return fmt.Sprintf(`resource "equinix_fabric_service_profile" "test" {
  name = "%s"
  description = "Generic SP"
  type = "%s"
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
      supported_bandwidths = [100,500]
  }
  marketing_info {
    promotion = false
  }
}`, name, typel2, portUUID, portType, portMetroCode)
}

func checkServiceProfileDelete(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).FabricClient
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, acceptance.TestAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_service_profile" {
			continue
		}
		_, err := waitAndCheckServiceProfileDeleted(rs.Primary.ID, client, ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func waitAndCheckServiceProfileDeleted(uuid string, client *v4.APIClient, ctx context.Context) (v4.ServiceProfile, error) {
	log.Printf("Waiting for service profile to be in deleted, uuid %s", uuid)
	stateConf := &retry.StateChangeConf{
		Target: []string{string(v4.DELETED_ServiceProfileStateEnum)},
		Refresh: func() (interface{}, string, error) {
			dbConn, _, err := client.ServiceProfilesApi.GetServiceProfileByUuid(ctx, uuid, nil)
			if err != nil {
				return "", "", err
			}
			updatableState := ""
			if *dbConn.State == v4.DELETED_ServiceProfileStateEnum {
				updatableState = string(*dbConn.State)
			}
			return dbConn, updatableState, nil
		},
		Timeout:    1 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	inter, err := stateConf.WaitForStateContext(ctx)
	dbConn := v4.ServiceProfile{}

	if err == nil {
		dbConn = inter.(v4.ServiceProfile)
	}
	return dbConn, err
}
