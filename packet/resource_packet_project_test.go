package packet

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("packet_project", &resource.Sweeper{
		Name:         "packet_project",
		Dependencies: []string{"packet_device"},
		F:            testSweepProjects,
	})
}

func testSweepProjects(region string) error {
	log.Printf("[DEBUG] Sweeping projects")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping projects: %s", err)
	}
	client := meta.(*packngo.Client)

	ps, _, err := client.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("Error getting project list for sweepeing projects: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if strings.HasPrefix(p.Name, "tfacc-") {
			pids = append(pids, p.ID)
		}
	}
	for _, pid := range pids {
		log.Printf("Removing project %s", pid)
		_, err := client.Projects.Delete(pid)
		if err != nil {
			return fmt.Errorf("Error deleting project %s", err)
		}
	}
	return nil
}

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
						"packet_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
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
			{
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
						"packet_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
			{
				Config: testAccCheckPacketProjectConfig_basic(rInt + 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketProjectExists("packet_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"packet_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt+1)),
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
						fmt.Sprintf("tfacc-project-%d", rInt)),
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
    name = "tfacc-project-%d"
	backend_transfer = true
}`, r)
}

func testAccCheckPacketProjectConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "tfacc-project-%d"
}`, r)
}

func testAccCheckPacketProjectConfig_BGP(r int, pass string) string {
	return fmt.Sprintf(`
resource "packet_project" "foobar" {
    name = "tfacc-project-%d"
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
	name = "tfacc-project-%s"
}

resource "packet_project" "foobar" {
		name = "tfacc-project-%s"
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
						"packet_project.foobar", "name", fmt.Sprintf("tfacc-project-%s", rn)),
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
