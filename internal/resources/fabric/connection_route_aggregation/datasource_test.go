package connection_route_aggregation_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func testAccFabricConnectionRouteAggregationDataSourcesConfig() string {
	return fmt.Sprintf(`
		resource "equinix_fabric_connection_route_aggregation" "test" {
  			route_aggregation_id = "f264c892-bfd4-4823-ab8f-7ee74cc7a0f0"
  			connection_id = "daa424a5-e7b1-4fb6-bdbc-200e9b4757d1"
		}

		data "equinix_fabric_connection_route_aggregation" "data_cra" {
  			depends_on = [equinix_fabric_connection_route_aggregation.test]
  			route_aggregation_id = "f264c892-bfd4-4823-ab8f-7ee74cc7a0f0"
  			connection_id = "daa424a5-e7b1-4fb6-bdbc-200e9b4757d1"
		}


		data "equinix_fabric_connection_route_aggregations" "data_cras" {
  			depends_on = [equinix_fabric_connection_route_aggregation.test]
  			connection_id = "daa424a5-e7b1-4fb6-bdbc-200e9b4757d1"
		}`)
}

func TestAccFabricConnectionRouteAggregationDataSources_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckConnectionRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricConnectionRouteAggregationDataSourcesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_aggregation.data_cra", "attachment_status", "ATTACHED"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection_route_aggregation.data_cra", "type", "BGP_IPv4_PREFIX_AGGREGATION"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_aggregation.data_cra", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_aggregation.data_cra", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_aggregations.data_cras", "data.0.attachment_status"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_aggregations.data_cras", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_aggregations.data_cras", "data.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection_route_aggregations.data_cras", "data.0.uuid"),
					resource.TestCheckResourceAttr("data.equinix_fabric_connection_route_aggregations.data_cras", "data.#", "1"),
					resource.TestCheckResourceAttr("data.equinix_fabric_connection_route_aggregations.data_cras", "pagination.%", "5"),
					resource.TestCheckResourceAttr("data.equinix_fabric_connection_route_aggregations.data_cras", "pagination.limit", "10"),
					resource.TestCheckResourceAttr("data.equinix_fabric_connection_route_aggregations.data_cras", "pagination.offset", "0"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
