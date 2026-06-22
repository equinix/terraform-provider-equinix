package routeaggregationrule_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const datasourceConfig = `
resource "equinix_fabric_route_aggregation" "test" {
  type        = "BGP_IPv4_PREFIX_AGGREGATION"
  name        = "test-aggregation"
  description = "Test Route Aggregation"
  project = {
    project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
  }
}

resource "equinix_fabric_route_aggregation_rule" "new-rar1" {
  route_aggregation_id = equinix_fabric_route_aggregation.test.id
  name                 = "route_agg_rule1_PFCR"
  description          = "route aggregation rule 1 PFCR"
  prefix               = "192.166.0.0/24"
}

resource "equinix_fabric_route_aggregation_rule" "new-rar2" {
  route_aggregation_id = equinix_fabric_route_aggregation.test.id
  name                 = "route_agg_rule2_PFCR"
  description          = "route aggregation rule 2 PFCR"
  prefix               = "192.169.0.0/24"
}

data "equinix_fabric_route_aggregation_rule" "data_rar" {
  route_aggregation_id      = equinix_fabric_route_aggregation.test.id
  route_aggregation_rule_id = equinix_fabric_route_aggregation_rule.new-rar1.id
}

data "equinix_fabric_route_aggregation_rules" "data_rars" {
  depends_on = [
    equinix_fabric_route_aggregation_rule.new-rar1,
    equinix_fabric_route_aggregation_rule.new-rar2
  ]
  route_aggregation_id = equinix_fabric_route_aggregation.test.id
  pagination = {
    limit  = 32
    offset = 0
  }

  sort = [{
    property  = "/name"
    direction = "ASC"
  }]
  outer_operator = "OR"
}

data "equinix_fabric_route_aggregation_rules" "data_rars_filtered" {
  depends_on = [
    equinix_fabric_route_aggregation_rule.new-rar1,
    equinix_fabric_route_aggregation_rule.new-rar2
  ]
  route_aggregation_id = equinix_fabric_route_aggregation.test.id
  pagination = {
    limit  = 32
    offset = 0
  }

  filter = [{
    property = "/name"
    operator = "="
    values   = ["route_agg_rule1_PFCR"]
    }
  ]

  sort = [{
    property  = "/name"
    direction = "ASC"
  }]
  outer_operator = "OR"
}


data "equinix_fabric_route_aggregation_rules" "data_rars_or" {
  depends_on = [
    equinix_fabric_route_aggregation_rule.new-rar1,
    equinix_fabric_route_aggregation_rule.new-rar2
  ]
  route_aggregation_id = equinix_fabric_route_aggregation.test.id
  pagination = {
    limit  = 32
    offset = 0
  }

  filter = [{
    property = "/name"
    operator = "="
    values   = ["route_agg_rule1_PFCR"]
    }, {
    property = "/name"
    operator = "="
    values   = ["route_agg_rule2_PFCR"]
    }
  ]

  sort = [{
    property  = "/name"
    direction = "ASC"
  }]
  outer_operator = "OR"
}

data "equinix_fabric_route_aggregation_rules" "data_rars_and" {
  depends_on = [
    equinix_fabric_route_aggregation_rule.new-rar1,
    equinix_fabric_route_aggregation_rule.new-rar2
  ]
  route_aggregation_id = equinix_fabric_route_aggregation.test.id
  pagination = {
    limit  = 32
    offset = 0
  }

  filter = [{
    property = "/name"
    operator = "="
    values   = ["route_agg_rule1_PFCR"]
    }, {
    property = "/type"
    operator = "="
    values   = ["BGP_IPv4_PREFIX_AGGREGATION_RULE"]
    }
  ]

  sort = [{
    property  = "/name"
    direction = "ASC"
  }]
  outer_operator = "AND"
}
`

func TestAccFabricRouteAggregationRuleDataSources_PNFV(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationRuleDelete,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("name"), knownvalue.StringExact("route_agg_rule1_PFCR")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("type"), knownvalue.StringExact("BGP_IPv4_PREFIX_AGGREGATION_RULE")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("description"), knownvalue.StringExact("route aggregation rule 1 PFCR")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("href"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("prefix"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rule.data_rar", tfjsonpath.New("state"), knownvalue.NotNull()),

					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation_rules.data_rars", tfjsonpath.New("pagination"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"offset": knownvalue.Int32Exact(0),
							"limit":  knownvalue.Int32Exact(32),
							"total":  knownvalue.Int32Exact(2),
						}),
					),

					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_aggregation_rules.data_rars",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":        knownvalue.NotNull(),
								"name":        knownvalue.StringExact("route_agg_rule1_PFCR"),
								"description": knownvalue.StringExact("route aggregation rule 1 PFCR"),
								"prefix":      knownvalue.StringExact("192.166.0.0/24"),
								"change_log":  knownvalue.NotNull(),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":        knownvalue.NotNull(),
								"name":        knownvalue.StringExact("route_agg_rule2_PFCR"),
								"description": knownvalue.StringExact("route aggregation rule 2 PFCR"),
								"prefix":      knownvalue.StringExact("192.169.0.0/24"),
							}),
						}),
					),

					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_aggregation_rules.data_rars_filtered",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":        knownvalue.NotNull(),
								"name":        knownvalue.StringExact("route_agg_rule1_PFCR"),
								"description": knownvalue.StringExact("route aggregation rule 1 PFCR"),
								"prefix":      knownvalue.StringExact("192.166.0.0/24"),
							}),
						}),
					),

					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_aggregation_rules.data_rars_or",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":        knownvalue.NotNull(),
								"name":        knownvalue.StringExact("route_agg_rule1_PFCR"),
								"description": knownvalue.StringExact("route aggregation rule 1 PFCR"),
								"prefix":      knownvalue.StringExact("192.166.0.0/24"),
							}),
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":        knownvalue.NotNull(),
								"name":        knownvalue.StringExact("route_agg_rule2_PFCR"),
								"description": knownvalue.StringExact("route aggregation rule 2 PFCR"),
								"prefix":      knownvalue.StringExact("192.169.0.0/24"),
							}),
						}),
					),

					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_aggregation_rules.data_rars_and",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":        knownvalue.NotNull(),
								"name":        knownvalue.StringExact("route_agg_rule1_PFCR"),
								"description": knownvalue.StringExact("route aggregation rule 1 PFCR"),
								"prefix":      knownvalue.StringExact("192.166.0.0/24"),
							}),
						}),
					),
				},
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
