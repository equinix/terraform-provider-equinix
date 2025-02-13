package streamsubscription_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckStreamSubscriptionDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_stream_subscription" {
			continue
		}

		streamID := rs.Primary.Attributes["stream_id"]

		if subscription, _, err := client.StreamSubscriptionsApi.GetStreamSubscriptionByUuid(ctx, streamID, rs.Primary.ID).Execute(); err == nil &&
			subscription.GetState() == fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONED {
			return fmt.Errorf("fabric stream %s still exists and is %s",
				rs.Primary.ID, string(fabricv4.STREAMSUBSCRIPTIONSTATE_PROVISIONED))
		}
	}
	return nil
}

func testAccFabricStreamSubscriptionConfig(streamTestData map[string]map[string]string) string {
	return fmt.Sprintf(`
		resource "equinix_fabric_stream" "new_stream" {
		  type = "TELEMETRY_STREAM"
		  name = "Subscription_Test_PFCR"
		  description = "Testing stream subscriptions resource"
		  project = {
			project_id = "291639000636552"
		  }
		}

		resource "equinix_fabric_stream_subscription" "splunk" {
		  type = "STREAM_SUBSCRIPTION"
		  name = "Splunk_PFCR"
		  description = "Stream Subscription Splunk TF Testing"
		  stream_id = equinix_fabric_stream.new_stream.id
		  enabled = false
		  sink = {
			type = "SPLUNK_HEC"
			uri = "%s"
			settings = {
			  event_index  = "%s"
			  metric_index = "%s"
			  source = "%s"
			}
			credential = {
			  type = "ACCESS_TOKEN"
			  access_token = "%s"
			}
		  }
		}

		resource "equinix_fabric_stream_subscription" "slack" {
		  type = "STREAM_SUBSCRIPTION"
		  name = "Slack_PFCR"
		  description = "Stream Subscription Slack TF Testing"
		  stream_id = equinix_fabric_stream.new_stream.id
		  enabled = false
		  sink = {
			type = "SLACK"
			uri = "https://hooks.slack.com/services/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"
		  }
		}

		resource "equinix_fabric_stream_subscription" "pager_duty" {
		  type = "STREAM_SUBSCRIPTION"
		  name = "PagerDuty_PFCR"
		  description = "Stream Subscription PagerDuty TF Testing"
		  stream_id = equinix_fabric_stream.new_stream.id
		  enabled = false
		  sink = {
			type = "PAGERDUTY"
			host = "%s"
			settings = {
				transform_alerts = true
			    change_uri = "%s"
			    alert_uri = "%s"
			}
			credential = {
			  type = "INTEGRATION_KEY"
			  integration_key = "%s"
			}
		  }
		}

		resource "equinix_fabric_stream_subscription" "datadog" {
		  type = "STREAM_SUBSCRIPTION"
		  name = "DataDog_PFCR"
		  description = "Stream Subscription DataDog TF Testing"
		  stream_id = equinix_fabric_stream.new_stream.id
		  enabled = false
		  sink = {
			type = "DATADOG"
			host = "%s"
			settings = {
				source = "Equinix"
				application_key = "%s"
			    event_uri = "%s"
			    metric_uri = "%s"
			}
			credential = {
			  type = "API_KEY"
			  api_key = "%s"
			}
		  }
		}

		resource "equinix_fabric_stream_subscription" "msteams" {
		  type = "STREAM_SUBSCRIPTION"
		  name = "MSTeams_PFCR"
		  description = "Stream Subscription Microsoft Teams TF Testing"
		  stream_id = equinix_fabric_stream.new_stream.id
		  enabled = false
		  sink = {
			type = "TEAMS"
			uri = "%s"
		  }
		}
	`,
		streamTestData["splunk"]["uri"],
		streamTestData["splunk"]["event_index"],
		streamTestData["splunk"]["metric_index"],
		streamTestData["splunk"]["source"],
		streamTestData["splunk"]["accessToken"],
		streamTestData["pagerDuty"]["host"],
		streamTestData["pagerDuty"]["change_uri"],
		streamTestData["pagerDuty"]["alert_uri"],
		streamTestData["pagerDuty"]["integrationKey"],
		streamTestData["datadog"]["host"],
		streamTestData["datadog"]["applicationKey"],
		streamTestData["datadog"]["event_uri"],
		streamTestData["datadog"]["metric_uri"],
		streamTestData["datadog"]["APIKey"],
		streamTestData["msteams"]["uri"],
	)
}

func TestAccFabricStreamSubscription_PFCR(t *testing.T) {
	streamTestData := testing_helpers.GetFabricStreamTestData(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamSubscriptionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricStreamSubscriptionConfig(streamTestData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "name", "Subscription_Test_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "type", "TELEMETRY_STREAM"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "project.project_id", "291639000636552"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream.new_stream", "description", "Testing stream subscriptions resource"),
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
