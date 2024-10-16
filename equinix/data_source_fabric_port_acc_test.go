package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"

	"github.com/equinix/equinix-sdk-go/services/fabricv4"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceFabricPort_PNFV(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var port fabricv4.Port
	var portType, portState, portEncapsulationType, portRedundancyPriority string
	if len(ports) > 0 {
		port = ports["pnfv"]["dot1q"][0]
		portType = string(port.GetType())
		portState = string(port.GetState())
		portEncapsulation := port.GetEncapsulation()
		portEncapsulationType = string(portEncapsulation.GetType())
		portRedundancy := port.GetRedundancy()
		portRedundancyPriority = string(portRedundancy.GetPriority())
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPort(port.GetUuid()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "id", port.GetUuid()),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "name", port.GetName()),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "bandwidth"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "used_bandwidth"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "type", portType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "encapsulation.0.type", portEncapsulationType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "state", portState),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "redundancy.0.priority", portRedundancyPriority),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "lag_enabled"),
				),
			},
		},
	})
}

func testDataSourceFabricPort(port_uuid string) string {
	return fmt.Sprintf(`
		data "equinix_fabric_port" "test" {
			uuid = "%s"
		}`,
		port_uuid)
}

func TestAccDataSourceFabricPorts_PNFV(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var port fabricv4.Port
	var portType, portState, portEncapsulationType, portRedundancyPriority string
	if len(ports) > 0 {
		port = ports["pnfv"]["dot1q"][0]
		portType = string(port.GetType())
		portState = string(port.GetState())
		portEncapsulation := port.GetEncapsulation()
		portEncapsulationType = string(portEncapsulation.GetType())
		portRedundancy := port.GetRedundancy()
		portRedundancyPriority = string(portRedundancy.GetPriority())
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPorts(port.GetName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "id", port.GetUuid()),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.name", port.GetName()),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.bandwidth"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.used_bandwidth"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.type", portType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.state", portState),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.encapsulation.0.type", portEncapsulationType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.redundancy.0.priority", portRedundancyPriority),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.lag_enabled"),
				),
			},
		},
	})
}

func testDataSourceFabricPorts(port_name string) string {
	return fmt.Sprintf(`
		data "equinix_fabric_ports" "test" {
		  filters {
			name = "%s"
		  }
		}`,
		port_name)
}

func TestAccDataSourceFabricPort_PPDS(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var port fabricv4.Port
	var portType, portState, portEncapsulationType, portRedundancyPriority string
	if len(ports) > 0 {
		port = ports["ppds"]["dot1q"][0]
		portType = string(port.GetType())
		portState = string(port.GetState())
		portEncapsulation := port.GetEncapsulation()
		portEncapsulationType = string(portEncapsulation.GetType())
		portRedundancy := port.GetRedundancy()
		portRedundancyPriority = string(portRedundancy.GetPriority())
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPort(port.GetUuid()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "id", port.GetUuid()),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "name", port.GetName()),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "bandwidth"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "used_bandwidth"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "type", portType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "encapsulation.0.type", portEncapsulationType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "state", portState),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "redundancy.0.priority", portRedundancyPriority),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "lag_enabled"),
				),
			},
		},
	})
}

func TestAccDataSourceFabricPorts_PPDS(t *testing.T) {
	ports := testing_helpers.GetFabricEnvPorts(t)
	var port fabricv4.Port
	var portType, portState, portEncapsulationType, portRedundancyPriority string
	if len(ports) > 0 {
		port = ports["ppds"]["dot1q"][0]
		portType = string(port.GetType())
		portState = string(port.GetState())
		portEncapsulation := port.GetEncapsulation()
		portEncapsulationType = string(portEncapsulation.GetType())
		portRedundancy := port.GetRedundancy()
		portRedundancyPriority = string(portRedundancy.GetPriority())
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPorts(port.GetName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "id", port.GetUuid()),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.name", port.GetName()),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.bandwidth"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.used_bandwidth"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.type", portType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.state", portState),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.encapsulation.0.type", portEncapsulationType),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.redundancy.0.priority", portRedundancyPriority),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.lag_enabled"),
				),
			},
		},
	})
}
