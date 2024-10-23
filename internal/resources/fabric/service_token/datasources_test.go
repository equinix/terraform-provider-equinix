package service_token_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccFabricServiceTokenDataSource_PNFV(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricServiceTokenConfigDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_token.service-token", "uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "expiration_date_time", "2024-11-18T06:43:49.980Z"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", "fcf0fcec-65f6-4544-8810-ae4756fab8c4"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_token.service-token", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id", "5"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_service_tokens.service-tokens", "data.0.uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.expiration_date_time", "2024-11-18T06:43:49.980Z"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", "fcf0fcec-65f6-4544-8810-ae4756fab8c4"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttr("data.equinix_fabric_service_tokens.service-tokens", "data.0.service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id", "5"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricServiceTokenConfigDataSourceConfig() string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			expiration_date_time = "2024-11-18T06:43:49.980Z"
			service_token_connection {
				type = "EVPL_VC"
				supported_bandwidths = [50, 200, 10000]
				z_side {
					access_point_selectors{
						type = "VD"
						virtual_device{
							type = "EDGE"
							uuid = "fcf0fcec-65f6-4544-8810-ae4756fab8c4"
						}
						interface{
							type = "NETWORK"
							id = 5
						}
					}
				}
			}
			notifications {
    			type   = "ALL"
    			emails = ["example@equinix.com", "test1@equinix.com"]
  			}
		}
		
		data "equinix_fabric_service_token" "service-token" {
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
		
    `)
}
