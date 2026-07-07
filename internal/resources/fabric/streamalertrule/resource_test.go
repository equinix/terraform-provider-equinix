package streamalertrule_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	eqconfig "github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func CheckStreamAlertRuleDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*eqconfig.Config).NewFabricClientForTesting(ctx)

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

var resourceConfig = `
variable "aside_vlan" {
  type = number
}

variable "aside_port_uuid" {
  type = string
}

variable "zside_vlan" {
  type = number
}

variable "zside_port_uuid" {
  type = string
}

variable "uri" {
  type = string
}

variable "event_index" {
  type = string
}

variable "metric_index" {
  type = string
}


variable "alert_source" {
  type = string
}

variable "access_token" {
  type      = string
  sensitive = true
}

variable "name" {
  type = string
}

variable "description" {
  type = string
}

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
    uri  = var.uri
    settings = {
      event_index  = var.event_index
      metric_index = var.metric_index
      source       = var.alert_source
    }
    credential = {
      type         = "ACCESS_TOKEN"
      access_token = var.access_token
    }
  }
  lifecycle {
    create_before_destroy = true
  }
}
resource "equinix_fabric_connection" "test_connection" {
  name = "Test Connection PFCR"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 50
  redundancy { priority = "PRIMARY" }
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.aside_port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = var.aside_vlan
      }
    }
  }
  z_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.zside_port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = var.zside_vlan
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
  stream_id   = equinix_fabric_stream.new_stream.uuid
  name        = var.name
  type        = "METRIC_ALERT"
  description = var.description
  detection_method = {
    type               = "THRESHOLD"
    operand            = "ABOVE"
    window_size        = "PT15M"
    warning_threshold  = "35000000"
    critical_threshold = "40000000"
  }
  metric_selector = {
    "include" : [
      "equinix.fabric.connection.bandwidth_tx.usage"
    ]
  }
  resource_selector = {
    "include" : [
      "*/connections/${equinix_fabric_connection.test_connection.id}"
    ]
  }
}
data "equinix_fabric_stream_alert_rule" "by_ids" {
  stream_id     = equinix_fabric_stream.new_stream.uuid
  alert_rule_id = equinix_fabric_stream_alert_rule.alert_rule.uuid
}

data "equinix_fabric_stream_alert_rules" "all" {
  depends_on = [
    equinix_fabric_stream.new_stream,
    equinix_fabric_stream_alert_rule.alert_rule
  ]
  stream_id = equinix_fabric_stream.new_stream.uuid
  pagination = {
    limit  = 20
    offset = 0
  }
}
`

func TestAccFabricStreamAlertRule_PFCR(t *testing.T) {
	streamData := testinghelpers.GetFabricStreamTestData(t)
	uri := streamData["splunk"]["uri"]
	accessToken := streamData["splunk"]["accessToken"]
	eventIndex := streamData["splunk"]["event_index"]
	metricIndex := streamData["splunk"]["metric_index"]
	source := streamData["splunk"]["source"]

	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortUUID, zSidePortUUID string
	if len(ports) > 0 {
		aSidePortUUID = ports["pfcr"]["dot1q"][0].GetUuid()
		zSidePortUUID = ports["pfcr"]["dot1q"][1].GetUuid()
	}

	asideVlan, err := testinghelpers.RandomVlan(aSidePortUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	zsideVlan, err := testinghelpers.RandomVlan(zSidePortUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckStreamAlertRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				ConfigVariables: config.Variables{
					"name":            config.StringVariable("alert_rule_PFCR"),
					"description":     config.StringVariable("stream alert rule acceptance test PFCR"),
					"aside_vlan":      config.IntegerVariable(asideVlan),
					"aside_port_uuid": config.StringVariable(aSidePortUUID),
					"zside_vlan":      config.IntegerVariable(zsideVlan),
					"zside_port_uuid": config.StringVariable(zSidePortUUID),
					"uri":             config.StringVariable(uri),
					"event_index":     config.StringVariable(eventIndex),
					"metric_index":    config.StringVariable(metricIndex),
					"alert_source":    config.StringVariable(source),
					"access_token":    config.StringVariable(accessToken),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_stream_alert_rule.alert_rule",
						map[string]knownvalue.Check{
							"href":    knownvalue.NotNull(),
							"uuid":    knownvalue.NotNull(),
							"state":   knownvalue.StringExact("ACTIVE"),
							"enabled": knownvalue.Bool(true),
							"detection_method": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":               knownvalue.StringExact("THRESHOLD"),
								"operand":            knownvalue.StringExact("ABOVE"),
								"window_size":        knownvalue.StringExact("PT15M"),
								"warning_threshold":  knownvalue.StringExact("35000000"),
								"critical_threshold": knownvalue.StringExact("40000000"),
							}),
							"change_log": knownvalue.NotNull(),
						}),

					testinghelpers.ExpectKnownAttributes("data.equinix_fabric_stream_alert_rule.by_ids",
						map[string]knownvalue.Check{
							"type":        knownvalue.StringExact("METRIC_ALERT"),
							"name":        knownvalue.StringExact("alert_rule_PFCR"),
							"description": knownvalue.StringExact("stream alert rule acceptance test PFCR"),
							"enabled":     knownvalue.Bool(true),
							"resource_selector": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"include": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringRegexp(regexp.MustCompile(`\*/connections/.+`))}),
							}),
							"href":  knownvalue.NotNull(),
							"uuid":  knownvalue.NotNull(),
							"state": knownvalue.StringExact("ACTIVE"),
							"detection_method": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":               knownvalue.StringExact("THRESHOLD"),
								"operand":            knownvalue.StringExact("ABOVE"),
								"window_size":        knownvalue.StringExact("PT15M"),
								"warning_threshold":  knownvalue.StringExact("35000000"),
								"critical_threshold": knownvalue.StringExact("40000000"),
							}),
							"change_log": knownvalue.NotNull(),
						}),

					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_stream_alert_rules.all",
						tfjsonpath.New("data").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"type":              knownvalue.NotNull(),
							"name":              knownvalue.NotNull(),
							"description":       knownvalue.NotNull(),
							"enabled":           knownvalue.NotNull(),
							"resource_selector": knownvalue.NotNull(),
							"href":              knownvalue.NotNull(),
							"uuid":              knownvalue.NotNull(),
							"state":             knownvalue.NotNull(),
							"detection_method": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":               knownvalue.NotNull(),
								"operand":            knownvalue.NotNull(),
								"window_size":        knownvalue.NotNull(),
								"warning_threshold":  knownvalue.NotNull(),
								"critical_threshold": knownvalue.NotNull(),
							}),
							"change_log": knownvalue.NotNull(),
						}),
				},
				ExpectNonEmptyPlan: true,
			},
			{
				Config: resourceConfig,
				ConfigVariables: config.Variables{
					"name":            config.StringVariable("up_alert_rule_PFCR"),
					"description":     config.StringVariable("updated stream alert rule acceptance test PFCR"),
					"aside_vlan":      config.IntegerVariable(asideVlan),
					"aside_port_uuid": config.StringVariable(aSidePortUUID),
					"zside_vlan":      config.IntegerVariable(zsideVlan),
					"zside_port_uuid": config.StringVariable(zSidePortUUID),
					"uri":             config.StringVariable(uri),
					"event_index":     config.StringVariable(eventIndex),
					"metric_index":    config.StringVariable(metricIndex),
					"alert_source":    config.StringVariable(source),
					"access_token":    config.StringVariable(accessToken),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_stream_alert_rule.alert_rule",
						map[string]knownvalue.Check{
							"href":    knownvalue.NotNull(),
							"uuid":    knownvalue.NotNull(),
							"state":   knownvalue.StringExact("ACTIVE"),
							"enabled": knownvalue.Bool(true),
							"detection_method": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":               knownvalue.StringExact("THRESHOLD"),
								"operand":            knownvalue.StringExact("ABOVE"),
								"window_size":        knownvalue.StringExact("PT15M"),
								"warning_threshold":  knownvalue.StringExact("35000000"),
								"critical_threshold": knownvalue.StringExact("40000000"),
							}),
							"change_log": knownvalue.NotNull(),
						}),

					testinghelpers.ExpectKnownAttributes("data.equinix_fabric_stream_alert_rule.by_ids",
						map[string]knownvalue.Check{
							"type":        knownvalue.StringExact("METRIC_ALERT"),
							"name":        knownvalue.StringExact("up_alert_rule_PFCR"),
							"description": knownvalue.StringExact("updated stream alert rule acceptance test PFCR"),
							"enabled":     knownvalue.Bool(true),
							"resource_selector": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"include": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringRegexp(regexp.MustCompile(`\*/connections/.+`))}),
							}),
							"href":  knownvalue.NotNull(),
							"uuid":  knownvalue.NotNull(),
							"state": knownvalue.StringExact("ACTIVE"),
							"detection_method": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":               knownvalue.StringExact("THRESHOLD"),
								"operand":            knownvalue.StringExact("ABOVE"),
								"window_size":        knownvalue.StringExact("PT15M"),
								"warning_threshold":  knownvalue.StringExact("35000000"),
								"critical_threshold": knownvalue.StringExact("40000000"),
							}),
							"change_log": knownvalue.NotNull(),
						}),

					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_stream_alert_rules.all",
						tfjsonpath.New("data").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"type":              knownvalue.NotNull(),
							"name":              knownvalue.NotNull(),
							"description":       knownvalue.NotNull(),
							"enabled":           knownvalue.NotNull(),
							"resource_selector": knownvalue.NotNull(),
							"href":              knownvalue.NotNull(),
							"uuid":              knownvalue.NotNull(),
							"state":             knownvalue.NotNull(),
							"detection_method": knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":               knownvalue.NotNull(),
								"operand":            knownvalue.NotNull(),
								"window_size":        knownvalue.NotNull(),
								"warning_threshold":  knownvalue.NotNull(),
								"critical_threshold": knownvalue.NotNull(),
							}),
							"change_log": knownvalue.NotNull(),
						}),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
