package equinix

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("equinix_metal_vrf", &resource.Sweeper{
		Name: "equinix_metal_vrf",
		Dependencies: []string{
			"equinix_metal_device",
			// TODO: add sweeper when offered
			// "equinix_metal_reserved_ip_block",
		},
		F: testSweepVRFs,
	})
}

func testSweepVRFs(region string) error {
	log.Printf("[DEBUG] Sweeping VRFs")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping VRFs: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping VRFs: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if isSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.VRFs.List(pid, nil)
		if err != nil {
			log.Printf("Error listing VRFs to sweep: %s", err)
			continue
		}
		for _, d := range ds {
			if isSweepableTestResource(d.Name) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing VRFs %s", did)
		_, err := metal.VRFs.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting VRFs %s", err)
		}
	}
	return nil
}

func TestAccMetalVRF_basic(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_vrf.foobar", "local_asn"),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalVRF_withIPRanges(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withIPRanges(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalVRFConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
		},
	})
}

func TestAccMetalVRF_withIPReservations(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withIPRanges(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withIPReservations(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrPair("equinix_metal_vrf.foobar", "id", "equinix_metal_reserved_ip_block.foobar", "vrf_id"),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "equinix_metal_reserved_ip_block.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalVRF_withGateway(t *testing.T) {
	var vrf packngo.VRF
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalVRFCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVRFConfig_withIPReservations(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
				),
			},
			{
				Config: testAccMetalVRFConfig_withGateway(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalVRFExists("equinix_metal_vrf.foobar", &vrf),
					resource.TestCheckResourceAttr(
						"equinix_metal_vrf.foobar", "name", fmt.Sprintf("tfacc-vrf-%d", rInt)),
					resource.TestCheckResourceAttrPair("equinix_metal_vrf.foobar", "id", "equinix_metal_gateway.foobar", "vrf_id"),
				),
			},
			{
				ResourceName:      "equinix_metal_vrf.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "equinix_metal_gateway.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalVRFCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vrf" {
			continue
		}
		if _, _, err := client.VRFs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal VRF still exists")
		}
	}

	return nil
}

func testAccMetalVRFExists(n string, vrf *packngo.VRF) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Config).metal

		foundResource, _, err := client.VRFs.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundResource.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundResource)
		}

		*vrf = *foundResource

		return nil
	}
}

func testAccMetalVRFConfig_basic(r int) string {
	testMetro := "da"

	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vrfs-%d"
}

resource "equinix_metal_vrf" "foobar" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	project_id = "${equinix_metal_project.foobar.id}"
}`, r, r, testMetro)
}

func testAccMetalVRFConfig_withIPRanges(r int) string {
	testMetro := "da"

	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vrfs-%d"
}

resource "equinix_metal_vrf" "foobar" {
	name = "tfacc-vrf-%d"
	metro = "%s"
	description = "tfacc-vrf-%d"
	local_asn = "65000"
	ip_ranges = ["192.168.100.0/25"]
	project_id = "${equinix_metal_project.foobar.id}"
}`, r, r, testMetro, r)
}

func testAccMetalVRFConfig_withIPReservations(r int) string {
	testMetro := "da"

	return testAccMetalVRFConfig_withIPRanges(r) + fmt.Sprintf(`

resource "equinix_metal_reserved_ip_block" "foobar" {
	vrf_id = "${equinix_metal_vrf.foobar.id}"
	cidr = 29
	description = "tfacc-reserved-ip-block-%d"
	network = "192.168.100.0"
	type = "vrf"
	metro = "%s"
	project_id = "${equinix_metal_project.foobar.id}"
}
`, r, testMetro)
}

func testAccMetalVRFConfig_withGateway(r int) string {
	testMetro := "da"

	return testAccMetalVRFConfig_withIPReservations(r) + fmt.Sprintf(`

resource "equinix_metal_vlan" "foobar" {
	description = "test VLAN for VRF"
	metro       = "%s"
	project_id  = equinix_metal_project.foobar.id
}

resource "equinix_metal_gateway" "foobar" {
    project_id        = equinix_metal_project.foobar.id
    vlan_id           = equinix_metal_vlan.foobar.id
    ip_reservation_id = equinix_metal_reserved_ip_block.foobar.id
}
`, testMetro)
}
