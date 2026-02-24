package advertisedroute_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricDataSourceAdvertisedRoutes_PFCR(t *testing.T) {
	offset := 6
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricAdvertisedRoutesDataSourcesConfig(offset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.type"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.protocol_type"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.prefix"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.next_hop"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.as_path.0"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.connection.uuid"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.connection.href"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.connection.name"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.change_log.created_date_time"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_advertised_routes.routes", "data.0.change_log.updated_date_time"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccFabricAdvertisedRoutesDataSourcesConfig(offset int) string {
	return fmt.Sprintf(`

	data "equinix_fabric_advertised_routes" "routes" {
		connection_id = "3946618a-4834-4acb-b2b1-3b6d9c634fcc"
		   filter =  {
   		property = "/type"
   		operator = "IN"
			values    = ["IPv4_BGP_ROUTE"]
 			}
		pagination = {
   		limit = 100
   		offset = "%[1]d"
 		}
		sort = {
   		property = "/changeLog/updatedDateTime"
   		direction = "DESC"
       }
		}
	`, offset)
}
