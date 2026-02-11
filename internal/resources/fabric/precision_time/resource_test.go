// Package precisiontime_test for EPT resources and data sources tests
package precisiontime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	testing_helpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricCreatePort2EPT_NTPConfiguration_PFCR(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var portUuid string
	if len(ports) > 0 {
		portUuid = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkEptServiceDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2EPTNPTConfig("Equinix Precision Time NTP UAT Global", "port2eptntp_PFCR", portUuid, "SV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time_service.ntp", "id"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time_service.ntp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time_service.ntp", "project.project_id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "name", "tf_acc_eptntp_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "type", "NTP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "package.code", "NTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ntp", "ipv4.default_gateway", "192.168.254.254"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_service.ntp", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_service.ntp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_service.ntp", "project.project_id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "name", "tf_acc_eptntp_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "type", "NTP"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "package.code", "NTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ntp", "ipv4.default_gateway", "192.168.254.254"),
				),
				ExpectNonEmptyPlan: true,
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
					vlan_tag= 101
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

	resource "equinix_fabric_precision_time_service" "ntp" {
	  type = "NTP"
	  name = "tf_acc_eptntp_PFCR"
	  package = {
		code = "NTP_STANDARD"
	  }
	  connections = [
		{
			uuid = equinix_fabric_connection.test.id
	  	}
     ]
	  ipv4 = {
		primary = "192.168.254.241"
		secondary = "192.168.254.242"
		network_mask = "255.255.255.240"
		default_gateway = "192.168.254.254"
	  }
	}

	data "equinix_fabric_precision_time_service" "ntp" {
	  ept_service_id = equinix_fabric_precision_time_service.ntp.id
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
		PreCheck:                 func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             checkEptServiceDelete,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricCreatePort2EPTPTPConfig("Equinix Precision Time PTP Global UAT", "port2eptptp_PFCR", portUuid, "SV"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time_service.ptp", "uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time_service.ptp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("equinix_fabric_precision_time_service.ptp", "project.project_id"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "name", "tf_acc_eptptp_PFCR"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "type", "PTP"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "package.code", "PTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"equinix_fabric_precision_time_service.ptp", "ipv4.default_gateway", "192.168.254.254"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_service.ptp", "uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_service.ptp", "connections.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_service.ptp", "project.project_id"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "name", "tf_acc_eptptp_PFCR"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "type", "PTP"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "package.code", "PTP_STANDARD"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "ipv4.primary", "192.168.254.241"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "ipv4.secondary", "192.168.254.242"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "ipv4.network_mask", "255.255.255.240"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_precision_time_service.ptp", "ipv4.default_gateway", "192.168.254.254"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_services.all", "data.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_services.all", "data.0.connections.0.uuid"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_services.all", "data.0.project.project_id"),
					resource.TestCheckResourceAttrSet("data.equinix_fabric_precision_time_services.all", "data.0.href"),
					resource.TestCheckResourceAttr("data.equinix_fabric_precision_time_services.all", "data.#", "1"),
					resource.TestCheckResourceAttr("data.equinix_fabric_precision_time_services.all", "pagination.%", "5"),
					resource.TestCheckResourceAttr("data.equinix_fabric_precision_time_services.all", "pagination.limit", "2"),
					resource.TestCheckResourceAttr("data.equinix_fabric_precision_time_services.all", "pagination.offset", "1"),
				),
				ExpectNonEmptyPlan: true,
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
					vlan_tag= "100"
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

	resource "equinix_fabric_precision_time_service" "ptp" {
	  type = "PTP"
	  name = "tf_acc_eptptp_PFCR"
	  package = {
		code = "PTP_STANDARD"
	  }
	  connections = [
		{
			uuid = equinix_fabric_connection.test.id
	  	}
     ]
	  ipv4 = {
		primary = "192.168.254.241"
		secondary = "192.168.254.242"
		network_mask = "255.255.255.240"
		default_gateway = "192.168.254.254"
	  }
	}

	data "equinix_fabric_precision_time_service" "ptp" {
	  ept_service_id = equinix_fabric_precision_time_service.ptp.id
	}

	data "equinix_fabric_precision_time_services" "all" {
		depends_on = [equinix_fabric_precision_time_service.ptp, equinix_fabric_connection.test]
		  pagination = {
			limit = 2
			offset = 1
		  }
		  filters = [{
			property = "/type"
			operator = "="
			values = ["PTP"]
		  }]
		  sort = [{
			direction = "DESC"
			property = "/uuid"
		  }]
	}
`, spName, name, portUuid, zSideMetro)
}

func checkEptServiceDelete(s *terraform.State) error {
	ctx := context.Background()
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewFabricClientForTesting(ctx)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_fabric_precision_time_service" {
			continue
		}

		if eptService, _, err := client.PrecisionTimeApi.GetTimeServicesById(ctx, rs.Primary.ID).Execute(); err == nil {
			if eptService.GetState() == fabricv4.PRECISIONTIMESERVICERESPONSESTATE_PROVISIONED {
				return fmt.Errorf("fabric EPT service %s still exists and is %s",
					rs.Primary.ID, eptService.GetState())
			}
		}
	}
	return nil
}
