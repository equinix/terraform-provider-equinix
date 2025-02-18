package connection_route_aggregation_test

import (
	"context"
	"fmt"
	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func testAccFabricConnectionRouteAggregationConfig(portUuid string) string {
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
			}
			notifications {
				type = "ALL"
				emails = [
					"test@equinix.com",
					"test1@equinix.com"
				]
			}
			project {
				project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
			}
			account {
				account_number = 77733367
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
			   project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
			}
			z_side {
				access_point {
					type = "COLO"
					port{
						uuid = "%s"
					}
					link_protocol {
						type= "DOT1Q"
						vlan_tag= 2571
					}
					location {
						metro_code = "DC"
					}
				}
			}
		}
		
		//resource "equinix_fabric_route_aggregation" "test" {
		//  type = "BGP_IPv4_PREFIX_AGGREGATION"
		//  name = "Route_Aggregation_Test"
		//  description = "Test Route Aggregation"
		//  project = {
		//	project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		//	}
		//}

		resource "equinix_fabric_connection_route_aggregation" "test" {
			route_aggregation_id = "8f8a2ddb-25f8-416e-ad0a-202a9d2af9e1"
			connection_id = equinix_fabric_connection.test.id
		}
	`, portUuid)
}

func TestAccFabricConnectionRouteAggregation_PFCR(t *testing.T) {
	portId := "c5720fcc-4ae6-ae6e-13e0-306a5c00adaf"
	//upRouteAggregationName := "stream_up_PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckConnectionRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricConnectionRouteAggregationConfig(portId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection_route_aggregation.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_connection_route_aggregation.test", "attachment_status"),
					resource.TestCheckResourceAttrSet("equinix_fabric_connection_route_aggregation.test", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_connection_route_aggregation.test", "type", "BGP_IPv4_PREFIX_AGGREGATION"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func CheckConnectionRouteAggregationDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection_route_aggregation" {
			continue
		}

		routeAggregationId := rs.Primary.Attributes["route_aggregation_id"]
		connectionId := rs.Primary.Attributes["connection_id"]

		if connectionRouteAggregation, _, err := client.RouteAggregationsApi.GetConnectionRouteAggregationByUuid(ctx, routeAggregationId, connectionId).Execute(); err == nil &&
			connectionRouteAggregation.GetAttachmentStatus() == fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHED {
			return fmt.Errorf("fabric connection route aggregation attchement %s still exists and is %s",
				rs.Primary.ID, string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_ATTACHED))
		}
	}
	return nil
}
