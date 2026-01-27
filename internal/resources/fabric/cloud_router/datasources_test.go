package cloud_router

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricDataSourceAdvertisedRoutes_PFCR(t *testing.T) {
	limit := 8
	offset := 6
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,// what is this based on?
		Steps: []resource.TestStep{
			{
				Config: testAccFabricAdvertisedRoutesDataSourcesConfig(limit, offset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_advertised_routes.advertised_route", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "type", "IPv4_BGP_ROUTE"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "protocolType", "BGP"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "state", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "prefix", "prefix"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "nextHop", "10.1.1.1/24"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "MED", "1"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "localPreference", "2"),
					resource.TestCheckResourceAttr(
						"data.equinix_advertised_routes.advertised_route", "asPath", ["aspath1", "aspath2"]),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.advertised_route", "connection.0.uuid"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.advertised_route", "connection.0.href"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.advertised_route", "connection.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.advertised_route", "changeLog.0.createdBy"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.advertised_route", "changeLog.0.createdByFullName"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connections.connections", "changeLog.0.createdByEmail"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.createdDateTime"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.updatedBy"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.updatedByFullName"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.updatedByEmail"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.updatedDateTime"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.deletedBy"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.deletedByFullName"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.deletedByEmail"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_advertised_routes.connections", "changeLog.0.deletedDateTime"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricAdvertisedRoutesDataSourcesConfig(limit, offset int) string {
	return fmt.Sprintf(`
	
	data "equinix_fabric_connection" "test" {
		type = "IPv4_BGP_ROUTE"
		protocolType = "BGP"
		state = "ACTIVE"
		prefix = "prefix"
		nextHop = "10.1.1.1/24"
		MED = "1"
		localPreference =  "2"
		asPath =  ["aspath1", "aspath2"]

	}
	data "equinix_advertised_routes" "routes" {
		pagination = {
				limit = "%[1]d",
				offset = "%[2]d"
			}
		outer_operator = "AND"
		filter {
			property = "/type"
			operator = "="
			values = ["IPv4_BGP_ROUTE"]
		}
		filter {
			property = "/state"
			operator = "="
			values = ["ACTIVE"]
		}
	}

	`,limit, offset)
}
