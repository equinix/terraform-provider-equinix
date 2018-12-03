package packet

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func TestAccPacketProject_Basic(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("foobar-%d", rInt)),
				),
			},
		},
	})
}

func TestAccPacketProject_Update(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("foobar-%d", rInt)),
				),
			},
			resource.TestStep{
				Config: testAccCheckPacketProjectConfig_basic(rInt + 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("foobar-%d", rInt+1)),
				),
			},
		},
	})
}
func testAccCheckPacketProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_project" {
			continue
		}
		if _, _, err := client.Projects.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Project still exists")
		}
	}

	return nil
}

func testAccCheckPacketProjectExists(n string, project *packngo.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*packngo.Client)

		foundProject, _, err := client.Projects.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if foundProject.ID != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundProject)
		}

		*project = *foundProject

		return nil
	}
}

func testAccCheckPacketProjectConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "foobar-%d"
}`, r)
}

func testAccCheckPacketProjectOrgConfig(r string) string {
	return fmt.Sprintf(`
resource "packet_organization" "test" {
	name = "foobar-%s"
}

resource "packet_project" "foobar" {
		name = "foobar-%s"
		organization_id = "${packet_organization.test.id}"
}`, r, r)
}

func TestAccPacketProjectOrg(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketProjectOrgConfig(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("foobar-%s", rn)),
				),
			},
		},
	})
}

func TestAccPacketProject_importBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
			},
			resource.TestStep{
				ResourceName:      "packet_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
