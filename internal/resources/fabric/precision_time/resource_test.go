// Package precisiontime_test for EPT resources and data sources tests
package precisiontime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	eqconfig "github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricCreatePort2EPT_NTPConfiguration_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][0].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUuid)

	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}
	newConnectionId := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkEptServiceDelete,
		Steps: []resource.TestStep{
			{
				Config: port2EPTNPTConfig,
				ConfigVariables: config.Variables{
					"port_uuid": config.StringVariable(portUuid),
					"vlan_tag":  config.IntegerVariable(targetVlan),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_precision_time_service.ntp", map[string]knownvalue.Check{
						"id":   knownvalue.NotNull(),
						"name": knownvalue.StringExact("tf_acc_eptntp_PFCR"),
						"type": knownvalue.StringExact("NTP"),
					}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_precision_time_service.ntp",
						tfjsonpath.New("package"),
						map[string]knownvalue.Check{
							"href": knownvalue.NotNull(),
							"code": knownvalue.StringExact("NTP_STANDARD"),
						}),

					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_precision_time_service.ntp",
						tfjsonpath.New("ipv4"),
						map[string]knownvalue.Check{
							"primary":         knownvalue.StringExact("192.168.254.241"),
							"secondary":       knownvalue.StringExact("192.168.254.242"),
							"network_mask":    knownvalue.StringExact("255.255.255.240"),
							"default_gateway": knownvalue.StringExact("192.168.254.254"),
						}),

					testinghelpers.ExpectKnownAttributes("data.equinix_fabric_precision_time_service.ntp", map[string]knownvalue.Check{
						"uuid": knownvalue.NotNull(),
						"name": knownvalue.StringExact("tf_acc_eptntp_PFCR"),
						"type": knownvalue.StringExact("NTP"),
					}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ntp",
						tfjsonpath.New("package"),
						map[string]knownvalue.Check{
							"href": knownvalue.NotNull(),
							"code": knownvalue.StringExact("NTP_STANDARD"),
						}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ntp",
						tfjsonpath.New("project"),
						map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"), // inherited from port
						}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ntp",
						tfjsonpath.New("ipv4"),
						map[string]knownvalue.Check{
							"primary":         knownvalue.StringExact("192.168.254.241"),
							"secondary":       knownvalue.StringExact("192.168.254.242"),
							"network_mask":    knownvalue.StringExact("255.255.255.240"),
							"default_gateway": knownvalue.StringExact("192.168.254.254"),
						}),

					newConnectionId.AddStateValue("equinix_fabric_connection.test", tfjsonpath.New("uuid")),
					newConnectionId.AddStateValue("equinix_fabric_precision_time_service.ntp", tfjsonpath.New("connections").AtSliceIndex(0).AtMapKey("uuid")),
					newConnectionId.AddStateValue("data.equinix_fabric_precision_time_service.ntp", tfjsonpath.New("connections").AtSliceIndex(0).AtMapKey("uuid")),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

var port2EPTNPTConfig = `
variable "port_uuid" {
  type = string
}

variable "vlan_tag" {
  type = number
}

data "equinix_fabric_service_profiles" "sp" {
  filter {
    property = "/name"
    operator = "="
    values   = ["Equinix Precision Time NTP UAT Global"]
  }
}


resource "equinix_fabric_connection" "test" {
  name = "port2eptntp_PFCR"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 1
  redundancy { priority = "PRIMARY" }
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = var.vlan_tag
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.sp.data.0.uuid
      }
      location {
        metro_code = "SV"
      }
    }
  }
}

resource "equinix_fabric_precision_time_service" "ntp" {
  type = "NTP"
  name = "tf_acc_eptntp_PFCR"
  package = {
    code = "NTP_STANDARD"
  }
  connections = [
    {
      uuid = equinix_fabric_connection.test.id
    }
  ]
  ipv4 = {
    primary         = "192.168.254.241"
    secondary       = "192.168.254.242"
    network_mask    = "255.255.255.240"
    default_gateway = "192.168.254.254"
  }
}

data "equinix_fabric_precision_time_service" "ntp" {
  ept_service_id = equinix_fabric_precision_time_service.ntp.id
}
`

func TestAccFabricCreatePort2EPT_PTPConfiguration_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][1].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUuid)

	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	newConnectionId := statecheck.CompareValue(compare.ValuesSame())
	newServiceId := statecheck.CompareValue(compare.ValuesSame())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkEptServiceDelete,
		Steps: []resource.TestStep{
			{
				Config: port2EPTPTPConfig,
				ConfigVariables: config.Variables{
					"port_uuid": config.StringVariable(portUuid),
					"vlan_tag":  config.IntegerVariable(targetVlan),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_precision_time_service.ptp", map[string]knownvalue.Check{
						"id":   knownvalue.NotNull(),
						"name": knownvalue.StringExact("tf_acc_eptptp_PFCR"),
						"type": knownvalue.StringExact("PTP"),
					}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_precision_time_service.ptp",
						tfjsonpath.New("package"),
						map[string]knownvalue.Check{
							"href": knownvalue.NotNull(),
							"code": knownvalue.StringExact("PTP_STANDARD"),
						}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ptp",
						tfjsonpath.New("project"),
						map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"), // inherited from port
						}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_precision_time_service.ptp",
						tfjsonpath.New("ipv4"),
						map[string]knownvalue.Check{
							"primary":         knownvalue.StringExact("192.168.254.241"),
							"secondary":       knownvalue.StringExact("192.168.254.242"),
							"network_mask":    knownvalue.StringExact("255.255.255.240"),
							"default_gateway": knownvalue.StringExact("192.168.254.254"),
						}),

					testinghelpers.ExpectKnownAttributes("data.equinix_fabric_precision_time_service.ptp", map[string]knownvalue.Check{
						"uuid": knownvalue.NotNull(),
						"name": knownvalue.StringExact("tf_acc_eptptp_PFCR"),
						"type": knownvalue.StringExact("PTP"),
					}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ptp",
						tfjsonpath.New("package"),
						map[string]knownvalue.Check{
							"href": knownvalue.NotNull(),
							"code": knownvalue.StringExact("PTP_STANDARD"),
						}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ptp",
						tfjsonpath.New("project"),
						map[string]knownvalue.Check{
							"project_id": knownvalue.StringExact("33ec651f-cc99-48e0-94d3-47466899cdc7"), // inherited from port
						}),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_service.ptp",
						tfjsonpath.New("ipv4"),
						map[string]knownvalue.Check{
							"primary":         knownvalue.StringExact("192.168.254.241"),
							"secondary":       knownvalue.StringExact("192.168.254.242"),
							"network_mask":    knownvalue.StringExact("255.255.255.240"),
							"default_gateway": knownvalue.StringExact("192.168.254.254"),
						}),

					statecheck.ExpectKnownValue("data.equinix_fabric_precision_time_services.all", tfjsonpath.New("data"), knownvalue.ListSizeExact(1)),
					testinghelpers.ExpectKnownAttributesAt("data.equinix_fabric_precision_time_services.all",
						tfjsonpath.New("pagination"),
						map[string]knownvalue.Check{
							"limit":  knownvalue.Int32Exact(2),
							"offset": knownvalue.Int32Exact(0),
						}),

					newConnectionId.AddStateValue("equinix_fabric_connection.test", tfjsonpath.New("uuid")),
					newConnectionId.AddStateValue("equinix_fabric_precision_time_service.ptp", tfjsonpath.New("connections").AtSliceIndex(0).AtMapKey("uuid")),
					newConnectionId.AddStateValue("data.equinix_fabric_precision_time_service.ptp", tfjsonpath.New("connections").AtSliceIndex(0).AtMapKey("uuid")),

					newServiceId.AddStateValue("equinix_fabric_precision_time_service.ptp", tfjsonpath.New("uuid")),
					newServiceId.AddStateValue("data.equinix_fabric_precision_time_service.ptp", tfjsonpath.New("uuid")),
					newServiceId.AddStateValue("data.equinix_fabric_precision_time_services.all", tfjsonpath.New("data").AtSliceIndex(0).AtMapKey("uuid")),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

var port2EPTPTPConfig = `
variable "port_uuid" {
  type = string
}

variable "vlan_tag" {
  type = number
}

data "equinix_fabric_service_profiles" "sp" {
  filter {
    property = "/name"
    operator = "="
    values   = ["Equinix Precision Time PTP Global UAT"]
  }
}

resource "equinix_fabric_connection" "test" {
  name = "port2eptptp_PFCR"
  type = "EVPL_VC"
  notifications {
    type   = "ALL"
    emails = ["example@equinix.com"]
  }
  bandwidth = 5
  redundancy { priority = "PRIMARY" }
  order {
    purchase_order_number = "1-323292"
  }
  a_side {
    access_point {
      type = "COLO"
      port {
        uuid = var.port_uuid
      }
      link_protocol {
        type     = "DOT1Q"
        vlan_tag = var.vlan_tag
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      profile {
        type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profiles.sp.data.0.uuid
      }
      location {
        metro_code = "SV"
      }
    }
  }
}

resource "equinix_fabric_precision_time_service" "ptp" {
  type = "PTP"
  name = "tf_acc_eptptp_PFCR"
  package = {
    code = "PTP_STANDARD"
  }
  connections = [
    {
      uuid = equinix_fabric_connection.test.id
    }
  ]
  ipv4 = {
    primary         = "192.168.254.241"
    secondary       = "192.168.254.242"
    network_mask    = "255.255.255.240"
    default_gateway = "192.168.254.254"
  }
}

data "equinix_fabric_precision_time_service" "ptp" {
  ept_service_id = equinix_fabric_precision_time_service.ptp.id
}

data "equinix_fabric_precision_time_services" "all" {
  pagination = {
    limit  = 2
    offset = 0
  }
  filters = [{
    property = "/uuid"
    operator = "="
    values   = [equinix_fabric_precision_time_service.ptp.uuid]
  }]
  sort = [{
    direction = "DESC"
    property  = "/uuid"
  }]
}
`

func checkEptServiceDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*eqconfig.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_precision_time_service" {
			continue
		}

		if eptService, _, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, rs.Primary.ID).Execute(); err == nil {
			if eptService.GetState() == fabricv4.PRECISIONTIMESERVICERESPONSESTATE_PROVISIONED {
				return fmt.Errorf("fabric EPT service %s still exists and is %s",
					rs.Primary.ID, eptService.GetState())
			}
		}
	}
	return nil
}
