package routeaggregation_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

var datasourceConfig = `
resource "equinix_fabric_route_aggregation" "new_ra_1" {
  type        = "BGP_IPv4_PREFIX_AGGREGATION"
  name        = "RA_datasource_1_PFCR"
  description = "route_agg_1_PFCR_description"
  project = {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
}

resource "equinix_fabric_route_aggregation" "new_ra_2" {
  depends_on = [equinix_fabric_route_aggregation.new_ra_1] # to ensure that 2 is created after 1
  type        = "BGP_IPv4_PREFIX_AGGREGATION"
  name        = "RA_datasource_2_PFCR"
  description = "route_agg_2_PFCR_description"
  project = {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
}

data "equinix_fabric_route_aggregation" "data_ra" {
  route_aggregation_id = equinix_fabric_route_aggregation.new_ra_1.id
}

data "equinix_fabric_route_aggregations" "data_ras" {
  depends_on = [equinix_fabric_route_aggregation.new_ra_1, equinix_fabric_route_aggregation.new_ra_2]

  filter = {
    property = "/name"
    operator = "LIKE"
    values   = ["RA_datasource_%_PFCR"]
  }
  pagination = {
    limit  = 2
    offset = 0
  }
  sort = {
    property  = "/changeLog/createdDateTime"
    direction = "DESC"
  }
}
`

func TestAccFabricRouteAggregationDataSources_PFCR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             CheckRouteAggregationDelete,
		Steps: []resource.TestStep{
			{
				Config: datasourceConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("href"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("uuid"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("name"), knownvalue.StringExact("RA_datasource_1_PFCR")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("description"), knownvalue.StringExact("route_agg_1_PFCR_description")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("type"), knownvalue.StringExact("BGP_IPv4_PREFIX_AGGREGATION")),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("change_log"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"created_by": knownvalue.NotNull(),
						}),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregation.data_ra", tfjsonpath.New("project"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
						}),
					),

					statecheck.ExpectKnownValue(
						"data.equinix_fabric_route_aggregations.data_ras",
						tfjsonpath.New("data"),
						knownvalue.ListPartial(map[int]knownvalue.Check{
							0: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":              knownvalue.NotNull(),
								"name":              knownvalue.StringExact("RA_datasource_2_PFCR"),
								"description":       knownvalue.StringExact("route_agg_2_PFCR_description"),
								"href":              knownvalue.NotNull(),
								"connections_count": knownvalue.Int32Exact(0),
								"rules_count":       knownvalue.Int32Exact(0),
								"uuid":              knownvalue.NotNull(),
								"change_log": knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"created_by": knownvalue.NotNull(),
								}),
								"project": knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
								}),
							}),
							1: knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"type":              knownvalue.NotNull(),
								"name":              knownvalue.StringExact("RA_datasource_1_PFCR"),
								"description":       knownvalue.StringExact("route_agg_1_PFCR_description"),
								"href":              knownvalue.NotNull(),
								"connections_count": knownvalue.Int32Exact(0),
								"rules_count":       knownvalue.Int32Exact(0),
								"uuid":              knownvalue.NotNull(),
								"change_log": knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"created_by": knownvalue.NotNull(),
								}),
								"project": knownvalue.ObjectPartial(map[string]knownvalue.Check{
									"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_route_aggregations.data_ras", tfjsonpath.New("pagination"),
						knownvalue.ObjectPartial(map[string]knownvalue.Check{
							"limit":  knownvalue.Int32Exact(2),
							"offset": knownvalue.Int32Exact(0),
						}),
					),
				},
			},
		},
	})

}
