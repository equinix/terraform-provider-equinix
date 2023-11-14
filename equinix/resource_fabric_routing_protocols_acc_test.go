package equinix

import (
	"context"
	"fmt"
	"testing"

	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFabricCreateRoutingProtocols(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: checkRoutingProtocolsDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreateRoutingProtocolsConfig(7, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("equinix_fabric_routing_protocols.test", "direct_routing_protocol.direct_ipv4.*", map[string]string{
						"equinix_iface_ip": fmt.Sprintf("190.1.1.1/29"),
					}),
				),
				ExpectNonEmptyPlan: true,
				ImportStateId:      "<connection_uuid>/<direct_uuid>/<bgp_uuid>",
				ImportState:        true,
				ImportStateVerify:  true,
			},
		},
	})
}

func testAccFabricCreateRoutingProtocolsConfig(i, y int) string {
	return fmt.Sprintf(`resource "equinix_fabric_routing_protocols" "test" {
	  connection_uuid = "b0eb3892-6404-442b-85c0-4f8fcf53d123"
	
	  direct_routing_protocol {
		name = "Direct-%d"
		direct_ipv4 {
		  equinix_iface_ip = "190.1.1.1/29"
		}
		direct_ipv6 {
		  equinix_iface_ip = "190::1:1/126"
		}
	  }
	
	  bgp_routing_protocol {
		name = "FCR-Con-BGP-RP"
		bgp_ipv4 {
		  customer_peer_ip = "190.1.1.%d"
		  enabled          = true
		}
		bgp_ipv6 {
		  customer_peer_ip = "190::1:2"
		  enabled          = true
		}
		customer_asn = "100"
		equinix_asn  = "1345"
	  }
	}`, i, y)
}

func checkRoutingProtocolsDelete(s *terraform.State) error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, v4.ContextAccessToken, testAccProvider.Meta().(*Config).FabricAuthToken)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_routing_protocols" {
			continue
		}
		err := waitUntilRoutingProtocolIsDeprovisioned(rs.Primary.ID, rs.Primary.Attributes["connection_uuid"], testAccProvider.Meta(), ctx)
		if err != nil {
			return fmt.Errorf("API call failed while waiting for resource deletion")
		}
	}
	return nil
}

func TestAccFabricReadRoutingProtocolsByUuid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadRoutingProtocolsConfig("d7ab9902-4d89-4fc0-9a5f-9faea39d82a3", "6534071e-6216-4943-b8c4-1d2f2f0cbb07", "c1c052d0-8da9-4918-94fc-3e507f3ebf6a"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_fabric_routing_protocols.test", "direct_routing_protocol.type", fmt.Sprint("DIRECT")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_routing_protocols.test", "direct_routing_protocol.state", fmt.Sprint("PROVISIONED")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_routing_protocols.test", "bgp_routing_protocol.type", fmt.Sprint("BGP")),
					resource.TestCheckResourceAttr(
						"equinix_fabric_routing_protocols.test", "bgp_routing_protocol.state", fmt.Sprint("PROVISIONED")),
				),
			},
		},
	})
}

func testAccFabricReadRoutingProtocolsConfig(connectionUuid, directRoutingProtocolUUID, bgpRoutingProtocolUUID string) string {
	return fmt.Sprintf(`data "equinix_fabric_routing_protocols" "test" {
	connection_uuid = "%s"
	direct_routing_protocol_uuid = "%s"
	bgp_routing_protocol_uuid = "%s"
	}`, connectionUuid, directRoutingProtocolUUID, bgpRoutingProtocolUUID)
}
