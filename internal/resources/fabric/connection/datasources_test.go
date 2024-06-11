package connection_test

import (
	"fmt"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

func TestAccFabricDataSourceConnection_PFCR(t *testing.T) {
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
				Config: testAccFabricDataSourceConnectionConfig(50, aSidePortUuid, zSidePortUuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connection.test", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "name", "ds_con_test_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "order.0.purchase_order_number", "1-129105284100"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", "2444"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "a_side.0.access_point.0.location.0.metro_code", "DC"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.vlan_tag", "2555"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.location.0.metro_code", "SV"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_connections.connections", "id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.name", "ds_con_test_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.bandwidth", "50"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.type", "EVPL_VC"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.redundancy.0.priority", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.order.0.purchase_order_number", "1-129105284100"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.a_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.a_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.a_side.0.access_point.0.link_protocol.0.vlan_tag", "2444"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.a_side.0.access_point.0.location.0.metro_code", "DC"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.link_protocol.0.vlan_tag", "2555"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.location.0.metro_code", "SV"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricDataSourceConnectionConfig(bandwidth int32, aSidePortUuid, zSidePortUuid string) string {
	return fmt.Sprintf(`

resource "equinix_fabric_connection" "test" {
	type = "EVPL_VC"
	name = "ds_con_test_PFCR"
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
				vlan_tag= 2444
			}
		}
	}
	z_side {
		access_point {
			type = "COLO"
			port {
			 	uuid = "%s"
			}
			link_protocol {
				type= "DOT1Q"
				vlan_tag= 2555
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
		values = ["ds_con_test_PFCR"]
	}
	filter {
		property = "/uuid"
		operator = "="
		values = [equinix_fabric_connection.test.id]
	}
}

`, bandwidth, aSidePortUuid, zSidePortUuid)
}
