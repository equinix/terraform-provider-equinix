package streamattachment_test

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckStreamAttachmentDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_stream_attachment" {
			continue
		}

		assetID := rs.Primary.Attributes["asset_id"]
		asset := rs.Primary.Attributes["asset"]
		streamID := rs.Primary.Attributes["stream_id"]

		if attachment, deleteResp, err := client.StreamsApi.GetStreamAssetByUuid(ctx, assetID, fabricv4.Asset(asset), streamID).Execute(); err == nil {
			if deleteResp == nil ||
				!slices.Contains([]int{http.StatusBadRequest, http.StatusForbidden, http.StatusNotFound}, deleteResp.StatusCode) ||
				attachment.GetAttachmentStatus() == fabricv4.STREAMASSETATTACHMENTSTATUS_ATTACHED {
				return fmt.Errorf("fabric stream %s still exists and is %s",
					rs.Primary.ID, string(fabricv4.STREAMASSETATTACHMENTSTATUS_ATTACHED))
			}
		}
	}
	return nil
}

func testAccFabricStreamAttachmentConfig() string {
	return `
		resource "equinix_fabric_stream" "new_stream" {
		  type = "TELEMETRY_STREAM"
		  name = "Attachment_Test_PFCR"
		  description = "Testing Stream Attachment resource"
		  project = {
			project_id = "291639000636552"
		  }
		}

		resource "equinix_fabric_cloud_router" "test"{
			type = "XF_ROUTER"
			name = "STREAM_TEST_PFCR"
			location{
				metro_code  = "SV"
			}
			package{
				code = "STANDARD"
			}
			order{
				purchase_order_number = "1-234567"
			}
			notifications{
				type = "ALL"
				emails = [
					"test@equinix.com",
					"test1@equinix.com"
				]
			}
			project{
				project_id = "291639000636552"
			}
			account {
				account_number = 201257
			}
		}

		resource "equinix_fabric_stream_attachment" "router" {
			asset_id = equinix_fabric_cloud_router.test.id
			asset = "routers"
			stream_id = equinix_fabric_stream.new_stream.id
		}

		data "equinix_fabric_stream_attachment" "by_ids" {
			depends_on = [equinix_fabric_stream_attachment.router]
			asset_id = equinix_fabric_cloud_router.test.id
			asset = "routers"
			stream_id = equinix_fabric_stream.new_stream.id
		}

		data "equinix_fabric_stream_attachments" "all" {
			depends_on = [equinix_fabric_stream_attachment.router]
			pagination = {
				limit = 100
				offset = 0
			}
			filters = [{
				property = "/streamUuid"
				operator = "="
				values = [equinix_fabric_stream.new_stream.id]
			}]
			sort = [{
				direction = "DESC"
				property = "/uuid"
			}]
		}
	`
}

func TestAccFabricStreamAttachment_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamAttachmentDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricStreamAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_attachment.router", "asset", "routers"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_attachment.router", "asset_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_attachment.router", "stream_id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_attachment.router", "metrics_enabled"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_attachment.router", "type"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_attachment.router", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_attachment.router", "href"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_attachment.router", "attachment_status", "ATTACHED"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_attachment.by_ids", "asset", "routers"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachment.by_ids", "asset_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachment.by_ids", "stream_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachment.by_ids", "metrics_enabled"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachment.by_ids", "type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachment.by_ids", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachment.by_ids", "href"),
					resource.TestCheckResourceAttr("data.equinix_fabric_stream_attachment.by_ids", "attachment_status", "ATTACHED"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachments.all", "data.0.metrics_enabled"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachments.all", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachments.all", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_attachments.all", "data.0.href"),
					resource.TestCheckResourceAttr("data.equinix_fabric_stream_attachments.all", "data.0.attachment_status", "ATTACHED"),
				),
			},
		},
	})

}
