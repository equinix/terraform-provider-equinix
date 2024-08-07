package precision_time

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricCreatePort2EPT_NTPConfiguration_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2EPTNPTConfig("NTP Configured Precision Time Service", "port2eptntp_PFCR", portUuid, "SV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time.ntp", "id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time.ntp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time.ntp", "project_id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "name", "tf_acc_eptntp_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "type", "NTP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "package.code", "PTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ntp", "ipv4.default_gateway", "192.168.254.254"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time.ntp", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time.ntp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time.ntp", "project_id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "name", "tf_acc_eptntp_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "type", "NTP"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "package.code", "PTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ntp", "ipv4.default_gateway", "192.168.254.254"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricCreatePort2EPTNPTConfig(spName, name, portUuid, zSideMetro string) string {
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
		bandwidth = 1
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
					vlan_tag= "1354"
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
	}

	resource "equinix_fabric_precision_time" "ntp" {
	  type = "NTP"
	  name = "tf_acc_eptntp_PFCR"
	  package {
		code = "NTP_STANDARD"
	  }
	  connections {
		uuid = equinix_fabric_connection.test.id
	  }
	  ipv4 {
		primary = "192.168.254.241"
		secondary = "192.168.254.242"
		network_mask = "255.255.255.240"
		default_gateway = "192.168.254.254"
	  }
	}

	data "equinix_fabric_precision_time" "ntp" {
	  uuid = equinix_precision_time.ntp.id
	}

`, spName, name, portUuid, zSideMetro)
}

func TestAccFabricCreatePort2EPT_PTPConfiguration_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2EPTPTPConfig("PTP Configured Precision Time Service", "port2eptptp_PFCR", portUuid, "SV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time.ptp", "id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time.ptp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time.ptp", "project_id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "name", "tf_acc_eptptp_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "type", "PTP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "package.code", "PTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time.ptp", "ipv4.default_gateway", "192.168.254.254"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time.ptp", "id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time.ptp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time.ptp", "project_id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "name", "tf_acc_eptptp_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "type", "PTP"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "package.code", "PTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time.ptp", "ipv4.default_gateway", "192.168.254.254"),
				),
				ExpectNonEmptyPlan: false,
			},
		},
	})

}

func testAccFabricCreatePort2EPTPTPConfig(spName, name, portUuid, zSideMetro string) string {
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
		bandwidth = 5
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
					vlan_tag= "1355"
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
	}

	resource "equinix_fabric_precision_time" "ptp" {
	  type = "PTP"
	  name = "tf_acc_eptptp_PFCR"
	  package {
		code = "PTP_STANDARD"
	  }
	  connections {
		uuid = equinix_fabric_connection.test.id
	  }
	  ipv4 {
		primary = "192.168.254.241"
		secondary = "192.168.254.242"
		network_mask = "255.255.255.240"
		default_gateway = "192.168.254.254"
	  }
	}

	data "equinix_fabric_precision_time" "ptp" {
	  uuid = equinix_precision_time.ptp.id
	}

`, spName, name, portUuid, zSideMetro)
}
