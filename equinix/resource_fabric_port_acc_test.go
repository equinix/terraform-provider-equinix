package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricReadPort(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadPortConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "name", "ops-user100-CX-SV1-NL-Qinq-STD-1G-PRI-NK-349"),
				),
			},
		},
	})
}

func testAccFabricReadPortConfig() string {
	return `data "equinix_fabric_port" "test" {
	uuid = "c4d9350e-783c-83cd-1ce0-306a5c00a600"
	}`
}

// Get Ports By Name
func testAccFabricReadGetPortsByNameConfig(name string) string {
	return fmt.Sprintf(`data "equinix_fabric_ports" "test" {
	filters {
		name = "%s"
		}
	}
`, name)
}

func TestAccFabricGetPortsByName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadGetPortsByNameConfig("ops-user100-CX-DC11-NL-Dot1q-BO-10G-SEC-JP-113"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.#", fmt.Sprint(1)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.name", "ops-user100-CX-DC11-NL-Dot1q-BO-10G-SEC-JP-113"),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.uuid", "c4d9350e-7791-791d-1ce0-306a5c00a600"),
					resource.TestCheckNoResourceAttr(
						"data.equinix_fabric_ports.test", "pagination"),
				),
			},
		},
	})
}
