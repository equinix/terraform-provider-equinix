package route_filter_rule_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const config = `
resource "equinix_fabric_route_filter" "route_filter" {
  name = "rf_ds_test_PFCR"
  project {
    project_id = "291639000636552"
  }
  type        = "BGP_IPv4_PREFIX_FILTER"
  description = "Route Filter Policy for X Purpose"
}

resource "equinix_fabric_route_filter_rule" "rule1" {
  route_filter_id = equinix_fabric_route_filter.route_filter.id
  name            = "RF_DS_Rule1_PFCR"
  prefix          = "192.168.0.0/24"
  prefix_match    = "exact"
  description     = "Route Filter Rule 1"
}

resource "equinix_fabric_route_filter_rule" "rule2" {
  route_filter_id = equinix_fabric_route_filter.route_filter.id
  name            = "RF_DS_Rule2_PFCR"
  prefix          = "192.168.0.0/25"
  prefix_match    = "orlonger"
  description     = "Route Filter Rule 2"
}

data "equinix_fabric_route_filter_rule" "rf_rule" {
  route_filter_id = equinix_fabric_route_filter.route_filter.id
  uuid            = equinix_fabric_route_filter_rule.rule1.id
}

data "equinix_fabric_route_filter_rules" "rf_rules" {
  depends_on = [
    equinix_fabric_route_filter_rule.rule1,
    equinix_fabric_route_filter_rule.rule2,
  ]
  route_filter_id = equinix_fabric_route_filter.route_filter.id

  sort {
    property  = "/name"
    direction = "ASC"
  }

  outer_operator = "AND"
}

data "equinix_fabric_route_filter_rules" "rf_rules_filtered" {
  depends_on = [
    equinix_fabric_route_filter_rule.rule1,
    equinix_fabric_route_filter_rule.rule2,
  ]

  route_filter_id = equinix_fabric_route_filter.route_filter.id

  filter {
    property = "/name"
    operator = "="
    values   = ["RF_DS_Rule2_PFCR"]
  }

  outer_operator = "AND"
}

data "equinix_fabric_route_filter_rules" "rf_rules_or" {
  depends_on = [
    equinix_fabric_route_filter_rule.rule1,
    equinix_fabric_route_filter_rule.rule2,
  ]
  route_filter_id = equinix_fabric_route_filter.route_filter.id

  filter {
    property = "/name"
    operator = "="
    values   = ["RF_DS_Rule1_PFCR"]
  }

  filter {
    property = "/name"
    operator = "="
    values   = ["RF_DS_Rule2_PFCR"]
  }

  sort {
    property  = "/name"
    direction = "ASC"
  }

  outer_operator = "OR"
}

data "equinix_fabric_route_filter_rules" "rf_rules_and" {
  depends_on = [
    equinix_fabric_route_filter_rule.rule1,
    equinix_fabric_route_filter_rule.rule2,
  ]

  route_filter_id = equinix_fabric_route_filter.route_filter.id

  filter {
    property = "/type"
    operator = "="
    values   = ["BGP_IPv4_PREFIX_FILTER_RULE"]
  }
  filter {
    property = "/name"
    operator = "="
    values   = ["RF_DS_Rule1_PFCR"]
  }

  outer_operator = "AND"
}
`

func TestAccFabricRouteFilterRule_DataSources_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: config,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rule.rf_rule", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rule.rf_rule", tfjsonpath.New("route_filter_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rule.rf_rule", tfjsonpath.New("name"), knownvalue.StringExact("RF_DS_Rule1_PFCR")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rule.rf_rule", tfjsonpath.New("description"), knownvalue.StringExact("Route Filter Rule 1")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rule.rf_rule", tfjsonpath.New("prefix"), knownvalue.StringExact("192.168.0.0/24")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rule.rf_rule", tfjsonpath.New("prefix_match"), knownvalue.StringExact("exact")),

					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules", tfjsonpath.New("route_filter_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_filter_rules.rf_rules",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":         knownvalue.NotNull(),
								"name":         knownvalue.StringExact("RF_DS_Rule1_PFCR"),
								"description":  knownvalue.StringExact("Route Filter Rule 1"),
								"prefix":       knownvalue.StringExact("192.168.0.0/24"),
								"prefix_match": knownvalue.StringExact("exact"),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":         knownvalue.NotNull(),
								"name":         knownvalue.StringExact("RF_DS_Rule2_PFCR"),
								"description":  knownvalue.StringExact("Route Filter Rule 2"),
								"prefix":       knownvalue.StringExact("192.168.0.0/25"),
								"prefix_match": knownvalue.StringExact("orlonger"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules_filtered", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules_filtered", tfjsonpath.New("route_filter_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_filter_rules.rf_rules_filtered",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":         knownvalue.NotNull(),
								"name":         knownvalue.StringExact("RF_DS_Rule2_PFCR"),
								"description":  knownvalue.StringExact("Route Filter Rule 2"),
								"prefix":       knownvalue.StringExact("192.168.0.0/25"),
								"prefix_match": knownvalue.StringExact("orlonger"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules_or", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules_or", tfjsonpath.New("route_filter_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_filter_rules.rf_rules_or",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":         knownvalue.NotNull(),
								"name":         knownvalue.StringExact("RF_DS_Rule1_PFCR"),
								"description":  knownvalue.StringExact("Route Filter Rule 1"),
								"prefix":       knownvalue.StringExact("192.168.0.0/24"),
								"prefix_match": knownvalue.StringExact("exact"),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":         knownvalue.NotNull(),
								"name":         knownvalue.StringExact("RF_DS_Rule2_PFCR"),
								"description":  knownvalue.StringExact("Route Filter Rule 2"),
								"prefix":       knownvalue.StringExact("192.168.0.0/25"),
								"prefix_match": knownvalue.StringExact("orlonger"),
							}),
						}),
					),

					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules_and", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filter_rules.rf_rules_and", tfjsonpath.New("route_filter_id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_filter_rules.rf_rules_and",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":         knownvalue.NotNull(),
								"name":         knownvalue.StringExact("RF_DS_Rule1_PFCR"),
								"description":  knownvalue.StringExact("Route Filter Rule 1"),
								"prefix":       knownvalue.StringExact("192.168.0.0/24"),
								"prefix_match": knownvalue.StringExact("exact"),
							}),
						}),
					),
				},
				ExpectNonEmptyPlan: false,
			},
		},
	})

}
