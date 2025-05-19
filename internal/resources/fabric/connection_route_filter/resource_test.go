package connection_route_filter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection_route_filter"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricConnectionRouteFilter_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUUID string
	if len(ports) > 0 {
		portUUID = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricConnectionRouteFilterConfig(portUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection_route_filter.test", "id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_connection_route_filter.test", "connection_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_connection_route_filter.test", "route_filter_id"),

					resource.TestCheckResourceAttr(
						"equinix_fabric_connection_route_filter.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection_route_filter.test", "type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection_route_filter.test", "attachment_status", string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION)),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_filter.test", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_filter.test", "connection_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_filter.test", "route_filter_id"),

					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_filter.test", "direction", "INBOUND"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_filter.test", "type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_filter.test", "attachment_status", string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION)),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_filters.test", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_filters.test", "connection_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_filters.test", "data.0.uuid"),

					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_filters.test", "data.0.direction", "INBOUND"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_filters.test", "data.0.type", "BGP_IPv4_PREFIX_FILTER"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_filters.test", "data.0.attachment_status", string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION)),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricConnectionRouteFilterConfig(portUUID string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_cloud_router" "test" {
			type = "XF_ROUTER"
			name = "RF_CR_PFCR"
			location {
				metro_code  = "DC"
			}
			package {
				code = "STANDARD"
			}
			order {
				purchase_order_number = "1-234567"
                term_length = 1
			}
			notifications {
				type = "ALL"
				emails = [
					"test@equinix.com",
					"test1@equinix.com"
				]
			}
			project {
				project_id = "291639000636552"
			}
			account {
				account_number = 201257
			}
		}

		resource "equinix_fabric_connection" "test" {
			type = "IP_VC"
			name = "RF_CR_Connection_PFCR"
			notifications {
				type = "ALL"
				emails = ["test@equinix.com","test1@equinix.com"]
			}
			order {
				purchase_order_number = "123485"
                term_length = 1
			}
			bandwidth = 50
			redundancy {
				priority= "PRIMARY"
			}
			a_side {
				access_point {
					type = "CLOUD_ROUTER"
					router {
						uuid = equinix_fabric_cloud_router.test.id
					}
				}
			}
			project {
			   project_id = "291639000636552"
			}
			z_side {
				access_point {
					type = "COLO"
					port{
						uuid = "%s"
					}
					link_protocol {
						type= "DOT1Q"
<<<<<<< Updated upstream
						vlan_tag= 1500
=======
						vlan_tag= 2100
>>>>>>> Stashed changes
					}
					location {
						metro_code = "DC"
					}
				}
			}
		}
		
		resource "equinix_fabric_route_filter" "test" {
			name = "rf_test_PFCR"
			project {
				project_id = "291639000636552"
			}
			type = "BGP_IPv4_PREFIX_FILTER"
			description = "Route Filter Policy for X Purpose"
		}

		resource "equinix_fabric_route_filter_rule" "test" {
			route_filter_id = equinix_fabric_route_filter.test.id
			name = "RF_Rule_PFCR"
			prefix = "192.168.0.0/24"
			prefix_match = "exact"
			description = "Route Filter Rule for X Purpose"
		}

		resource "equinix_fabric_connection_route_filter" "test" {
			depends_on = [ equinix_fabric_route_filter_rule.test ]
			connection_id = equinix_fabric_connection.test.id
			route_filter_id = equinix_fabric_route_filter.test.id
			direction = "INBOUND"
		}

		data "equinix_fabric_connection_route_filter" "test" {
			depends_on = [ equinix_fabric_connection_route_filter.test ]
			connection_id = equinix_fabric_connection.test.id
			route_filter_id = equinix_fabric_route_filter.test.id
		}

		data "equinix_fabric_connection_route_filters" "test" {
			depends_on = [ equinix_fabric_connection_route_filter.test ]
			connection_id = equinix_fabric_connection.test.id
		}

	`, portUUID)
}

func CheckConnectionRouteFilterDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection_route_filter" {
			continue
		}

		connectionID := rs.Primary.Attributes["connection_id"]

		err := connection_route_filter.WaitForDeletion(connectionID, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
