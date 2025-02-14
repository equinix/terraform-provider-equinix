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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricZsideVirtualDeviceServiceToken_PNFV(t *testing.T) {
	connectionTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var virtualDevice string
	if len(connectionTestData) > 0 {
		virtualDevice = connectionTestData["pnfv"]["virtualDevice"]
	}
	serviceTokenName, serviceTokenUpdatedName := "Service_token_PNFV", "UP_Service_Token_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "zside vd token", "Updated zside vd token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricZsideVirtualDeviceServiceTokenConfig(serviceTokenName, serviceTokenDescription, virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-01-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccFabricZsideVirtualDeviceServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription, virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-01-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.interface.0.id"),
				),
				ExpectNonEmptyPlan: false,
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
	serviceTokenName, serviceTokenUpdatedName := "token_port_PNFV", "UP_Token_port_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "aside port token", "Updated aside port token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricAsidePortServiceTokenConfig(serviceTokenName, serviceTokenDescription, portUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-01-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.bandwidth_limit", "1000"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.port.0.uuid", portUUID),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.link_protocol.0.vlan_tag", "2987"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccFabricAsidePortServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription, portUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-01-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.bandwidth_limit", "1000"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.port.0.uuid", portUUID),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.a_side.0.access_point_selectors.0.link_protocol.0.vlan_tag", "2987"),
				),
				ExpectNonEmptyPlan: false,
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
	serviceTokenName, serviceTokenUpdatedName := "token_zport_PNFV", "UP_Token_zport_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "zside port token", "Updated zside port token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricZsidePortServiceTokenConfig(serviceTokenName, serviceTokenDescription, portUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-01-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.port.0.uuid", portUUID),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.link_protocol.0.vlan_tag", "2087"),
				),
				ExpectNonEmptyPlan: false,
			},
			{
				Config: testAccFabricZsidePortServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription, portUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-01-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.supported_bandwidths.#", "3"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.port.0.uuid", portUUID),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.link_protocol.0.vlan_tag", "2087"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccFabricZsideNetworkServiceToken_PNFV(t *testing.T) {
	connectionTestData := testinghelpers.GetFabricEnvConnectionTestData(t)
	var networkUUID string
	if len(connectionTestData) > 0 {
		networkUUID = connectionTestData["pfcr"]["network"]
	}
	serviceTokenName, serviceTokenUpdatedName := "token_zwan_PNFV", "UP_Token_zwan_PNFV"
	serviceTokenDescription, serviceTokenUpdatedDescription := "zside network token", "Updated zside network token"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckServiceTokenDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricZsideNetworkServiceTokenConfig(serviceTokenName, serviceTokenDescription, networkUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-02-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.network.0.uuid", networkUUID),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricZsideNetworkServiceTokenConfig(serviceTokenUpdatedName, serviceTokenUpdatedDescription, networkUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_service_token.test", "uuid"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "name", serviceTokenUpdatedName),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "type", "VC_TOKEN"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "description", serviceTokenUpdatedDescription),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "expiration_date_time", "2025-02-18T06:43:49.981Z"),
					resource.TestCheckResourceAttr("equinix_fabric_service_token.test", "service_token_connection.0.z_side.0.access_point_selectors.0.network.0.uuid", networkUUID),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricAsidePortServiceTokenConfig(serviceTokenName string, serviceTokenDescription string, portUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			name = "%s"
			description = "%s"
			expiration_date_time = "2025-01-18T06:43:49.981Z"
			service_token_connection {
				type = "EVPL_VC"
				bandwidth_limit = 1000
				a_side {
					access_point_selectors{
						type = "COLO"
						port {
							uuid = "%s"
						}
						link_protocol {
							type = "DOT1Q"
							vlan_tag = "2987"
						}
					}
				}
			}
			notifications {
    			type   = "ALL"
    			emails = ["example@equinix.com", "test1@equinix.com"]
  			}

		}
    `, serviceTokenName, serviceTokenDescription, portUUID)
}

func testAccFabricZsidePortServiceTokenConfig(serviceTokenName string, serviceTokenDescription string, portUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			name = "%s"
			description = "%s"
			expiration_date_time = "2025-01-18T06:43:49.981Z"
			service_token_connection {
				type = "EVPL_VC"
				supported_bandwidths = [50, 200, 10000]
				z_side {
					access_point_selectors{
						type = "COLO"
						port {
							uuid = "%s"
						}
						link_protocol {
							type = "DOT1Q"
							vlan_tag = "2087"
						}
					}
				}
			}
			notifications {
    			type   = "ALL"
    			emails = ["example@equinix.com", "test1@equinix.com"]
  			}

		}
    `, serviceTokenName, serviceTokenDescription, portUUID)
}

func testAccFabricZsideVirtualDeviceServiceTokenConfig(serviceTokenName string, serviceTokenDescription string, virtualDeviceUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test"{
			type = "VC_TOKEN"
			name = "%s"
			description = "%s"
			expiration_date_time = "2025-01-18T06:43:49.981Z"
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
    			emails = ["example@equinix.com", "test1@equinix.com"]
  			}

		}
    `, serviceTokenName, serviceTokenDescription, virtualDeviceUUID)
}

func testAccFabricZsideNetworkServiceTokenConfig(serviceTokenName string, serviceTokenDescription string, networkUUID string) string {
	return fmt.Sprintf(
		`resource "equinix_fabric_service_token" "test" {
						type = "VC_TOKEN"
						name = "%s"
						description = "%s"
						expiration_date_time = "2025-02-18T06:43:49.981Z"
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
			`, serviceTokenName, serviceTokenDescription, networkUUID)
}

func CheckServiceTokenDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_service_token" {
			continue
		}

		err := servicetoken.WaitForDeletion(ctx, rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
