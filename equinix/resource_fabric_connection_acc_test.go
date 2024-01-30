package equinix_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"os"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
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

const (
	FabricConnectionsTestDataEnvVar = "TF_ACC_FABRIC_CONNECTIONS_TEST_DATA"
)

func GetFabricEnvConnectionTestData(t *testing.T) map[string]map[string]string {
	connectionTestDataJson := os.Getenv(FabricConnectionsTestDataEnvVar)
	var connectionTestData map[string]map[string]string
	if err := json.Unmarshal([]byte(connectionTestDataJson), &connectionTestData); err != nil {
		t.Fatalf("Failed reading connection data from environment: %v, %s", err, connectionTestDataJson)
	}
	return connectionTestData
}

func TestAccFabricCreatePort2SPConnection_PFCR(t *testing.T) {
	t.Skip("Skipping while focused on port connection")
	ports := GetFabricEnvPorts(t)
	connectionsTestData := GetFabricEnvConnectionTestData(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2SPConnectionConfig(connectionsTestData["pfcr"]["publicSPName"], "port2sp_PFCR", ports["pfcr"]["dot1q"][0].Uuid),
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
						"equinix_fabric_connection.test", "z_side.0.access_point.0.profile.0.name", connectionsTestData["public-sp"]["spName"]),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.location.0.metro_code", "SV"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricCreatePort2SPConnectionConfig(spName, name, portUuid string) string {
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
					metro_code= "SV"
				}
			}
		}
	}`, spName, name, portUuid)
}

func TestAccFabricCreatePort2PortConnection_PFCR(t *testing.T) {
	t.Skip("Is Successful; skipping because of duration")
	ports := GetFabricEnvPorts(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2PortConnectionConfig(50, ports["pfcr"]["dot1q"][0].Uuid, ports["pfcr"]["dot1q"][1].Uuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port_test__PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
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
				Config: testAccFabricCreatePort2PortConnectionConfig(100, ports["pfcr"]["dot1q"][0].Uuid, ports["pfcr"]["dot1q"][1].Uuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "port_test_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "100"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
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
	t.Skip("Skipping while focused on port connection")
	ports := GetFabricEnvPorts(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateCloudRouter2PortConnectionConfig("fcr_test_PFCR", ports["pfcr"]["dot1q"][1].Uuid),
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
						"equinix_fabric_connection.test", "a_side.0.access_point.0.router"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.vlan_tag", "2325"),
				),
				ExpectNonEmptyPlan: true,
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
	t.Skip("Skipping while PFNV account gets fixed")
	connectionTestData := GetFabricEnvConnectionTestData(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateVirtualDevice2NetworkConnectionConfig("vd2network_PNFV", connectionTestData["pnfv"]["virtualDevice"], connectionTestData["pnfv"]["network"]),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "name", "vd2network_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "order.0.purchase_order_number", "123485"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "project.0.project_id", "291639000636552"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.type", "VD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.virtual_device.0.type", "NETWORK"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.virtual_device.0.uuid", connectionTestData["pnfv"]["virtualDevice"]),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.interface.0.type", "NETWORK"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "a_side.0.access_point.0.interface.0.id", "7"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.type", "NETWORK"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_connection.test", "z_side.0.access_point.0.network.0.uuid", connectionTestData["pnfv"]["network"]),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricCreateVirtualDevice2NetworkConnectionConfig(name, virtualDeviceUuid, networkUuid string) string {
	return fmt.Sprintf(`

	resource "equinix_fabric_connection" "test" {
		type = "EVPL_VC"
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
					type = "NETWORK"
					id = 7
				}
			}
		}
		z_side {
			access_point {
				type = "NETWORK"
				network {
					uuid "%s"
				}
			}
		}
	}`, name, virtualDeviceUuid, networkUuid)
}

func CheckConnectionDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, acceptance.TestAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_connection" {
			continue
		}
		err := equinix.WaitUntilConnectionDeprovisioned(rs.Primary.ID, acceptance.TestAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}
