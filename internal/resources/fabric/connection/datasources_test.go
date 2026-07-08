package connection_test

import (
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricDataSourceConnection_PFCR(t *testing.T) {
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: dataSourceConnectionConfig,
				ConfigVariables: config.Variables{
					"aside_vlan":      config.IntegerVariable(asideVlan),
					"aside_port_uuid": config.StringVariable(aSidePortUUID),
					"zside_vlan":      config.IntegerVariable(zsideVlan),
					"zside_port_uuid": config.StringVariable(zSidePortUUID),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_connections.connections", tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":      knownvalue.StringExact("ds_con_test_PFCR"),
								"bandwidth": knownvalue.Int32Exact(50),
								"type":      knownvalue.StringExact("EVPL_VC"),
								"redundancy": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"priority": knownvalue.StringExact("PRIMARY"),
									}),
								}),
								"order": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"purchase_order_number": knownvalue.StringExact("1-129105284100"),
									}),
								}),
								"a_side": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"access_point": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type": knownvalue.StringExact("COLO"),
												"link_protocol": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectPartial(map[string]knownvalue.Check{
														"type":     knownvalue.StringExact("DOT1Q"),
														"vlan_tag": knownvalue.Int32Exact(int32(asideVlan)),
													}),
												}),
												"location": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectPartial(map[string]knownvalue.Check{
														"metro_code": knownvalue.StringExact("DC"),
													}),
												}),
											}),
										}),
									}),
								}),
								"z_side": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"access_point": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectPartial(map[string]knownvalue.Check{
												"type": knownvalue.StringExact("COLO"),
												"link_protocol": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectPartial(map[string]knownvalue.Check{
														"type":     knownvalue.StringExact("DOT1Q"),
														"vlan_tag": knownvalue.Int32Exact(int32(zsideVlan)),
													}),
												}),
												"location": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectPartial(map[string]knownvalue.Check{
														"metro_code": knownvalue.StringExact("SV"),
													}),
												}),
											}),
										}),
									}),
								}),
							}),
						}),
					),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

var dataSourceConnectionConfig = `
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

resource "equinix_fabric_connection" "test" {

  type = "EVPL_VC"
  name = "ds_con_test_PFCR"
  notifications {
    type   = "ALL"
    emails = ["test@equinix.com", "test1@equinix.com"]
  }
  order {
    purchase_order_number = "1-129105284100"
  }
  bandwidth = 50
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

data "equinix_fabric_connection" "test" {
  uuid = equinix_fabric_connection.test.id
}

data "equinix_fabric_connections" "connections" {
  outer_operator = "AND"
  filter {
    property = "/name"
    operator = "="
    values   = ["ds_con_test_PFCR"]
  }
  filter {
    property = "/uuid"
    operator = "="
    values   = [equinix_fabric_connection.test.id]
  }
}
`
