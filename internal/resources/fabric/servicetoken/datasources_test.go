package servicetoken_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"time"
)

func TestAccFabricServiceTokenDataSource_PNFV(t *testing.T) {
	connectionTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var virtualDevice string
	if len(connectionTestData) > 0 {
		virtualDevice = connectionTestData["pnfv"]["virtualDevice"]
	}
	expiration := time.Now().Add(30 * 24 * time.Hour).UTC().Format("2006-01-02T15:04:05.000Z")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricServiceTokenConfigDataSourceConfig(expiration, virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "expiration_date_time", expiration),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_token.service-token-for-zside-virtual-device-for-zside-virtual-device", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_tokens.service-tokens", "data.0.uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.expiration_date_time", expiration),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccFabricServiceTokenConfigDataSourceConfig(expiration string, virtualDeviceUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			expiration_date_time = "%s"
			service_token_connection {
				type = "EVPL_VC"
				supported_bandwidths = [50, 200, 10000]
				z_side {
					access_point_selectors{
						type = "VD"
						virtual_device{
							type = "EDGE"
							uuid = "%s"
						}
						interface{
							type = "NETWORK"
						}
					}
				}
			}
			notifications {
    			type   = "ALL"
    			emails = ["panthers_auto@equinix.com", "test1@equinix.com", "example@equinix.com"]
  			}
            project {
    			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
  			}
		}
		
		data "equinix_fabric_service_token" "service-token-for-zside-virtual-device-for-zside-virtual-device" {
			uuid = equinix_fabric_service_token.test.id
		}
		
		data "equinix_fabric_service_tokens" "service-tokens"{
			filter {
				property = "/uuid"
				operator = "="
				values 	 = [equinix_fabric_service_token.test.id]
			}
			filter {
				property = "/state"
				operator = "="
				values 	 = ["INACTIVE"]
			}
			pagination {
				offset = 0
				limit = 5
				total = 25
			}
		}
		
    `, expiration, virtualDeviceUUID)
}
