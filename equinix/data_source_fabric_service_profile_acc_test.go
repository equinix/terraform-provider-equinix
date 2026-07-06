package equinix_test

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

const testAccFabricReadServiceProfileConfig = `
variable "port_uuid" {
  type = string
}

variable "port_type" {
  type = string
}

variable "metro_code" {
  type = string
}

resource "equinix_fabric_service_profile" "test" {
  name         = "SP_DataSource_PFCR"
  description  = "Generic Read SP"
  type         = "L2_PROFILE"
  notifications {
    emails = ["opsuser100@equinix.com"]
    type   = "BANDWIDTH_ALERT"
  }
  tags           = ["VoIP", "Saas"]
  visibility     = "PRIVATE"
  allowed_emails = ["panthersfcr@test.com"]
  ports {
    uuid = var.port_uuid
    type = var.port_type
    location {
      metro_code = var.metro_code
    }
    cross_connect_id          = ""
    seller_region             = ""
    seller_region_description = ""
  }
  access_point_type_configs {
    type                             = "COLO"
    connection_redundancy_required   = false
    allow_bandwidth_auto_approval    = false
    allow_remote_connections         = false
    connection_label                 = "test"
    enable_auto_generate_service_key = false
    bandwidth_alert_threshold        = 10
    allow_custom_bandwidth           = true
    api_config {
      api_available        = false
      equinix_managed_vlan = true
      bandwidth_from_api   = false
      integration_id       = "test"
      equinix_managed_port = true
    }
    authentication_key {
      required    = false
      label       = "Service Key"
      description = "XYZ"
    }
    supported_bandwidths = [100, 500]
  }
  marketing_info {
    promotion = false
  }
}

data "equinix_fabric_service_profile" "test" {
  uuid = equinix_fabric_service_profile.test.uuid
}

data "equinix_fabric_service_profiles" "test" {
  and_filters = true
  filter {
    property = "/uuid"
    operator = "="
    values   = [equinix_fabric_service_profile.test.uuid]
  }
}
`

func TestAccFabricServiceProfileDataSources_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)

	var portUUID, portMetroCode, portType string
	if len(ports) > 0 {
		port := ports["pfcr"]["dot1q"][0]
		portUUID = port.GetUuid()
		portMetroCodeLocation := port.GetLocation()
		portMetroCode = portMetroCodeLocation.GetMetroCode()
		portType = string(port.GetType())
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkServiceProfileDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadServiceProfileConfig,
				ConfigVariables: config.Variables{
					"port_uuid":  config.StringVariable(portUUID),
					"port_type":  config.StringVariable(portType),
					"metro_code": config.StringVariable(portMetroCode),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.equinix_fabric_service_profile.test", tfjsonpath.New("name"), knownvalue.StringExact("SP_DataSource_PFCR")),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_service_profile.test",
						tfjsonpath.New("uuid"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue("data.equinix_fabric_service_profile.test", tfjsonpath.New("description"), knownvalue.StringExact("Generic Read SP")),
					statecheck.ExpectKnownValue("data.equinix_fabric_service_profile.test", tfjsonpath.New("state"), knownvalue.StringExact("ACTIVE")),
					statecheck.ExpectKnownValue("data.equinix_fabric_service_profile.test", tfjsonpath.New("visibility"), knownvalue.StringExact("PRIVATE")),
					statecheck.ExpectKnownValue("data.equinix_fabric_service_profile.test", tfjsonpath.New("href"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue("data.equinix_fabric_service_profile.test", tfjsonpath.New("self_profile"), knownvalue.Bool(true)),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_service_profile.test",
						tfjsonpath.New("access_point_type_configs"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"uuid":                             knownvalue.NotNull(),
								"type":                             knownvalue.NotNull(),
								"allow_remote_connections":         knownvalue.NotNull(),
								"allow_custom_bandwidth":           knownvalue.NotNull(),
								"enable_auto_generate_service_key": knownvalue.NotNull(),
								"connection_redundancy_required":   knownvalue.NotNull(),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"data.equinix_fabric_service_profile.test",
						tfjsonpath.New("metros"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"code":         knownvalue.StringExact(portMetroCode),
								"name":         knownvalue.NotNull(),
								"display_name": knownvalue.NotNull(),
							}),
						}),
					),

					statecheck.ExpectKnownValue(
						"data.equinix_fabric_service_profiles.test",
						tfjsonpath.New("data"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectPartial(map[string]knownvalue.Check{
								"name":           knownvalue.StringExact("SP_DataSource_PFCR"),
								"type":           knownvalue.StringExact("L2_PROFILE"),
								"uuid":           knownvalue.NotNull(),
								"description":    knownvalue.StringExact("Generic Read SP"),
								"state":          knownvalue.StringExact("ACTIVE"),
								"visibility":     knownvalue.StringExact("PRIVATE"),
								"href":           knownvalue.NotNull(),
								"notifications":  knownvalue.NotNull(),
								"marketing_info": knownvalue.NotNull(),
								"allowed_emails": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("panthersfcr@test.com"),
								}),
								"tags": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("VoIP"),
									knownvalue.StringExact("Saas"),
								}),
								"account":    knownvalue.NotNull(),
								"change_log": knownvalue.NotNull(),
								"access_point_type_configs": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"uuid":                             knownvalue.NotNull(),
										"type":                             knownvalue.NotNull(),
										"allow_remote_connections":         knownvalue.NotNull(),
										"allow_custom_bandwidth":           knownvalue.NotNull(),
										"enable_auto_generate_service_key": knownvalue.NotNull(),
										"connection_redundancy_required":   knownvalue.NotNull(),
									}),
								}),
								"metros": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectPartial(map[string]knownvalue.Check{
										"name":         knownvalue.NotNull(),
										"display_name": knownvalue.NotNull(),
									}),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}
