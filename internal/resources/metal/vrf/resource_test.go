package vrf_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	metalDedicatedConnIDEnvVar = "TF_ACC_METAL_DEDICATED_CONNECTION_ID"
)

func TestAccMetalVRF_basic(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_vrf.test", "local_asn"),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalVRF_withIPRanges(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withIPRanges(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
		},
	})
}

func TestAccMetalVRF_withIPReservations(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withIPRanges(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withIPReservations(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrPair("equinix_metal_vrf.test", "id", "equinix_metal_reserved_ip_block.test", "vrf_id"),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:            "equinix_metal_reserved_ip_block.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_state"},
			},
		},
	})
}

func TestAccMetalVRF_withGateway(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withIPReservations(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withGateway(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.test", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrPair("equinix_metal_vrf.test", "id", "equinix_metal_gateway.test", "vrf_id"),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "equinix_metal_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalVRFConfig_withConnection(t *testing.T) {
	var vrf metalv1.Vrf
	rInt := acctest.RandInt()
	nniVlan := acctest.RandIntRange(1024, 1093)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV6ProviderFactories: acceptance.ProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withVC(rInt, nniVlan),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.test", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_virtual_circuit.test", "name", fmt.Sprintf("tfacc-vc-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"nni_vlan", strconv.Itoa(nniVlan)),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_virtual_circuit.test",
						"vrf_id", "equinix_metal_vrf.test", "id"),
					resource.TestCheckNoResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"vlan_id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"peer_asn", "65530"),
					resource.TestCheckResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"subnet", "192.168.100.16/31"),
					resource.TestCheckResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"metal_ip", "192.168.100.16"),
					resource.TestCheckResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"customer_ip", "192.168.100.17"),
				),
			},
			{
				ResourceName:      "equinix_metal_virtual_circuit.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalVRFConfig_withVCGateway(rInt, nniVlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(
						"equinix_metal_virtual_circuit.test",
						"vlan_id"),
				),
			},
			{
				ResourceName:            "equinix_metal_reserved_ip_block.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"wait_for_state"},
			},
			{
				ResourceName:      "equinix_metal_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "equinix_metal_gateway.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalVRFCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vrf" {
			continue
		}
		if _, _, err := client.VRFsApi.FindVrfById(context.Background(), rs.Primary.ID).Execute(); err == nil {
			return fmt.Errorf("Metal VRF still exists")
		}
	}

	return nil
}

func testAccMetalVRFExists(n string, vrf *metalv1.Vrf) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

		foundResource, _, err := client.VRFsApi.FindVrfById(context.Background(), rs.Primary.ID).Execute()
		if err != nil {
			return err
		}
		if foundResource.GetId() != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundResource)
		}

		*vrf = *foundResource

		return nil
	}
}

func testAccMetalVRFConfig_basic(r int) string {
	testMetro := "da"

	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-vrfs-%d"
}

resource "equinix_metal_vrf" "test" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	project_id = "${equinix_metal_project.test.id}"
}`, r, r, testMetro)
}

func testAccMetalVRFConfig_withIPRanges(r int) string {
	testMetro := "da"

	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-vrfs-%d"
}

resource "equinix_metal_vrf" "test" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	description = "tfacc-vrf-%d"
	local_asn = "65000"
	ip_ranges = ["192.168.100.0/25"]
	project_id = equinix_metal_project.test.id
}`, r, r, testMetro, r)
}

func testAccMetalVRFConfig_withIPReservations(r int) string {
	testMetro := "da"

	return testAccMetalVRFConfig_withIPRanges(r) + fmt.Sprintf(`

resource "equinix_metal_reserved_ip_block" "test" {
	vrf_id = equinix_metal_vrf.test.id
	cidr = 29
	description = "tfacc-reserved-ip-block-%d"
	network = "192.168.100.0"
	type = "vrf"
	metro = "%s"
	project_id = equinix_metal_project.test.id
}
`, r, testMetro)
}

func testAccMetalVRFConfig_withGateway(r int) string {
	testMetro := "da"

	return testAccMetalVRFConfig_withIPReservations(r) + fmt.Sprintf(`

resource "equinix_metal_vlan" "test" {
	description = "tfacc-vlan-vrf"
	metro       = "%s"
	project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_gateway" "test" {
    project_id        = equinix_metal_project.test.id
    vlan_id           = equinix_metal_vlan.test.id
    ip_reservation_id = equinix_metal_reserved_ip_block.test.id
}
`, testMetro)
}

func testAccMetalVRFConfig_withVC(r, nniVlan int) string {
	// Dedicated connection in DA metro
	testConnection := os.Getenv(metalDedicatedConnIDEnvVar)
	return testAccMetalVRFConfig_withIPRanges(r) + fmt.Sprintf(`

	data "equinix_metal_connection" "test" {
		connection_id = "%s"
	}

	resource "equinix_metal_virtual_circuit" "test" {
		name = "tfacc-vc-%d"
		description = "tfacc-vc-%d"
		connection_id = data.equinix_metal_connection.test.id
		project_id = equinix_metal_project.test.id
		port_id = data.equinix_metal_connection.test.ports[0].id
		nni_vlan = %d
		vrf_id = equinix_metal_vrf.test.id
		peer_asn = 65530
		subnet = "192.168.100.16/31"
		metal_ip = "192.168.100.16"
		customer_ip = "192.168.100.17"
	}
	`, testConnection, r, r, nniVlan)
}

func testAccMetalVRFConfig_withVCGateway(r, nniVlan int) string {
	// Dedicated connection in DA metro
	testConnection := os.Getenv(metalDedicatedConnIDEnvVar)
	return testAccMetalVRFConfig_withGateway(r) + fmt.Sprintf(`
	data "equinix_metal_connection" "test" {
		connection_id = "%s"
	}

	resource "equinix_metal_virtual_circuit" "test" {
		name = "tfacc-vc-%d"
		description = "tfacc-vc-%d"
		connection_id = data.equinix_metal_connection.test.id
		project_id = equinix_metal_project.test.id
		port_id = data.equinix_metal_connection.test.ports[0].id
		nni_vlan = %d
		vrf_id = equinix_metal_vrf.test.id
		peer_asn = 65530
		subnet = "192.168.100.16/31"
		metal_ip = "192.168.100.16"
		customer_ip = "192.168.100.17"
	}`, testConnection, r, r, nniVlan)
}
