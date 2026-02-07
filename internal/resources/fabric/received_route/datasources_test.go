package received_route_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricDataSourceReceivedRoutes_PFCR(t *testing.T) {
	limit := 8
	offset := 6
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReceivedRoutesDataSourcesConfig(limit, offset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_received_routes.received_route", "id"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "type"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "protocol_type"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "prefix"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "next_hop"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "med"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "local_preference"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "as_path"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "connection.0.uuid"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "connection.0.href"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "connection.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "change_log.0.created_by"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.received_route", "change_log.0.created_by_full_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connections.connections", "change_log.0.createdByEmail"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.createdDateTime"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.updatedBy"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.updatedByFullName"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.updatedByEmail"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.updatedDateTime"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.deletedBy"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.deletedByFullName"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.deletedByEmail"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_received_routes.connections", "change_log.0.deletedDateTime"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricReceivedRoutesDataSourcesConfig(limit, offset int) string {
	return fmt.Sprintf(`

	data "equinix_received_routes" "routes" {
		connection_id = "6b6fde52-843f-475d-a252-2c9b294aa70d"
		   filter =  {
   		property = "/type"
   		operator = "IN"
			values    = ["IPv4_BGP_ROUTE"]
 			}
		pagination = {
   		limit = 100
   		offset = 0
 		}
		sort = {
   		property = "/changeLog/updatedDateTime"
   		direction = "DESC"
       }
		}
	`)
}
