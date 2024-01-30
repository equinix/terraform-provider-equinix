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
						"data.equinix_fabric_port.test", "name", port["pnfv"]["dot1q"][0].Name),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "type", string(*port["pnfv"]["dot1q"][0].Type_)),
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
