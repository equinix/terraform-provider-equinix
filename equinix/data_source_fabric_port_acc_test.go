package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceFabricPort_PNFV(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPort("c4d85dbe-fa99-a999-f7e0-306a5c00af26"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "name", "eqx-nfv-201091-CX-SV1-L-Dot1q-BO-20G-PRI-JP-395"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "type", "XF_PORT"),
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
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceFabricPorts("eqx-nfv-201091-CX-SV1-L-Dot1q-BO-20G-PRI-JP-395"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "id", "c4d85dbe-fa99-a999-f7e0-306a5c00af26"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.type", "XF_PORT"),
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
