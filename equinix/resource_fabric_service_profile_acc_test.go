package equinix

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
)

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

func TestAccFabricReadServiceProfileByUuid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfileConfig("bfb74121-7e2c-4f74-99b3-69cdafb03b41"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "name", fmt.Sprint("Azure ExpressRoute")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profile.test", "uuid", fmt.Sprint("bfb74121-7e2c-4f74-99b3-69cdafb03b41")),
				),
			},
		},
	})
}

func TestAccFabricSearchServiceProfilesByName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfilesListConfig("Azure ExpressRoute"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.#", fmt.Sprint(1)), // Check  total number of ServiceProfile list returned in the Response payloads
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.name", fmt.Sprint("Azure ExpressRoute")),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_service_profiles.test", "data.0.uuid", fmt.Sprint("bfb74121-7e2c-4f74-99b3-69cdafb03b41")),
					resource.TestCheckNoResourceAttr(
						"data.equinix_fabric_service_profiles.test", "pagination"),
				),
			},
		},
	})
}

func TestAccFabricCreateServiceProfile(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkServiceProfileDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateServiceProfileConfig("fabric_tf_acc_test_CCEPL_01"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", fmt.Sprint("fabric_tf_acc_test_CCEPL_01")),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricCreateServiceProfileConfig("fabric_tf_acc_test_CCEPL_02"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_service_profile.test", "name", fmt.Sprint("fabric_tf_acc_test_CCEPL_02")),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricCreateServiceProfileConfig(name string) string {
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
      uuid = "c4d9350e-77c5-7c5d-1ce0-306a5c00a600"
      type = "XF_PORT"
      location {
        metro_code = "SV"
      }
      cross_connect_id = ""
      seller_region = ""
      seller_region_description = ""
  }
  access_point_type_configs {
      type= "COLO"
      connection_redundancy_required= false
      allow_bandwidth_auto_approval= false
      allow_remote_connections= false
      connection_label= "test"
      enable_auto_generate_service_key= false
      bandwidth_alert_threshold= 10
      allow_custom_bandwidth= true
      api_config {
        api_available= false
        equinix_managed_vlan= true
        bandwidth_from_api= false
        integration_id= "test"
        equinix_managed_port= true
      }
      authentication_key{
        required= false
        label= "Service Key"
        description= "XYZ"
      }
      supported_bandwidths= [100,500]
  }
  marketing_info {
    promotion = false
  }
}`, name)
}

func checkServiceProfileDelete(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).FabricClient
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*config.Config).FabricAuthToken)
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
