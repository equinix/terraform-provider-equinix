package streamsubscription_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"

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
  type        = "TELEMETRY_STREAM"
  name        = "Subscription_Test_PFCR"
  description = "Testing stream subscriptions resource"
  project = {
    project_id = "291639000636552"
  }
}

resource "equinix_fabric_stream" "new_stream2" {
  type        = "TELEMETRY_STREAM"
  name        = "Subscription_Test2_PFCR"
  description = "Testing stream subscriptions resource limits"
  project = {
    project_id = "291639000636552"
  }
}

resource "equinix_fabric_stream" "new_stream3" {
  type        = "TELEMETRY_STREAM"
  name        = "Subscription_Test3_PFCR"
  description = "Testing stream subscriptions resource limits"
  project = {
    project_id = "291639000636552"
  }
}

resource "equinix_fabric_stream_subscription" "slack" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "Slack_PFCR"
  description = "Stream Subscription Slack TF Testing"
  stream_id   = equinix_fabric_stream.new_stream.id
  enabled     = false
  sink = {
    type = "SLACK"
    uri  = "%s"
  }
}

resource "equinix_fabric_stream_subscription" "pager_duty" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "PagerDuty_PFCR"
  description = "Stream Subscription PagerDuty TF Testing"
  stream_id   = equinix_fabric_stream.new_stream.id
  enabled     = false
  sink = {
    type = "PAGERDUTY"
    host = "%s"
    settings = {
      change_uri = "%s"
      alert_uri  = "%s"
    }
    credential = {
      type            = "INTEGRATION_KEY"
      integration_key = "%s"
    }
  }
}

resource "equinix_fabric_stream_subscription" "datadog" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "DataDog_PFCR"
  description = "Stream Subscription DataDog TF Testing"
  stream_id   = equinix_fabric_stream.new_stream2.id
  enabled     = false
  sink = {
    type = "DATADOG"
    host = "%s"
    settings = {
      source          = "Equinix"
      application_key = "%s"
      event_uri       = "%s"
      metric_uri      = "%s"
    }
    credential = {
      type    = "API_KEY"
      api_key = "%s"
    }
  }
}

resource "equinix_fabric_stream_subscription" "msteams" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "MSTeams_PFCR"
  description = "Stream Subscription Microsoft Teams TF Testing"
  stream_id   = equinix_fabric_stream.new_stream2.id
  enabled     = false
  sink = {
    type = "TEAMS"
    uri  = "%s"
  }
}

resource "equinix_fabric_stream_subscription" "webhook" {
  type        = "STREAM_SUBSCRIPTION"
  name        = "Webhook_PFCR"
  description = "Stream Subscription Webhook TF Testing"
  stream_id   = equinix_fabric_stream.new_stream2.id
  enabled     = false
  sink = {
    type = "WEBHOOK"
    settings = {
      format     = "CLOUDEVENT"
      event_uri  = "%s"
      metric_uri = "%s"
    }
  }
}

data "equinix_fabric_stream_subscription" "by_ids" {
  stream_id       = equinix_fabric_stream.new_stream2.id
  subscription_id = equinix_fabric_stream_subscription.webhook.id
}

data "equinix_fabric_stream_subscriptions" "all" {
  depends_on = [
    equinix_fabric_stream_subscription.slack,
    equinix_fabric_stream_subscription.pager_duty,
    equinix_fabric_stream_subscription.datadog,
    equinix_fabric_stream_subscription.msteams,
    equinix_fabric_stream_subscription.webhook,
  ]
  stream_id = equinix_fabric_stream.new_stream.id
  pagination = {
    limit  = 20
    offset = 0
  }
}
	`,
		streamTestData["slack"]["uri"],
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
		streamTestData["webhook"]["event_uri"],
		streamTestData["webhook"]["metric_uri"],
	)
}

func TestAccFabricStreamSubscription_PFCR(t *testing.T) {
	streamTestData := testinghelpers.GetFabricStreamTestData(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamSubscriptionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricStreamSubscriptionConfig(streamTestData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_subscription.by_ids", "name", "Webhook_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_subscription.by_ids", "type", "STREAM_SUBSCRIPTION"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_subscription.by_ids", "description", "Stream Subscription Webhook TF Testing"),
					resource.TestCheckResourceAttr("data.equinix_fabric_stream_subscription.by_ids", "sink.type", "WEBHOOK"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscription.by_ids", "stream_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscription.by_ids", "sink.settings.event_uri"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscription.by_ids", "sink.settings.metric_uri"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscription.by_ids", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscriptions.all", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscriptions.all", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscriptions.all", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscriptions.all", "data.0.sink.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_subscriptions.all", "data.0.uuid"),

					resource.TestCheckResourceAttr("equinix_fabric_stream_subscription.webhook", "sink.type", "WEBHOOK"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_subscription.webhook", "uuid"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}
