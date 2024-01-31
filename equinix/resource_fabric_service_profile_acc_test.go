package equinix_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
)

const (
	FabricDedicatedPortEnvVar = "TF_ACC_FABRIC_DEDICATED_PORTS"
)

type EnvPorts map[string]map[string][]v4.Port

func GetFabricEnvPorts(t *testing.T) EnvPorts {
	var ports EnvPorts
	portJson := os.Getenv(FabricDedicatedPortEnvVar)
	if err := json.Unmarshal([]byte(portJson), &ports); err != nil {
		t.Fatalf("Failed reading port data from environment: %v, %s", err, portJson)
	}
	return ports
}

func testAccFabricReadServiceProfileConfig(uuid string) string {
	return fmt.Sprintf(`data "equinix_fabric_service_profile" "test" {
	uuid = "%s"
	}
`, uuid)
}

func testAccFabricReadServiceProfilesListConfig(name string) string {
	return fmt.Sprintf(`data "equinix_fabric_service_profiles" "test" {
		filter {
			property = "/name"
			operator = "="
			values = ["%s"]
		}
	}
`, name)
}

func TestAccFabricReadServiceProfileByUuid_SP_FCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfileConfig("bfb74121-7e2c-4f74-99b3-69cdafb03b41"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "name", fmt.Sprint("Azure ExpressRoute")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "uuid", fmt.Sprint("bfb74121-7e2c-4f74-99b3-69cdafb03b41")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "state", fmt.Sprint("ACTIVE")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "visibility", fmt.Sprint("PUBLIC")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "access_point_type_configs.#", fmt.Sprint(1)),
					resource.TestMatchResourceAttr(
						"data.equinix_fabric_service_profile.test", "account.0.organization_name", regexp.MustCompile("^azure")),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "metros.0.display_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profile.test", "self_profile"),
				),
			},
		},
	})
}

func TestAccFabricSearchServiceProfilesByName_SP_FCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfilesListConfig("Azure ExpressRoute"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.#", fmt.Sprint(1)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.name", fmt.Sprint("Azure ExpressRoute")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.uuid", fmt.Sprint("bfb74121-7e2c-4f74-99b3-69cdafb03b41")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.type", fmt.Sprint("L2_PROFILE")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.state", fmt.Sprint("ACTIVE")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.visibility", fmt.Sprint("PUBLIC")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.account.#", fmt.Sprint(1)),
					resource.TestMatchResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.account.0.organization_name", regexp.MustCompile("^azure")),

					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.allow_remote_connections"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.allow_custom_bandwidth"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.enable_auto_generate_service_key"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.access_point_type_configs.0.connection_redundancy_required"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.metros.0.display_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_profiles.test", "data.0.self_profile"),
				),
			},
		},
	})
}

func TestAccFabricCreateServiceProfile_SP_FCR_b(t *testing.T) {
	spName1 := "fabric_tf_acc_test_CCEPL_01"
	spName2 := "fabric_tf_acc_test_CCEPL_02"
	ports := GetFabricEnvPorts(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkServiceProfileDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateServiceProfileConfig(spName1, ports["pfcr"]["dot1q"][0]["uuid"], ports["pfcr"]["dot1q"][0]["type"], ports["pfcr"]["dot1q"][0]["location"]["metroCode"]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", fmt.Sprint(spName1)),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricCreateServiceProfileConfig(spName2, ports["pfcr"]["qinq"][0]["uuid"], ports["pfcr"]["qinq"][0]["type"], ports["pfcr"]["qinq"][0]["location"]["metroCode"]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", fmt.Sprint(spName2)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricCreateServiceProfileConfig(name string, portUUID string, portType string, portMetroCode string) string {
	return fmt.Sprintf(`resource "equinix_fabric_service_profile" "test" {
  name = "%s"
  description = "Generic SP"
  type = "L2_PROFILE"
  notifications {
      emails = ["opsuser100@equinix.com"]
      type = "BANDWIDTH_ALERT"
      send_interval = ""
  }
  tags = ["Storage","Compute"]
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
}`, name, portUUID, portType, portMetroCode)
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
