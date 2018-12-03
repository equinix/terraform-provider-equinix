package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketVolume_Basic(t *testing.T) {
	var volume packngo.Volume

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "plan", "storage_1"),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "billing_cycle", "hourly"),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "size", "100"),
				),
			},
		},
	})
}

func TestAccPacketVolume_Update(t *testing.T) {
	var volume, v1, v2, v3, v4 packngo.Volume

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 10, "descstr", "storage_1", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "locked", "true"),
				),
			},
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 10, "descstr", "storage_1", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &v1),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "locked", "false"),
					testAccCheckPacketSameVolume(t, &volume, &v1),
				),
			},
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 10, "descstr2", "storage_2", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &v2),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "description", "descstr2"),
					testAccCheckPacketSameVolume(t, &volume, &v2),
				),
			},
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 20, "descstr2", "storage_2", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &v3),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "size", "20"),
					testAccCheckPacketSameVolume(t, &volume, &v3),
				),
			},
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 22, "descstr2", "storage_2", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &v4),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "locked", "true"),
					testAccCheckPacketSameVolume(t, &volume, &v4),
				),
			},
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 25, "descstr2", "storage_2", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketVolumeExists("packet_volume.foobar", &v4),
					resource.TestCheckResourceAttr(
						"packet_volume.foobar", "locked", "false"),
					testAccCheckPacketSameVolume(t, &volume, &v4),
				),
			},
		},
	})
}

func testAccCheckPacketVolumeDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_volume" {
			continue
		}
		if _, _, err := client.Volumes.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckPacketSameVolume(t *testing.T, before, after *packngo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected volume to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func TestAccPacketVolume_importBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketVolumeConfig_var(rs, 10, "descstr", "storage_1", false),
			},
			resource.TestStep{
				ResourceName:      "packet_volume.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPacketVolumeExists(n string, volume *packngo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundVolume, _, err := client.Volumes.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundVolume.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundVolume)
		}

		*volume = *foundVolume

		return nil
	}
}

func testAccCheckPacketVolumeConfig_basic(projectSuffix string) string {
	return fmt.Sprintf(
		`resource "packet_project" "foobar" {
    name = "%s"
}

resource "packet_volume" "foobar" {
    plan = "storage_1"
    billing_cycle = "hourly"
    size = 100
    project_id = "${packet_project.foobar.id}"
    facility = "ewr1"
    snapshot_policies = { snapshot_frequency = "1day", snapshot_count = 7 }
}`, projectSuffix)
}

func testAccCheckPacketVolumeConfig_var(projSuffix string, size int, desc string, planID string, locked bool) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "TerraformTestProject-%s"
}

resource "packet_volume" "foobar" {
    billing_cycle = "hourly"
    size = %d
    description = "%s"
    project_id = "${packet_project.foobar.id}"
    facility = "ewr1"
    plan = "%s"
    locked = %t
}
`, projSuffix, size, desc, planID, locked)
}
