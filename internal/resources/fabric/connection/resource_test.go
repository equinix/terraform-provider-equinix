package connection_test

import (
	"context"
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/equinix/terraform-provider-equinix/internal/resources/fabric/connection"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_fabric_connection_PNFV", &resource.Sweeper{
		Name: "equinix_fabric_connection",
		F:    testSweepConnections,
	})
}

func testSweepConnections(region string) error {
	return nil
}

func TestAccFabricCreatePort2SPConnection_PPDS(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	connectionsTestData := testing_helpers.GetFabricEnvConnectionTestData(t)
	var publicSPName, portUuid string
	if len(ports) > 0 && len(connectionsTestData) > 0 {
		publicSPName = connectionsTestData["ppds"]["publicSPName"]
		portUuid = ports["ppds"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2SPConnectionConfig(publicSPName, "port2sp_PPDS", portUuid, "CH"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port2sp_PPDS"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "1-323292"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", "2019"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "SP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.type", "L2_PROFILE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.name", publicSPName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.location.0.metro_code", "CH"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func TestAccFabricCreatePort2SPConnection_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	connectionsTestData := testing_helpers.GetFabricEnvConnectionTestData(t)
	var publicSPName, portUuid string
	if len(ports) > 0 && len(connectionsTestData) > 0 {
		publicSPName = connectionsTestData["pfcr"]["publicSPName"]
		portUuid = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2SPConnectionConfig(publicSPName, "port2sp_PFCR", portUuid, "SV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port2sp_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "1-323292"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", "2019"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "SP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.type", "L2_PROFILE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.name", publicSPName),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.location.0.metro_code", "SV"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricCreatePort2SPConnectionConfig(spName, name, portUuid, zSideMetro string) string {
	return fmt.Sprintf(`

	data "equinix_fabric_service_profiles" "this" {
	  filter {
		property = "/name"
		operator = "="
		values   = ["%s"]
	  }
	}


	resource "equinix_fabric_connection" "test" {
		name = "%s"
		type = "EVPL_VC"
		notifications{
			type="ALL" 
			emails=["example@equinix.com"]
		} 
		bandwidth = 50
		redundancy {priority= "PRIMARY"}
		order {
			purchase_order_number= "1-323292"
		}
		a_side {
			access_point {
				type= "COLO"
				port {
					uuid= "%s"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= "2019"
				}
			}
		}
		z_side {
			access_point {
				type= "SP"
				profile {
					type= "L2_PROFILE"
					uuid= data.equinix_fabric_service_profiles.this.data.0.uuid
				}
				location {
					metro_code= "%s"
				}
			}
		}
	}`, spName, name, portUuid, zSideMetro)
}

func TestAccFabricCreatePort2PortConnection_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var aSidePortUuid, zSidePortUuid string
	if len(ports) > 0 {
		aSidePortUuid = ports["pfcr"]["dot1q"][0].GetUuid()
		zSidePortUuid = ports["pfcr"]["dot1q"][1].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2PortConnectionConfig(50, aSidePortUuid, zSidePortUuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port_test_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "1-129105284100"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", "2397"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.vlan_tag", "2398"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccFabricCreatePort2PortConnectionConfig(100, aSidePortUuid, zSidePortUuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port_test_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "100"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "1-129105284100"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", "2397"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.vlan_tag", "2398"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})

}

func testAccFabricCreatePort2PortConnectionConfig(bandwidth int32, aSidePortUuid, zSidePortUuid string) string {
	return fmt.Sprintf(`resource "equinix_fabric_connection" "test" {
		type = "EVPL_VC"
		name = "port_test_PFCR"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
			purchase_order_number = "1-129105284100"
		}
		bandwidth = %d
		a_side {
			access_point {
				type = "COLO"
				port {
				 uuid = "%s"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= 2397
				}
				location {
					metro_code = "SV"
				}
			}
		}
		z_side {
			access_point {
				type = "COLO"
				port{
				 uuid = "%s"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= 2398
				}
				location {
					metro_code= "SV"
				}
			}
		}
	}`, bandwidth, aSidePortUuid, zSidePortUuid)
}

func TestAccFabricCreateCloudRouter2PortConnection_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][1].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateCloudRouter2PortConnectionConfig("fcr_test_PFCR", portUuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "fcr_test_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "IP_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "123485"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "project.0.project_id", "291639000636552"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "CLOUD_ROUTER"),
					resource.TestCheckResourceAttrSet(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.router.0.uuid"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.vlan_tag", "2325"),
				),
			},
		},
	})
}

func testAccFabricCreateCloudRouter2PortConnectionConfig(name, portUuid string) string {
	return fmt.Sprintf(`

	resource "equinix_fabric_cloud_router" "this" {
		type = "XF_ROUTER"
		name = "Test_PFCR"
		location{
			metro_code  = "SV"
		}
		order{
			purchase_order_number = "1-234567"
		}
		notifications{
			type = "ALL"
			emails = [
				"test@equinix.com",
				"test1@equinix.com"
			]
		}
		project{
			project_id = "291639000636552"
		}
		account {
			account_number = 201257
		}
		package {
			code = "STANDARD"
		}
	}

	resource "equinix_fabric_connection" "test" {
		type = "IP_VC"
		name = "%s"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
			purchase_order_number = "123485"
		}
		bandwidth = 50
		redundancy {
			priority= "PRIMARY"
		}
		a_side {
			access_point {
				type = "CLOUD_ROUTER"
				router {
					uuid = equinix_fabric_cloud_router.this.id
				}
			}
		}
		project{
		   project_id = "291639000636552"
		}
		z_side {
			access_point {
				type = "COLO"
				port{
					uuid = "%s"
				}
				link_protocol {
					type= "DOT1Q"
					vlan_tag= 2325
				}
				location {
					metro_code = "SV"
				}
			}
		}
	}`, name, portUuid)
}

func TestAccFabricCreateVirtualDevice2NetworkConnection_PNFV(t *testing.T) {
	connectionTestData := testing_helpers.GetFabricEnvConnectionTestData(t)
	var virtualDevice string
	if len(connectionTestData) > 0 {
		virtualDevice = connectionTestData["pnfv"]["virtualDevice"]
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateVirtualDevice2NetworkConnectionConfig("vd2network_PNFV", virtualDevice),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "vd2network_PNFV"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPLAN_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "123485"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "VD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.virtual_device.0.type", "EDGE"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.virtual_device.0.uuid", virtualDevice),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.interface.0.type", "CLOUD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.interface.0.id", "6"),
					resource.TestCheckResourceAttrSet(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "NETWORK"),
					resource.TestCheckResourceAttrSet(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.network.0.uuid"),
				),
			},
		},
	})

}

func testAccFabricCreateVirtualDevice2NetworkConnectionConfig(name, virtualDeviceUuid string) string {
	return fmt.Sprintf(`

	resource "equinix_fabric_network" "this" {
		type = "EVPLAN"
		name = "Tf_Network_PNFV"
		scope = "REGIONAL"
		notifications {
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		location {
			region = "AMER"
		}
		project{
			project_id = "4f855852-eb47-4721-8e40-b386a3676abf"
		}
	}

	resource "equinix_fabric_connection" "test" {
		type = "EVPLAN_VC"
		name = "%s"
		notifications{
			type = "ALL"
			emails = ["test@equinix.com","test1@equinix.com"]
		}
		order {
			purchase_order_number = "123485"
		}
		bandwidth = 50
		redundancy {
			priority= "PRIMARY"
		}
		a_side {
			access_point {
				type = "VD"
				virtual_device {
					type = "EDGE"
					uuid = "%s"
				}
				interface {
					type = "CLOUD"
					id = 6
				}
			}
		}
		z_side {
			access_point {
				type = "NETWORK"
				network {
					uuid = equinix_fabric_network.this.id
				}
			}
		}
	}`, name, virtualDeviceUuid)
}

func CheckConnectionDelete(s *terraform.State) error {
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection" {
			continue
		}

		err := connection.WaitUntilConnectionDeprovisioned(rs.Primary.ID, acceptance.TestAccProvider.Meta(), &schema.ResourceData{}, ctx, 10*time.Minute)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
