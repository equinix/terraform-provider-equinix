package metro_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricMetroDataSource_PFCR(t *testing.T) {
	metroName := "Melbourne"
	metroCode := "ME"
	limit := 8
	offset := 6

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricMetroDataSourcesConfig(metroCode, limit, offset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.equinix_fabric_metro.metro", "code", metroCode),
					resource.TestCheckResourceAttr("data.equinix_fabric_metro.metro", "name", metroName),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "region"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "equinix_asn"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "local_vc_bandwidth_max"),
					resource.TestCheckResourceAttr("data.equinix_fabric_metro.metro", "geo_coordinates.%", "2"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "geo_coordinates.latitude"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "geo_coordinates.longitude"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "connected_metros.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "connected_metros.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "connected_metros.0.avg_latency"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metro.metro", "connected_metros.0.remote_vc_bandwidth_max"),

					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.region"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.equinix_asn"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.local_vc_bandwidth_max"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.geo_coordinates.latitude"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.geo_coordinates.longitude"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.connected_metros.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.connected_metros.0.code"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.connected_metros.0.avg_latency"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_metros.metros", "data.0.connected_metros.0.remote_vc_bandwidth_max"),
					resource.TestCheckResourceAttr("data.equinix_fabric_metros.metros", "pagination.%", "5"),
					resource.TestCheckResourceAttr("data.equinix_fabric_metros.metros", "pagination.limit", strconv.Itoa(limit)),
					resource.TestCheckResourceAttr("data.equinix_fabric_metros.metros", "pagination.offset", strconv.Itoa(offset)),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func testAccFabricMetroDataSourcesConfig(code string, limit, offset int) string {
	return fmt.Sprintf(`

       data "equinix_fabric_metro" "metro" {
		  metro_code = "%[1]s"
		}

      data "equinix_fabric_metros" "metros" {
  		pagination = {
    		limit = "%[2]d",
            offset = "%[3]d"
		}
	}
	`, code, limit, offset)
}
