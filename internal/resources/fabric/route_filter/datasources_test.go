package route_filter_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var datasourceConfig = ` 
resource "equinix_fabric_route_filter" "test" {
  name = "RF_DS_Policy_PFCR"
  project {
    project_id = "291639000636552"
  }
  type        = "BGP_IPv4_PREFIX_FILTER"
  description = "Route Filter Policy for X Purpose"
}

data "equinix_fabric_route_filter" "rf_policy" {
  uuid = equinix_fabric_route_filter.test.id
}

data "equinix_fabric_route_filters" "rf_policies" {
  filter {
    property = "/type"
    operator = "="
    values   = ["BGP_IPv4_PREFIX_FILTER"]
  }
  filter {
    property = "/state"
    operator = "="
    values   = ["PROVISIONED"]
  }
  filter {
    property = "/project/projectId"
    operator = "="
    values   = ["291639000636552"]
  }
  filter {
    property = "/name"
    operator = "="
    values   = [equinix_fabric_route_filter.test.name]
  }
  pagination {
    offset = 0
    limit  = 5
  }
  sort {
    direction = "ASC"
    property  = "/name"
  }
}
`

func TestAccFabricRouteFilterPolicy_DataSources_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_route_filter.test", map[string]knownvalue.Check{
						"id":                      knownvalue.NotNull(),
						"name":                    knownvalue.StringExact("RF_DS_Policy_PFCR"),
						"type":                    knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER"),
						"state":                   knownvalue.StringExact("PROVISIONED"),
						"not_matched_rule_action": knownvalue.StringExact("DENY"),
						"rules_count":             knownvalue.Int32Exact(0),
						"description":             knownvalue.StringExact("Route Filter Policy for X Purpose"),
					}),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_filters.rf_policies", tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("RF_DS_Policy_PFCR"),
								"type":                    knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER"),
								"state":                   knownvalue.StringExact("PROVISIONED"),
								"not_matched_rule_action": knownvalue.StringExact("DENY"),
								"rules_count":             knownvalue.Int32Exact(0),
								"description":             knownvalue.StringExact("Route Filter Policy for X Purpose"),
							}),
						})),
				},
				ExpectNonEmptyPlan: false,
			},
		},
	})

}
