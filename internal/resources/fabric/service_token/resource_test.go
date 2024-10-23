package service_token_test

import (
	"context"
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/service_token"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
	"time"
)

func TestAccFabricServiceToken_PNFV(t *testing.T) {
	serviceTokenName, serviceTokenUpdatedName := "Service_token_PNFV", "UP_Service_Token_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "zside vd token", "Updated zside vd token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricServiceTokenConfig(serviceTokenName, serviceTokenDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2024-11-18T06:43:49.980Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", "fcf0fcec-65f6-4544-8810-ae4756fab8c4"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id", "5"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2024-11-18T06:43:49.980Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", "fcf0fcec-65f6-4544-8810-ae4756fab8c4"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id", "5"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricServiceTokenConfig(serviceTokenName string, serviceTokenDescription string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			name = "%s"
			description = "%s"
			expiration_date_time = "2024-11-18T06:43:49.98Z"
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
    `, serviceTokenName, serviceTokenDescription)
}

func CheckServiceTokenDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_service_token" {
			continue
		}

		err := service_token.WaitForDeletion(rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
