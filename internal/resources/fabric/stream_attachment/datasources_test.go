package stream_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccFabricStreamDataSourcesConfig(name, description string) string {
	return fmt.Sprintf(`
	
		resource "equinix_fabric_stream" "new_stream_1" {
		  type = "TELEMETRY_STREAM"
		  name = "%[1]s"
		  description = "%[2]s"
		}

		resource "equinix_fabric_stream" "new_stream_2" {
		  type = "TELEMETRY_STREAM"
		  name = "%[1]s"
		  description = "%[2]s"
		}
		
		resource "equinix_fabric_stream" "new_stream_3" {
		  type = "TELEMETRY_STREAM"
		  name = "%[1]s"
		  description = "%[2]s"
		  project = {
			project_id = "291639000636552"
		  }
		}

		data "equinix_fabric_stream" "data_stream" {
		  stream_id = equinix_fabric_stream.new_stream_3.id
		}
		
		data "equinix_fabric_streams" "data_streams" {
		  depends_on = [equinix_fabric_stream.new_stream_1, equinix_fabric_stream.new_stream_2, equinix_fabric_stream.new_stream_3]
		  pagination = {
			limit = 2
			offset = 1
		  }
		}
	`, name, description)
}

func TestAccFabricStreamDataSources_PFCR(t *testing.T) {
	streamName := "stream_PFCR"
	streamDescription := "stream acceptance test PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricStreamDataSourcesConfig(streamName, streamDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream.data_stream", "name", streamName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream.data_stream", "type", "TELEMETRY_STREAM"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream.data_stream", "project.project_id", "291639000636552"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream.data_stream", "description", streamDescription),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream.data_stream", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream.data_stream", "href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream.data_stream", "assets_count"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream.data_stream", "stream_subscriptions_count"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream.data_stream", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream.data_stream", "change_log.created_by"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.project.project_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.project.project_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.href"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.assets_count"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.stream_subscriptions_count"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_streams.data_streams", "data.0.change_log.created_by"),
					resource.TestCheckResourceAttr("data.equinix_fabric_streams.data_streams", "data.#", "2"),
					resource.TestCheckResourceAttr("data.equinix_fabric_streams.data_streams", "pagination.%", "5"),
					resource.TestCheckResourceAttr("data.equinix_fabric_streams.data_streams", "pagination.limit", "2"),
					resource.TestCheckResourceAttr("data.equinix_fabric_streams.data_streams", "pagination.offset", "1"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}
