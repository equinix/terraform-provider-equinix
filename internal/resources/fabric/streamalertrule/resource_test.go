package streamalertrule_test

import (
	"context"
	"fmt"
	"testing"

	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func CheckStreamAlertRuleDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_stream_alert_rule" {
			continue
		}

		streamID := rs.Primary.Attributes["stream_id"]

		if streamAlertRule, _, err := client.StreamAlertRulesApi.GetStreamAlertRuleByUuid(ctx, streamID, rs.Primary.ID).Execute(); err == nil &&
			streamAlertRule.GetState() == fabricv4.STREAMALERTRULESTATE_ACTIVE {
			return fmt.Errorf("fabric stream alert rule %s still exists and is %s",
				rs.Primary.ID, string(fabricv4.STREAMALERTRULESTATE_ACTIVE))
		}
	}
	return nil
}

func testAccFabricStreamAlertRuleConfig(uri, event_index, metric_index, source, access_token, aSidePortUUID, zSidePortUUID, alertRuleName, alertRuleDescription string) string {
	return fmt.Sprintf(`
        resource "equinix_fabric_stream" "new_stream" {
		  type = "TELEMETRY_STREAM"
		  name = "Stream_Test_PFCR"

		  description = "Stream Description"
		  project = {
			project_id = "291639000636552"
		  }
		}

        resource "equinix_fabric_stream_subscription" "SPLUNK" {
			  depends_on = [
				equinix_fabric_stream.new_stream
			  ]
			
			  type        = "STREAM_SUBSCRIPTION"
			  name        = "Stream_Sub_PFCR"
			  description = "Stream Subscription for Splunk PFCR"
			  stream_id   = equinix_fabric_stream.new_stream.uuid
			  enabled     = true
			  sink = {
				type = "SPLUNK_HEC"
				uri  = "%s"
				settings = {
				  event_index  = "%s"
				  metric_index = "%s"
				  source       = "%s"
				}
				credential = {
				  type         = "ACCESS_TOKEN"
				  access_token = "%s"
				}
			  }
			}
        resource "equinix_fabric_connection" "test_connection" {
			name = "Test Connection PFCR"
			type = "EVPL_VC"
			notifications{
				type="ALL" 
				emails=["example@equinix.com"]
			} 
			bandwidth = 50
			redundancy {priority= "PRIMARY"}
			order {
				purchase_order_number= "1-323292"
			}
			a_side {
				access_point {
					type= "COLO"
					port {
						uuid= "%s"
					}
					link_protocol {
						type= "DOT1Q"
						vlan_tag= "1232"
					}
				}
			}
			z_side {
				access_point {
				  type= "COLO"
				  port {
					uuid= "%s"
				  }
				  link_protocol {
					type= "DOT1Q"
					vlan_tag= "1278"
				  }
				}
		}
}

    resource "equinix_fabric_stream_attachment" "asset" {
		    depends_on = [
				equinix_fabric_stream.new_stream,
				equinix_fabric_connection.test_connection 
		    ]
		    asset_id  = equinix_fabric_connection.test_connection.id
		    asset     = "connections"
		    stream_id = equinix_fabric_stream.new_stream.uuid
    }

    resource "equinix_fabric_stream_alert_rule" "alert_rule" {
		  depends_on = [
			equinix_fabric_stream_attachment.asset
		  ]
		  stream_id          = equinix_fabric_stream.new_stream.uuid
		  name               = "%s"
		  type               = "METRIC_ALERT"
		  description        = "%s"
		  operand            = "ABOVE"
		  window_size        = "PT15M"
		  warning_threshold  = "35000000"
		  critical_threshold = "40000000"
		  metric_name        = "equinix.fabric.connection.bandwidth_tx.usage"
		  resource_selector   = {
			"include" : [
			  "*/connections/${equinix_fabric_connection.test_connection.id}"
			]
		  }
		}
        data "equinix_fabric_stream_alert_rule" "by_ids" {
		  depends_on = [
			equinix_fabric_stream.new_stream,
			equinix_fabric_stream_alert_rule.alert_rule
		  ]
		  stream_id = equinix_fabric_stream.new_stream.uuid
		  alert_rule_id = equinix_fabric_stream_alert_rule.alert_rule.uuid
		}
        data "equinix_fabric_stream_alert_rules" "all" {
			depends_on = [
				 equinix_fabric_stream.new_stream,
				 equinix_fabric_stream_alert_rule.alert_rule
			   ]
			   stream_id = equinix_fabric_stream.new_stream.uuid
			   pagination = {
				 limit = 20
				 offset = 0
			   }
			 }
	`, uri, event_index, metric_index, source, access_token, aSidePortUUID, zSidePortUUID, alertRuleName, alertRuleDescription)
}

func TestAccFabricStreamAlertRule_PFCR(t *testing.T) {
	streamData := testinghelpers.GetFabricStreamTestData(t)
	uri := streamData["splunk"]["uri"]
	accessToken := streamData["splunk"]["accessToken"]
	eventIndex := streamData["splunk"]["event_index"]
	metricIndex := streamData["splunk"]["metric_index"]
	source := streamData["splunk"]["metric_index"]

	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortUUID, zSidePortUUID string
	if len(ports) > 0 {
		aSidePortUUID = ports["pfcr"]["dot1q"][0].GetUuid()
		zSidePortUUID = ports["pfcr"]["dot1q"][1].GetUuid()
	}
	alertRuleName, updatedAlertRuleName := "alert_rule_PFCR", "up_alert_rule_PFCR"
	alertRuleDescription, updatedAlertRuleDescription := "stream alert rule acceptance test PFCR", "updated stream alert rule acceptance test PFCR"

	//alertRuleName := "alert_rule_PFCR"
	//alertRuleDescription := "stream alert rule acceptance test PFCR"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamAlertRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricStreamAlertRuleConfig(uri, eventIndex, metricIndex, source, accessToken, aSidePortUUID, zSidePortUUID, alertRuleName, alertRuleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_alert_rule.alert_rule", "name", alertRuleName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_alert_rule.alert_rule", "type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_alert_rule.alert_rule", "description", alertRuleDescription),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "metric_name", "equinix.fabric.connection.bandwidth_tx.usage"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "operand", "ABOVE"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "enabled", "true"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "window_size", "PT15M"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "warning_threshold", "35000000"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "critical_threshold", "40000000"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_alert_rule.alert_rule", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_alert_rule.alert_rule", "href"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_alert_rule.by_ids", "name", alertRuleName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_alert_rule.by_ids", "type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_alert_rule.by_ids", "description", alertRuleDescription),
					resource.TestCheckResourceAttr("data.equinix_fabric_stream_alert_rule.by_ids", "enabled", "true"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "window_size"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "critical_threshold"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "warning_threshold"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "metric_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "operand"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "href"),

					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.metric_name"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricStreamAlertRuleConfig(uri, eventIndex, metricIndex, source, accessToken, aSidePortUUID, zSidePortUUID, updatedAlertRuleName, updatedAlertRuleDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_alert_rule.alert_rule", "name", updatedAlertRuleName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_alert_rule.alert_rule", "type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_stream_alert_rule.alert_rule", "description", updatedAlertRuleDescription),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "metric_name", "equinix.fabric.connection.bandwidth_tx.usage"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "operand", "ABOVE"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "enabled", "true"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "window_size", "PT15M"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "warning_threshold", "35000000"),
					resource.TestCheckResourceAttr("equinix_fabric_stream_alert_rule.alert_rule", "critical_threshold", "40000000"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_alert_rule.alert_rule", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_stream_alert_rule.alert_rule", "href"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_alert_rule.by_ids", "name", updatedAlertRuleName),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_alert_rule.by_ids", "type", "METRIC_ALERT"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_stream_alert_rule.by_ids", "description", updatedAlertRuleDescription),
					resource.TestCheckResourceAttr("data.equinix_fabric_stream_alert_rule.by_ids", "enabled", "true"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "window_size"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "critical_threshold"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "warning_threshold"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "metric_name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "operand"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rule.by_ids", "href"),

					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.name"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.type"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.description"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_stream_alert_rules.all", "data.0.metric_name"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
