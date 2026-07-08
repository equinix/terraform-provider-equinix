package servicetoken_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/servicetoken"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFabricZsideVirtualDeviceServiceToken_PNFV(t *testing.T) {
	connectionTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var virtualDevice string
	if len(connectionTestData) > 0 {
		virtualDevice = connectionTestData["pnfv"]["virtualDevice"]
	}
	expiration := time.Now().Add(30 * 24 * time.Hour).UTC().Format("2006-01-02T15:04:05.000Z")
	serviceTokenName, serviceTokenUpdatedName := "Service_token_PNFV", "UP_Service_Token_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "zside vd token", "Updated zside vd token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricZsideVirtualDeviceServiceTokenConfig(serviceTokenName, serviceTokenDescription, expiration, virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", expiration),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id"),
				),
			},
			{
				Config: testAccFabricZsideVirtualDeviceServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription, expiration, virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", expiration),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id"),
				),
			},
		},
	})
}

func TestAccFabricAsidePortServiceToken_PNFV(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUUID string
	if len(ports) > 0 {
		portUUID = ports["pnfv"]["dot1q"][0].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	expiration := time.Now().Add(30 * 24 * time.Hour).UTC().Format("2006-01-02T15:04:05.000Z")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: asidePortServiceTokenConfig,
				ConfigVariables: config.Variables{
					"port_uuid":   config.StringVariable(portUUID),
					"vlan_tag":    config.IntegerVariable(targetVlan),
					"name":        config.StringVariable("token_port_PNFV"),
					"description": config.StringVariable("aside port token"),
					"expiration":  config.StringVariable(expiration),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_service_token.test", map[string]knownvalue.Check{
						"uuid":                 knownvalue.NotNull(),
						"name":                 knownvalue.StringExact("token_port_PNFV"),
						"description":          knownvalue.StringExact("aside port token"),
						"expiration_date_time": knownvalue.StringExact(expiration),
						"type":                 knownvalue.StringExact("VC_TOKEN"),
					}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"bandwidth_limit": knownvalue.Int32Exact(int32(1000)),
						}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("a_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("link_protocol").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"type":     knownvalue.StringExact("DOT1Q"),
							"vlan_tag": knownvalue.Int32Exact(int32(targetVlan)),
						}),
				},

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.port.0.uuid", portUUID),
				),
			},
			{
				Config: asidePortServiceTokenConfig,
				ConfigVariables: config.Variables{
					"port_uuid":   config.StringVariable(portUUID),
					"vlan_tag":    config.IntegerVariable(targetVlan),
					"name":        config.StringVariable("UP_Token_port_PNFV"),
					"description": config.StringVariable("Updated aside port token"),
					"expiration":  config.StringVariable(expiration),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_service_token.test", map[string]knownvalue.Check{
						"uuid":                 knownvalue.NotNull(),
						"name":                 knownvalue.StringExact("UP_Token_port_PNFV"),
						"description":          knownvalue.StringExact("Updated aside port token"),
						"expiration_date_time": knownvalue.StringExact(expiration),
						"type":                 knownvalue.StringExact("VC_TOKEN"),
					}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"bandwidth_limit": knownvalue.Int32Exact(int32(1000)),
						}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("a_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("link_protocol").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"type":     knownvalue.StringExact("DOT1Q"),
							"vlan_tag": knownvalue.Int32Exact(int32(targetVlan)),
						}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("a_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("port").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"uuid": knownvalue.StringExact(portUUID),
						}),
				},
			},
		},
	})
}

func TestAccFabricZsidePortServiceToken_PNFV(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var portUUID string
	if len(ports) > 0 {
		portUUID = ports["pnfv"]["dot1q"][0].GetUuid()
	}

	targetVlan, err := testinghelpers.RandomVlan(portUUID)
	if err != nil {
		t.Fatalf("unable to get a available VLAN: %s", err)
		return
	}

	expiration := time.Now().Add(30 * 24 * time.Hour).UTC().Format("2006-01-02T15:04:05.000Z")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: zsidePortServiceTokenConfig,
				ConfigVariables: config.Variables{
					"port_uuid":   config.StringVariable(portUUID),
					"vlan_tag":    config.IntegerVariable(targetVlan),
					"name":        config.StringVariable("token_port_PNFV"),
					"description": config.StringVariable("zside port token"),
					"expiration":  config.StringVariable(expiration),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_service_token.test", map[string]knownvalue.Check{
						"uuid":                 knownvalue.NotNull(),
						"name":                 knownvalue.StringExact("token_port_PNFV"),
						"description":          knownvalue.StringExact("zside port token"),
						"expiration_date_time": knownvalue.StringExact(expiration),
						"type":                 knownvalue.StringExact("VC_TOKEN"),
					}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"supported_bandwidths": knownvalue.ListSizeExact(3),
						}),

					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("z_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("link_protocol").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"type":     knownvalue.StringExact("DOT1Q"),
							"vlan_tag": knownvalue.Int32Exact(int32(targetVlan)),
						}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("z_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("port").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"uuid": knownvalue.StringExact(portUUID),
						}),
				},
			},
			{
				Config: zsidePortServiceTokenConfig,
				ConfigVariables: config.Variables{
					"port_uuid":   config.StringVariable(portUUID),
					"vlan_tag":    config.IntegerVariable(targetVlan),
					"name":        config.StringVariable("UP_Token_port_PNFV"),
					"description": config.StringVariable("Updated zside port token"),
					"expiration":  config.StringVariable(expiration),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					testinghelpers.ExpectKnownAttributes("equinix_fabric_service_token.test", map[string]knownvalue.Check{
						"uuid":                 knownvalue.NotNull(),
						"name":                 knownvalue.StringExact("UP_Token_port_PNFV"),
						"description":          knownvalue.StringExact("Updated zside port token"),
						"expiration_date_time": knownvalue.StringExact(expiration),
						"type":                 knownvalue.StringExact("VC_TOKEN"),
					}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"supported_bandwidths": knownvalue.ListSizeExact(3),
						}),
					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("z_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("port").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"uuid": knownvalue.StringExact(portUUID),
						}),

					testinghelpers.ExpectKnownAttributesAt("equinix_fabric_service_token.test",
						tfjsonpath.New("service_token_connection").AtSliceIndex(0).AtMapKey("z_side").AtSliceIndex(0).AtMapKey("access_point_selectors").AtSliceIndex(0).AtMapKey("link_protocol").AtSliceIndex(0),
						map[string]knownvalue.Check{
							"type":     knownvalue.StringExact("DOT1Q"),
							"vlan_tag": knownvalue.Int32Exact(int32(targetVlan)),
						}),
				},
			},
		},
	})
}

func TestAccFabricZsideNetworkServiceToken_PNFV(t *testing.T) {
	connectionTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var networkUUID string
	if len(connectionTestData) > 0 {
		networkUUID = connectionTestData["pnfv"]["network"]
	}
	expiration := time.Now().Add(30 * 24 * time.Hour).UTC().Format("2006-01-02T15:04:05.000Z")
	serviceTokenName, serviceTokenUpdatedName := "token_zwan_PNFV", "UP_Token_zwan_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "zside network token", "Updated zside network token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricZsideNetworkServiceTokenConfig(serviceTokenName, serviceTokenDescription, expiration, networkUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", expiration),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.network.0.uuid", networkUUID),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricZsideNetworkServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription, expiration, networkUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", expiration),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.network.0.uuid", networkUUID),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

var asidePortServiceTokenConfig = `
variable "port_uuid" {
  type = string
}

variable "vlan_tag" {
  type = number
}

variable "name" {
  type = string
}

variable "description" {
  type = string
}

variable "expiration" {
  type = string
}

resource "equinix_fabric_service_token" "test" {
  type                 = "VC_TOKEN"
  name                 = var.name
  description          = var.description
  expiration_date_time = var.expiration
  service_token_connection {
    type            = "EVPL_VC"
    bandwidth_limit = 1000
    a_side {
      access_point_selectors {
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
  }
  notifications {
    type   = "ALL"
    emails = ["panthers_auto@equinix.com", "test1@equinix.com", "example@equinix.com"]
  }

}
`

var zsidePortServiceTokenConfig = `
variable "port_uuid" {
  type = string
}

variable "vlan_tag" {
  type = number
}

variable "name" {
  type = string
}

variable "description" {
  type = string
}

variable "expiration" {
  type = string
}


resource "equinix_fabric_service_token" "test" {
  type                 = "VC_TOKEN"
  name                 = var.name
  description          = var.description
  expiration_date_time = var.expiration
  service_token_connection {
    type                 = "EVPL_VC"
    supported_bandwidths = [50, 200, 10000]
    z_side {
      access_point_selectors {
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
  }
  notifications {
    type   = "ALL"
    emails = ["panthers_auto@equinix.com", "test1@equinix.com", "example@equinix.com"]
  }

}
`

func testAccFabricZsideVirtualDeviceServiceTokenConfig(serviceTokenName string, serviceTokenDescription string, expiration string, virtualDeviceUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			name = "%s"
			description = "%s"
			expiration_date_time = "%s"
			service_token_connection {
				type = "EVPL_VC"
				supported_bandwidths = [50, 200, 10000]
				z_side {
					access_point_selectors{
						type = "VD"
						virtual_device{
							type = "EDGE"
							uuid = "%s"
						}
						interface{
							type = "NETWORK"
						}
					}
				}
			}
			notifications {
    			type   = "ALL"
    			emails = ["panthers_auto@equinix.com", "test1@equinix.com", "example@equinix.com"]
  			}

		}
    `, serviceTokenName, serviceTokenDescription, expiration, virtualDeviceUUID)
}

func testAccFabricZsideNetworkServiceTokenConfig(serviceTokenName string, serviceTokenDescription string, expiration string, networkUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test" {
						type = "VC_TOKEN"
						name = "%s"
						description = "%s"
						expiration_date_time = "%s"
						service_token_connection {
							type = "EVPLAN_VC"
							supported_bandwidths = [50, 200, 10000]
							z_side {
								access_point_selectors{
									type = "NETWORK"
									 network {
										uuid = "%s"
									}
								}
							}
						}
						notifications {
    						type   = "ALL"
    						emails = ["example@equinix.com", "test1@equinix.com"]
  						}
				   }
			`, serviceTokenName, serviceTokenDescription, expiration, networkUUID)
}

func CheckServiceTokenDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_service_token" {
			continue
		}

		err := servicetoken.WaitForDeletion(ctx, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for service token deletion. ID: %s, Err: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
