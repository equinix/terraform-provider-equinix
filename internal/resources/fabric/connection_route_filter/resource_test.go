package connection_route_filter_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection_route_filter"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricConnectionRouteFilter_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUUID string
	if len(ports) > 0 {
		portUUID = ports["pfcr"]["dot1q"][0].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUUID)

	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	newConnectionId := statecheck.CompareValue(compare.ValuesSame())
	newRouteFilterId := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionRouteFilterDelete,
		Steps: []resource.TestStep{
			{
				Config: connectionRouteFilterConfig,
				ConfigVariables: config.Variables{
					"port_uuid": config.StringVariable(portUUID),
					"vlan_tag":  config.IntegerVariable(targetVlan),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_connection_route_filter.test", map[string]knownvalue.Check{
						"id":                knownvalue.NotNull(),
						"direction":         knownvalue.StringExact("INBOUND"),
						"type":              knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER"),
						"attachment_status": knownvalue.StringExact(string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION)),
					}),
					newConnectionId.AddStateValue("equinix_fabric_connection_route_filter.test", tfjsonpath.New("connection_id")),
					newRouteFilterId.AddStateValue("equinix_fabric_connection_route_filter.test", tfjsonpath.New("route_filter_id")),

					testinghelpers.ExpectKnownAttributes("data.equinix_fabric_connection_route_filter.test", map[string]knownvalue.Check{
						"id":                knownvalue.NotNull(),
						"direction":         knownvalue.StringExact("INBOUND"),
						"type":              knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER"),
						"attachment_status": knownvalue.StringExact(string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION)),
					}),
					newConnectionId.AddStateValue("data.equinix_fabric_connection_route_filter.test", tfjsonpath.New("connection_id")),
					newRouteFilterId.AddStateValue("data.equinix_fabric_connection_route_filter.test", tfjsonpath.New("route_filter_id")),

					statecheck.ExpectKnownValue("data.equinix_fabric_connection_route_filters.test", tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":              knownvalue.NotNull(),
								"direction":         knownvalue.StringExact("INBOUND"),
								"type":              knownvalue.StringExact("BGP_IPv4_PREFIX_FILTER"),
								"attachment_status": knownvalue.StringExact(string(fabricv4.CONNECTIONROUTEAGGREGATIONDATAATTACHMENTSTATUS_PENDING_BGP_CONFIGURATION)),
							}),
						}),
					),

					newConnectionId.AddStateValue("data.equinix_fabric_connection_route_filters.test", tfjsonpath.New("connection_id")),
				},
			},
		},
	})

}

var connectionRouteFilterConfig = `
variable "port_uuid" {
  type = string
}

variable "vlan_tag" {
  type = number
}

resource "equinix_fabric_cloud_router" "test" {
  type = "XF_ROUTER"
  name = "RF_CR_PFCR"
  location {
    metro_code = "DC"
  }
  package {
    code = "STANDARD"
  }
  order {
    purchase_order_number = "1-234567"
    term_length           = 1
  }
  notifications {
    type = "ALL"
    emails = [
      "test@equinix.com",
      "test1@equinix.com"
    ]
  }
  project {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
  account {
    account_number = 201257
  }
}

resource "equinix_fabric_connection" "test" {
  type = "IP_VC"
  name = "RF_CR_Connection_PFCR"
  notifications {
    type   = "ALL"
    emails = ["test@equinix.com", "test1@equinix.com"]
  }
  order {
    purchase_order_number = "123485"
    term_length           = 1
  }
  bandwidth = 50
  redundancy {
    priority = "PRIMARY"
  }
  a_side {
    access_point {
      type = "CLOUD_ROUTER"
      router {
        uuid = equinix_fabric_cloud_router.test.id
      }
    }
  }
  project {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
  z_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = var.vlan_tag
      }
      location {
        metro_code = "DC"
      }
    }
  }
}

resource "equinix_fabric_route_filter" "test" {
  name = "rf_test_PFCR"
  project {
    project_id = "33ec651f-cc99-48e0-94d3-47466899cdc7"
  }
  type        = "BGP_IPv4_PREFIX_FILTER"
  description = "Route Filter Policy for X Purpose"
}

resource "equinix_fabric_route_filter_rule" "test" {
  route_filter_id = equinix_fabric_route_filter.test.id
  name            = "RF_Rule_PFCR"
  prefix          = "192.168.0.0/24"
  prefix_match    = "exact"
  description     = "Route Filter Rule for X Purpose"
}

resource "equinix_fabric_connection_route_filter" "test" {
  depends_on      = [equinix_fabric_route_filter_rule.test]
  connection_id   = equinix_fabric_connection.test.id
  route_filter_id = equinix_fabric_route_filter.test.id
  direction       = "INBOUND"
}

data "equinix_fabric_connection_route_filter" "test" {
  depends_on      = [equinix_fabric_connection_route_filter.test]
  connection_id   = equinix_fabric_connection.test.id
  route_filter_id = equinix_fabric_route_filter.test.id
}

data "equinix_fabric_connection_route_filters" "test" {
  depends_on    = [equinix_fabric_connection_route_filter.test]
  connection_id = equinix_fabric_connection.test.id
}
`

func CheckConnectionRouteFilterDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection_route_filter" {
			continue
		}

		connectionID := rs.Primary.Attributes["connection_id"]

		err := connection_route_filter.WaitForDeletion(ctx, connectionID, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for route filter deletion. ID: %s, Err: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
