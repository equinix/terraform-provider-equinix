package route_aggregation_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccFabricRouteAggregationConfig(name string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_route_aggregation" "test" {
		  type = "BGP_IPv4_PREFIX_AGGREGATION"
		  name = "%s"
		  description = "Test Route Aggregation"
		  project = {
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		  }
		}
	`, name)
}

func TestAccFabricRouteAggregation_PFCR(t *testing.T) {
	routeAggregationName := "stream_PFCR"
	//upRouteAggregationName := "stream_up_PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricRouteAggregationConfig(routeAggregationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_route_aggregation.test", "name", routeAggregationName),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "state"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_route_aggregation.test", "project.project_id"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "name", routeAggregationName),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "type", "BGP_IPv4_PREFIX_AGGREGATION"),
					resource.TestCheckResourceAttr("equinix_fabric_route_aggregation.test", "description", "Test Route Aggregation"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func CheckRouteAggregationDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_route_aggregation" {
			continue
		}

		_, resp, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, rs.Primary.ID).Execute()
		if err != nil {
			// Check if the response exists and contains status 400 or 404
			if resp != nil && (resp.StatusCode == 400 || resp.StatusCode == 404) {
				fmt.Printf("Resource %s not found, treating as deleted\n", rs.Primary.ID)
				return nil
			}

			// Handle specific API error messages
			var apiErr *fabricv4.GenericOpenAPIError
			if errors.As(err, &apiErr) {
				errorBody := apiErr.Body()
				var errorResponse map[string]interface{}
				if jsonErr := json.Unmarshal(errorBody, &errorResponse); jsonErr == nil {
					if errorCode, exists := errorResponse["errorCode"]; exists && errorCode == "EQ-3044301" {
						fmt.Printf("Detected EQ-3044301 for resource %s, treating as deleted\n", rs.Primary.ID)
						return nil // Successfully handled the expected deletion case
					}
				}
			}

			return fmt.Errorf("unexpected API error checking deletion: %v", err)
		}

		if routeAggregation, _, err := client.RouteAggregationsApi.GetRouteAggregationByUuid(ctx, rs.Primary.ID).Execute(); err == nil {
			if routeAggregation.GetState() == fabricv4.ROUTEAGGREGATIONSTATE_PROVISIONED {
				return fmt.Errorf("fabric route aggregation %s still exists and is %s",
					rs.Primary.ID, routeAggregation.GetState())
			}
		}
	}
	return nil
}
