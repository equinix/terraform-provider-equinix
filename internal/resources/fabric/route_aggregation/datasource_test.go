package route_aggregation_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFabricRouteAggregationDataSourcesConfig(name, description string) string {
	return fmt.Sprintf(`

		resource "equinix_fabric_route_aggregation" "new_ra_1" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "%[1]s"
		  description = "%[2]s"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}

		resource "equinix_fabric_route_aggregation" "new_ra_2" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "%[1]s"
		  description = "%[2]s"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}

		data "equinix_fabric_route_aggregation" "data_ra" {
		 route_aggregation_id = equinix_fabric_route_aggregation.new_ra_2.id
		}

		data "equinix_fabric_route_aggregations" "data_ras" {
		 depends_on = [equinix_fabric_route_aggregation.new_ra_1, equinix_fabric_route_aggregation.new_ra_2]
		   filter =  {
   		property = "/type"
   		operator = "="
			values    = ["BGP_IPv4_PREFIX_AGGREGATION"]
 			}
		pagination = {
   		limit = 2
   		offset = 1
 		}
		 sort = {
   		property = "/changeLog/updatedDateTime"
   		direction = "DESC"
       }
		}
	`, name, description)
}

func TestAccFabricRouteAggregationDataSources_PFCR(t *testing.T) {
	routeAggregationName := "route_agg_PFCR"
	routeAggregatioDescription := "route_agg_PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteAggregationDataSourcesConfig(routeAggregationName, routeAggregatioDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation.data_ra", "name", routeAggregationName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation.data_ra", "type", "BGP_IPv4_PREFIX_AGGREGATION"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation.data_ra", "project.project_id", "4f855852-eb47-4721-8e40-b386a3676abf"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_route_aggregation.data_ra", "description", routeAggregatioDescription),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation.data_ra", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation.data_ra", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregation.data_ra", "change_log.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.project.project_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.connections_count"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.rules_count"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_route_aggregations.data_ras", "data.0.change_log.created_by"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregations.data_ras", "data.#", "2"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregations.data_ras", "pagination.%", "5"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregations.data_ras", "pagination.limit", "2"),
					resource.TestCheckResourceAttr("data.equinix_fabric_route_aggregations.data_ras", "pagination.offset", "1"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}
