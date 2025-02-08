package stream_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckStreamDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_stream" {
			continue
		}

		if stream, _, err := client.StreamsApi.GetStreamByUuid(ctx, rs.Primary.ID).Execute(); err == nil &&
			stream.GetState() == string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONED) {
			return fmt.Errorf("fabric stream %s still exists and is %s",
				rs.Primary.ID, string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONED))
		}
	}
	return nil
}

func testAccFabricStreamConfig(name, description string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_stream" "new_stream" {
		  type = "TELEMETRY_STREAM"
		  name = "%s"
		  description = "%s"
		  project = {
			project_id = "291639000636552"
		  }
		}
	`, name, description)
}

func TestAccFabricStream_PFCR(t *testing.T) {
	streamName, updatedStreamName := "stream_PFCR", "stream_up_PFCR"
	streamDescription, updatedStreamDescription := "stream acceptance test PFCR", "updated stream acceptance test PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricStreamConfig(streamName, streamDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "name", streamName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "type", "TELEMETRY_STREAM"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "project.project_id", "291639000636552"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "description", streamDescription),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "assets_count"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "stream_subscriptions_count"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "change_log.created_by"),
				),
			},
			{
				Config: testAccFabricStreamConfig(updatedStreamName, updatedStreamDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "name", updatedStreamName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "project.project_id", "291639000636552"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "description", updatedStreamDescription),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "href"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "assets_count"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "stream_subscriptions_count"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream.new_stream", "change_log.created_by"),
				),
			},
		},
	})

}
