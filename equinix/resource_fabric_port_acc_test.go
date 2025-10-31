package equinix_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	testinghelpers "github.com/equinix/terraform-provider-equinix/internal/fabric/testing_helpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFabricReadPort_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortUUID string
	if len(ports) > 0 {
		aSidePortUUID = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadPortConfig(aSidePortUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_port.test", "state", "ACTIVE"),
				),
			},
		},
	})
}

func testAccFabricReadPortConfig(portUUID string) string {
	return fmt.Sprintf(`data "equinix_fabric_port" "test" {
	uuid = "%s"
	}
`, portUUID)
}

// Get Ports By Name
func testAccFabricReadGetPortsByNameConfig(name string) string {
	return fmt.Sprintf(`data "equinix_fabric_ports" "test" {
	filter {
		name = "%s"
		}
	}
`, name)
}

func TestAccFabricGetPortsByName_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortName string
	if len(ports) > 0 {
		aSidePortName = ports["pfcr"]["dot1q"][0].GetName()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadGetPortsByNameConfig(aSidePortName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.#", fmt.Sprint(1)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.state", "ACTIVE"),
					resource.TestCheckNoResourceAttr(
						"data.equinix_fabric_ports.test", "pagination"),
				),
			},
		},
	})
}

// Get Ports By UUID
func testAccFabricReadGetPortsByUUIDConfig(name string) string {
	return fmt.Sprintf(`

data "equinix_fabric_ports" "test" {
	filter {
		uuid = "%s"
		}
	}
`, name)
}

func TestAccFabricGetPortsByUUID_PFCR(t *testing.T) {
	ports := testinghelpers.GetFabricEnvPorts(t)
	var aSidePortUUID string
	if len(ports) > 0 {
		aSidePortUUID = ports["pfcr"]["dot1q"][0].GetUuid()
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.TestAccPreCheck(t) },
		Providers: acceptance.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFabricReadGetPortsByUUIDConfig(aSidePortUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.#", fmt.Sprint(1)),
					resource.TestCheckResourceAttr(
						"data.equinix_fabric_ports.test", "data.0.state", "ACTIVE"),
					resource.TestCheckNoResourceAttr(
						"data.equinix_fabric_ports.test", "pagination"),
				),
			},
		},
	})
}
