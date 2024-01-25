package vlan_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("equinix_metal_vlan", &resource.Sweeper{
		Name:         "equinix_metal_vlan",
		Dependencies: []string{"equinix_metal_virtual_circuit", "equinix_metal_vrf", "equinix_metal_device"},
		F:            testSweepVlans,
	})
}

func testSweepVlans(region string) error {
	log.Printf("[DEBUG] Sweeping vlans")
	config, err := acceptance.GetConfigForNonStandardMetalTest()
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping vlans: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping vlans: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if acceptance.IsSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.ProjectVirtualNetworks.List(pid, nil)
		if err != nil {
			log.Printf("Error listing vlans to sweep: %s", err)
			continue
		}
		for _, d := range ds.VirtualNetworks {
			if acceptance.IsSweepableTestResource(d.Description) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing vlan %s", did)
		_, err := metal.ProjectVirtualNetworks.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting vlan %s", err)
		}
	}
	return nil
}

func testAccCheckMetalVlanConfig_metro(projSuffix, metro, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = equinix_metal_project.foobar.id
    metro = "%s"
    description = "%s"
    vxlan = 5
}
`, projSuffix, metro, desc)
}

func TestAccMetalVlan_metro(t *testing.T) {
	rs := acctest.RandString(10)
	metro := "sv"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalVlanConfig_metro(rs, metro, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "metro", metro),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "facility", ""),
				),
			},
		},
	})
}

func TestAccMetalVlan_basic(t *testing.T) {
	var vlan metalv1.VirtualNetwork
	rs := acctest.RandString(10)
	fac := "ny5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVlanConfig_var(rs, fac, "tfacc-vlan"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalVlanExists("equinix_metal_vlan.foovlan", &vlan),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "description", "tfacc-vlan"),
					resource.TestCheckResourceAttr(
						"equinix_metal_vlan.foovlan", "facility", fac),
				),
			},
		},
	})
}

func testAccCheckMetalVlanExists(n string, vlan *metalv1.VirtualNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.Config).Metalgo

		foundVlan, _, err := client.VLANsApi.GetVirtualNetwork(context.Background(), rs.Primary.ID).Execute()
		if err != nil {
			return err
		}
		if foundVlan.GetId() != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundVlan)
		}

		*vlan = *foundVlan

		return nil
	}
}

func testAccMetalVlanCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_vlan" {
			continue
		}
		if _, _, err := client.ProjectVirtualNetworks.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Vlan still exists")
		}
	}

	return nil
}

func testAccMetalVlanConfig_var(projSuffix, facility, desc string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-vlan-%s"
}

resource "equinix_metal_vlan" "foovlan" {
    project_id = "${equinix_metal_project.foobar.id}"
    facility = "%s"
    description = "%s"
}
`, projSuffix, facility, desc)
}

func TestAccMetalVlan_importBasic(t *testing.T) {
	rs := acctest.RandString(10)
	fac := "ny5"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalVlanCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalVlanConfig_var(rs, fac, "tfacc-vlan"),
			},
			{
				ResourceName:      "equinix_metal_vlan.foovlan",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}