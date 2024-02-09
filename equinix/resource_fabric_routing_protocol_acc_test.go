package equinix_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccFabricCreateDirectRoutingProtocol_PFCR_A(t *testing.T) {
	ports := GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][1].Uuid
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: checkRoutingProtocolDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateRoutingProtocolConfig("fcr_test_PFCR", portUuid, "190.1.1.1/30", "190::1:1/126"),
				//Config: testAccFabricCreateRoutingProtocolConfig("fcr_test_PFCR", portUuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemAttr("equinix_fabric_routing_protocol.test", "anything.*", "Anything"),
					//resource.TestCheckResourceAttr(
					//	"equinix_fabric_connection.test", "name", "fcr_test_PFCR"),
					//resource.TestCheckTypeSetElemNestedAttrs("equinix_fabric_routing_protocol.test", "direct_ipv4.*", map[string]string{
					//	"equinix_iface_ip": fmt.Sprintf("190.1.1.1/30"),
					//}),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

//func testAccFabricCreateRoutingProtocolConfig(name, portUuid string) string {
//	return fmt.Sprintf(`
//
//resource "equinix_fabric_cloud_router" "this" {
//	type = "XF_ROUTER"
//	name = "Test_PFCR"
//	location{
//		metro_code  = "SV"
//	}
//	order{
//		purchase_order_number = "1-234567"
//	}
//	notifications{
//		type = "ALL"
//		emails = ["test@equinix.com", "test1@equinix.com"]
//	}
//	project{
//		project_id = "291639000636552"
//	}
//	account {
//		account_number = 201257
//	}
//	package {
//		code = "STANDARD"
//	}
//}
//
//resource "equinix_fabric_connection" "test" {
//	type = "IP_VC"
//	name = "%s"
//	notifications{
//		type = "ALL"
//		emails = ["test@equinix.com","test1@equinix.com"]
//	}
//	order {
//		purchase_order_number = "123485"
//	}
//	bandwidth = 50
//	redundancy {
//		priority= "PRIMARY"
//	}
//	a_side {
//		access_point {
//			type = "CLOUD_ROUTER"
//			router {
//				uuid = equinix_fabric_cloud_router.this.id
//			}
//		}
//	}
//	project{
//		project_id = "291639000636552"
//	}
//	z_side {
//		access_point {
//			type = "COLO"
//			port{
//				uuid = "%s"
//			}
//			link_protocol {
//				type= "DOT1Q"
//				vlan_tag= 2325
//			}
//			location {
//				metro_code = "SV"
//			}
//		}
//	}
//}`, name, portUuid)
//}

func testAccFabricCreateRoutingProtocolConfig(name, portUuid, ip4, ip6 string) string {
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
		emails = ["test@equinix.com", "test1@equinix.com"]
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

resource "equinix_fabric_connection" "this" {
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
				vlan_tag= 2334
			}
			location {
				metro_code = "SV"
			}
		}
	}
}

resource "equinix_fabric_routing_protocol" "this" {
	connection_uuid = equinix_fabric_connection.this.id
	direct_ipv4{
		equinix_iface_ip = "%s"
	}
	direct_ipv6{
		equinix_iface_ip = "%s"
	}
	type = "DIRECT"
	name = "fabric_tf_acc_test_rpDirect"
}`, name, portUuid, ip4, ip6)
}

//resource "equinix_fabric_routing_protocol" "test" {
//connection_uuid = "equinix_fabric_connection.this.uuid"
//type = "BGP"
//bgp_ipv4{
//customer_peer_ip = "%s"
//}
//bgp_ipv6{
//customer_peer_ip = "%s"
//}
//customer_asn = "100"
//}
//
//data "equinix_fabric_routing_protocol" "test" {
//	connection_uuid = equinix_fabric_routing_protocol.test.uuid
//}`, name, portUuid, ip4, ip6)
//}

//func TestAccFabricCreateDirectRoutingProtocol(t *testing.T) {
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
//		Providers:    acceptance.TestAccProviders,
//		CheckDestroy: checkRoutingProtocolDelete,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccFabricCreateRoutingProtocolDirectConfig("99d6bdc8-206f-4bff-a899-0dba708c03db", "190.1.1.1/30", "172::1:1/126"),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckTypeSetElemNestedAttrs("equinix_fabric_routing_protocol.test", "direct_ipv4.*", map[string]string{
//						"equinix_iface_ip": fmt.Sprintf("190.1.1.1/30"),
//					}),
//				),
//				ExpectNonEmptyPlan: true,
//			}, {
//				Config: testAccFabricCreateRoutingProtocolDirectConfig("99d6bdc8-206f-4bff-a899-0dba708c03db", "190.1.1.1/26", "172::1:1/126"),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckTypeSetElemNestedAttrs("equinix_fabric_routing_protocol.test", "direct_ipv4.*", map[string]string{
//						"equinix_iface_ip": fmt.Sprintf("190.1.1.1/26"),
//					}),
//				),
//				ExpectNonEmptyPlan: true,
//			},
//		},
//	})
//}

//func testAccFabricCreateRoutingProtocolDirectConfig(connectionUuid string, ipv4 string, ipv6 string) string {
//	return fmt.Sprintf(`	resource "equinix_fabric_routing_protocol" "test" {
//		connection_uuid = "%s"
//
//		type = "DIRECT"
//		name = "fabric_tf_acc_test_rpDirect"
//		direct_ipv4{
//			equinix_iface_ip = "%s"
//		}
//		direct_ipv6{
//			equinix_iface_ip = "%s"
//		}
//	}`, connectionUuid, ipv4, ipv6)
//}
//
//func TestAccFabricCreateBgpRoutingProtocol(t *testing.T) {
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:     func() { acceptance.TestAccPreCheck(t) },
//		Providers:    acceptance.TestAccProviders,
//		CheckDestroy: checkRoutingProtocolDelete,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccFabricCreateRoutingProtocolBgpConfig("99d6bdc8-206f-4bff-a899-0dba708c03db", "190.1.1.2", ""),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckTypeSetElemNestedAttrs("equinix_fabric_routing_protocol.test", "bgp_ipv4.*", map[string]string{
//						"customer_peer_ip": fmt.Sprintf("190.1.1.2"),
//					}),
//				),
//				ExpectNonEmptyPlan: true,
//			},
//			{
//				Config: testAccFabricCreateRoutingProtocolBgpConfig("99d6bdc8-206f-4bff-a899-0dba708c03db", "190.1.1.3", "172::1:2"),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckTypeSetElemNestedAttrs("equinix_fabric_routing_protocol.test", "bgp_ipv4.*", map[string]string{
//						"customer_peer_ip": fmt.Sprintf("190.1.1.3"),
//					}),
//				),
//				ExpectNonEmptyPlan: true,
//			},
//		},
//	})
//}

//func testAccFabricCreateRoutingProtocolBgpConfig(connectionUuid string, ipv4 string, ipv6 string) string {
//	return fmt.Sprintf(`	resource "equinix_fabric_routing_protocol" "test" {
//		connection_uuid = "%s"
//
//		type = "BGP"
//		bgp_ipv4{
//			customer_peer_ip = "%s"
//		}
//		bgp_ipv6{
//			customer_peer_ip = "%s"
//		}
//		customer_asn = "100"
//	}`, connectionUuid, ipv4, ipv6)
//}

func checkRoutingProtocolDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, acceptance.TestAccProvider.Meta().(*config.Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_routing_protocol" {
			continue
		}
		err := equinix.WaitUntilRoutingProtocolIsDeprovisioned(rs.Primary.ID, rs.Primary.Attributes["connection_uuid"], acceptance.TestAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

//func TestAccFabricReadRoutingProtocolByUuid(t *testing.T) {
//	resource.ParallelTest(t, resource.TestCase{
//		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
//		Providers: acceptance.TestAccProviders,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccFabricReadRoutingProtocolConfig("99d6bdc8-206f-4bff-a899-0dba708c03db", "00f48313-ab13-4524-aaad-93c31b5b8848"),
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestCheckResourceAttr(
//						"equinix_fabric_routing_protocol.test", "type", fmt.Sprint("DIRECT")),
//					resource.TestCheckResourceAttr(
//						"equinix_fabric_routing_protocol.test", "state", fmt.Sprint("PROVISIONED")),
//				),
//			},
//		},
//	})
//}
//
//func testAccFabricReadRoutingProtocolConfig(connectionUuid string, routingProtocolUuid string) string {
//	return fmt.Sprintf(`data "equinix_fabric_routing_protocol" "test" {
//	connection_uuid = "%s"
//	uuid = "%s"
//	}`, connectionUuid, routingProtocolUuid)
//}
