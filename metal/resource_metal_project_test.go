package metal

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("metal_project", &resource.Sweeper{
		Name:         "metal_project",
		Dependencies: []string{"metal_device"},
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
		return fmt.Errorf("Error getting project list for sweeping projects: %s", err)
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

func TestAccMetalProject_Basic(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
		},
	})
}

func TestAccMetalProject_BGPBasic(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectConfig_BGP(rInt, "2SFsdfsg43"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
				),
			},
		},
	})
}

func TestAccMetalProject_BackendTransferUpdate(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "backend_transfer", "false"),
				),
			},
			{
				Config: testAccCheckMetalProjectConfig_BT(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "backend_transfer", "true"),
				),
			},
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "backend_transfer", "false"),
				),
			},
		},
	})
}

func TestAccMetalProject_Update(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt + 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt+1)),
				),
			},
		},
	})
}

func testAccCheckMetalSameProject(t *testing.T, before, after *packngo.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.ID != after.ID {
			t.Fatalf("Expected device to be the same, but it was recreated: %s -> %s", before.ID, after.ID)
		}
		return nil
	}
}

func TestAccMetalProject_BGPUpdate(t *testing.T) {
	var p1, p2, p3 packngo.Project
	rInt := acctest.RandInt()
	res := "metal_project.foobar"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists(res, &p1),
					resource.TestCheckResourceAttr(res, "name",
						fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
			{
				Config: testAccCheckMetalProjectConfig_BGP(rInt, "fdsfsdf432F"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists(res, &p2),
					resource.TestCheckResourceAttr(res, "bgp_config.0.md5", "fdsfsdf432F"),
					testAccCheckMetalSameProject(t, &p1, &p2),
				),
			},
			{
				Config: testAccCheckMetalProjectConfig_BGP(rInt, "fdsfsdf432G"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists(res, &p3),
					resource.TestCheckResourceAttr(res, "bgp_config.0.md5", "fdsfsdf432G"),
					testAccCheckMetalSameProject(t, &p2, &p3),
				),
			},
			{
				Config:      testAccCheckMetalProjectConfig_basic(rInt),
				ExpectError: regexp.MustCompile("can not be removed"),
			},
		},
	})
}

func testAccCheckMetalProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_project" {
			continue
		}
		if _, _, err := client.Projects.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Project still exists")
		}
	}

	return nil
}

func testAccCheckMetalProjectExists(n string, project *packngo.Project) resource.TestCheckFunc {
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

func testAccCheckMetalProjectConfig_BT(r int) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-project-%d"
	backend_transfer = true
}`, r)
}

func testAccCheckMetalProjectConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-project-%d"
}`, r)
}

func testAccCheckMetalProjectConfig_BGP(r int, pass string) string {
	return fmt.Sprintf(`
resource "metal_project" "foobar" {
    name = "tfacc-project-%d"
	bgp_config {
		deployment_type = "local"
		md5 = "%s"
		asn = 65000
	}
}`, r, pass)
}

func testAccCheckMetalProjectOrgConfig(r string) string {
	return fmt.Sprintf(`
resource "metal_organization" "test" {
	name = "tfacc-project-%s"
}

resource "metal_project" "foobar" {
		name = "tfacc-project-%s"
		organization_id = "${metal_organization.test.id}"
}`, r, r)
}

func TestAccMetalProjectOrg(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectOrgConfig(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalProjectExists("metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%s", rn)),
				),
			},
		},
	})
}

func TestAccMetalProject_importBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalProjectConfig_basic(rInt),
			},
			{
				ResourceName:      "metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
