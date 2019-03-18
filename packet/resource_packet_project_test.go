package packet

import (
	"fmt"
	"regexp"
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
			{
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

func TestAccPacketProject_BGPBasic(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckPacketProjectConfig_BGP(rInt, "2SFsdfsg43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
				),
			},
		},
	})
}

func TestAccPacketProject_BackendTransferUpdate(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "backend_transfer", "false"),
				),
			},
			{
				Config: testAccCheckPacketProjectConfig_BT(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "backend_transfer", "true"),
				),
			},
			{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "backend_transfer", "false"),
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
			{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("foobar-%d", rInt)),
				),
			},
			{
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

func testAccCheckPacketSameProject(t *testing.T, before, after *packngo.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected device to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func TestAccPacketProject_BGPUpdate(t *testing.T) {
	var p1, p2, p3 packngo.Project
	rInt := acctest.RandInt()
	res := "packet_project.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists(res, &p1),
					resource.TestCheckResourceAttr(res, "name",
						fmt.Sprintf("foobar-%d", rInt)),
				),
			},
			{
				Config: testAccCheckPacketProjectConfig_BGP(rInt, "fdsfsdf432F"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists(res, &p2),
					resource.TestCheckResourceAttr(res, "bgp_config.0.md5", "fdsfsdf432F"),
					testAccCheckPacketSameProject(t, &p1, &p2),
				),
			},
			{
				Config: testAccCheckPacketProjectConfig_BGP(rInt, "fdsfsdf432G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists(res, &p3),
					resource.TestCheckResourceAttr(res, "bgp_config.0.md5", "fdsfsdf432G"),
					testAccCheckPacketSameProject(t, &p2, &p3),
				),
			},
			{
				Config:      testAccCheckPacketProjectConfig_basic(rInt),
				ExpectError: regexp.MustCompile("can not be removed"),
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

func testAccCheckPacketProjectConfig_BT(r int) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "foobar-%d"
	backend_transfer = true
}`, r)
}

func testAccCheckPacketProjectConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "foobar-%d"
}`, r)
}

func testAccCheckPacketProjectConfig_BGP(r int, pass string) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "foobar-%d"
	bgp_config {
		deployment_type = "local"
		md5 = "%s"
		asn = 65000
	}
}`, r, pass)
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
			{
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
			{
				Config: testAccCheckPacketProjectConfig_basic(rInt),
			},
			{
				ResourceName:      "packet_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
