package equinix_test

import (
	"encoding/json"
	"fmt"
	v4 "github.com/equinix-labs/fabric-go/fabric/v4"
	"os"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	FabricDedicatedPortEnvVar = "TF_ACC_FABRIC_DEDICATED_PORTS"
)

type EnvPorts map[string]map[string][]v4.Port

func GetFabricEnvPorts(t *testing.T) EnvPorts {
	var ports EnvPorts
	portJson := os.Getenv(FabricDedicatedPortEnvVar)
	if err := json.Unmarshal([]byte(portJson), &ports); err != nil {
		t.Fatalf("Failed reading port data from environment: %v, %s", err, portJson)
	}
	return ports
}

func TestAccDataSourceFabricPort_PNFV(t *testing.T) {
	port := GetFabricEnvPorts(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPort(port["pnfv"]["dot1q"][0].Uuid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "id", port["pnfv"]["dot1q"][0].Uuid),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "name", port["pnfv"]["dot1q"][0].Name),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "bandwidth"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_port.test", "used_bandwidth"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "type", string(*port["pnfv"]["dot1q"][0].Type_)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "encapsulation.0.type", port["pnfv"]["dot1q"][0].Encapsulation.Type_),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "state", string(*port["pnfv"]["dot1q"][0].State)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "redundancy.0.priority", string(*port["pnfv"]["dot1q"][0].Redundancy.Priority)),
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
	port := GetFabricEnvPorts(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPorts(port["pnfv"]["dot1q"][0].Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "id", port["pnfv"]["dot1q"][0].Uuid),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.name", port["pnfv"]["dot1q"][0].Name),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.bandwidth"),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.used_bandwidth"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.type", string(*port["pnfv"]["dot1q"][0].Type_)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.state", string(*port["pnfv"]["dot1q"][0].State)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.encapsulation.0.type", port["pnfv"]["dot1q"][0].Encapsulation.Type_),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.redundancy.0.priority", string(*port["pnfv"]["dot1q"][0].Redundancy.Priority)),
					resource.TestCheckResourceAttrSet(
						"data.equinix_fabric_ports.test", "data.0.lag_enabled"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.type", string(*port["pnfv"]["dot1q"][0].Type_)),
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
