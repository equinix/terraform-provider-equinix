package connection_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricDataSourceConnection_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortUUID, zSidePortUUID string
	if len(ports) > 0 {
		aSidePortUUID = ports["pfcr"]["dot1q"][0].GetUuid()
		zSidePortUUID = ports["pfcr"]["dot1q"][1].GetUuid()
	}
	vlanID := generateUniqueVlanId()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers:    acceptance.TestAccProviders,
		CheckDestroy: CheckConnectionDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricDataSourceConnectionConfig(50, aSidePortUUID, zSidePortUUID, vlanID, vlanID),
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
						"data.equinix_fabric_connection.test", "a_side.0.access_point.0.link_protocol.0.vlan_tag", vlanID),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "a_side.0.access_point.0.location.0.metro_code", "DC"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connection.test", "z_side.0.access_point.0.link_protocol.0.vlan_tag", vlanID),
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
						"data.equinix_fabric_connections.connections", "data.0.a_side.0.access_point.0.link_protocol.0.vlan_tag", vlanID),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.a_side.0.access_point.0.location.0.metro_code", "DC"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.type", "COLO"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.link_protocol.0.type", "DOT1Q"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.link_protocol.0.vlan_tag", vlanID),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_connections.connections", "data.0.z_side.0.access_point.0.location.0.metro_code", "SV"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccFabricDataSourceConnectionConfig(bandwidth int32, aSidePortUUID, aSideVlanID, zSidePortUUID, zSideVlanID string) string {
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
				vlan_tag= "%s"
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
				vlan_tag= "%s"
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

`, bandwidth, aSidePortUUID, zSidePortUUID, aSideVlanID, zSideVlanID)
}

func generateUniqueVlanId() string {
	timestampComponent := int(time.Now().UnixNano() % 4000)
	vlanId := 1 + (timestampComponent % 4092)
	if vlanId < 1 {
		vlanId = 1
	} else if vlanId > 4092 {
		vlanId = 4092
	}
	return strconv.Itoa(vlanId)
}
