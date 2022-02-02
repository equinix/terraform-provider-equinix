package equinix

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("equinix_metal_project", &resource.Sweeper{
		Name:         "equinix_metal_project",
		Dependencies: []string{"equinix_metal_device"},
		F:            testSweepProjects,
	})
}

func testSweepProjects(region string) error {
	log.Printf("[DEBUG] Sweeping projects")
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("Error getting client for sweeping projects: %s", err)
	}
	client := meta.Client()

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

func TestAccMetalProject_basic(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
		},
	})
}

type mockProjectService struct {
	CreateFn          func(project *packngo.ProjectCreateRequest) (*packngo.Project, *packngo.Response, error)
	UpdateFn          func(projectID string, project *packngo.ProjectUpdateRequest) (*packngo.Project, *packngo.Response, error)
	ListFn            func(project *packngo.ListOptions) ([]packngo.Project, *packngo.Response, error)
	DeleteFn          func(projectID string) (*packngo.Response, error)
	GetFn             func(projectID string, opts *packngo.GetOptions) (*packngo.Project, *packngo.Response, error)
	ListBGPSessionsFn func(projectID string, opts *packngo.ListOptions) ([]packngo.BGPSession, *packngo.Response, error)
	ListEventsFn      func(projectID string, opts *packngo.ListOptions) ([]packngo.Event, *packngo.Response, error)
	ListSSHKeysFn     func(projectID string, opts *packngo.ListOptions) ([]packngo.SSHKey, *packngo.Response, error)
}

func (m *mockProjectService) Create(project *packngo.ProjectCreateRequest) (*packngo.Project, *packngo.Response, error) {
	return m.CreateFn(project)
}

func (m *mockProjectService) List(project *packngo.ListOptions) ([]packngo.Project, *packngo.Response, error) {
	return m.ListFn(project)
}

func (m *mockProjectService) Delete(projectID string) (*packngo.Response, error) {
	return m.DeleteFn(projectID)
}

func (m *mockProjectService) Get(projectID string, opts *packngo.GetOptions) (*packngo.Project, *packngo.Response, error) {
	return m.GetFn(projectID, opts)
}

func (m *mockProjectService) ListBGPSessions(projectID string, opts *packngo.ListOptions) ([]packngo.BGPSession, *packngo.Response, error) {
	return m.ListBGPSessionsFn(projectID, opts)
}
func (m *mockProjectService) Update(projectID string, project *packngo.ProjectUpdateRequest) (*packngo.Project, *packngo.Response, error) {
	return m.UpdateFn(projectID, project)
}

func (m *mockProjectService) ListSSHKeys(projectID string, opts *packngo.ListOptions) ([]packngo.SSHKey, *packngo.Response, error) {
	return m.ListSSHKeysFn(projectID, opts)
}

func (m *mockProjectService) ListEvents(projectID string, opts *packngo.ListOptions) ([]packngo.Event, *packngo.Response, error) {
	return m.ListEventsFn(projectID, opts)
}

var _ packngo.ProjectService = (*mockProjectService)(nil)

// TODO(displague) How do we test this without TF_ACC set?
func TestAccMetalProject_errorHandling(t *testing.T) {
	rInt := acctest.RandInt()

	mockProjectService := &mockProjectService{
		CreateFn: func(project *packngo.ProjectCreateRequest) (*packngo.Project, *packngo.Response, error) {
			httpResp := &http.Response{Status: "422 Unprocessable Entity", StatusCode: 422}
			return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
		},
	}
	mockMetal := Provider()
	mockMetal.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return &packngo.Client{Projects: mockProjectService}, nil
	}

	mockProviders := map[string]*schema.Provider{
		"equinix": mockMetal,
	}
	resource.ParallelTest(t, resource.TestCase{
		Providers: mockProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalProjectConfig_basic(rInt),
				ExpectError: regexp.MustCompile(`Error: HTTP 422`),
			},
		},
	})
}

// TODO(displague) How do we test this without TF_ACC set?
func TestAccMetalProject_apiErrorHandling(t *testing.T) {
	rInt := acctest.RandInt()

	mockProjectService := &mockProjectService{
		CreateFn: func(project *packngo.ProjectCreateRequest) (*packngo.Project, *packngo.Response, error) {
			httpResp := &http.Response{Status: "422 Unprocessable Entity", StatusCode: 422, Header: http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"12345"}}}
			return nil, &packngo.Response{Response: httpResp}, &packngo.ErrorResponse{Response: httpResp}
		},
	}
	mockMetal := Provider()
	mockMetal.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return &packngo.Client{Projects: mockProjectService}, nil
	}

	mockProviders := map[string]*schema.Provider{
		"equinix": mockMetal,
	}
	resource.ParallelTest(t, resource.TestCase{
		Providers: mockProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccMetalProjectConfig_basic(rInt),
				ExpectError: regexp.MustCompile(`Error: API Error HTTP 422`),
			},
		},
	})
}
func TestAccMetalProject_BGPBasic(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_BGP(rInt, "2SFsdfsg43"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
				),
			},
			{
				ResourceName:      "equinix_metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalProject_backendTransferUpdate(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "backend_transfer", "false"),
				),
			},
			{
				Config: testAccMetalProjectConfig_BT(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "backend_transfer", "true"),
				),
			},
			{
				ResourceName:      "equinix_metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "backend_transfer", "false"),
				),
			},
		},
	})
}

func TestAccMetalProject_update(t *testing.T) {
	var project packngo.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
			{
				Config: testAccMetalProjectConfig_basic(rInt + 1),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt+1)),
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
	res := "equinix_metal_project.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists(res, &p1),
					resource.TestCheckResourceAttr(res, "name",
						fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
			{
				Config: testAccMetalProjectConfig_BGP(rInt, "fdsfsdf432F"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists(res, &p2),
					resource.TestCheckResourceAttr(res, "bgp_config.0.md5", "fdsfsdf432F"),
					testAccCheckMetalSameProject(t, &p1, &p2),
				),
			},
			{
				ResourceName:      "equinix_metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalProjectConfig_BGP(rInt, "fdsfsdf432G"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists(res, &p3),
					resource.TestCheckResourceAttr(res, "bgp_config.0.md5", "fdsfsdf432G"),
					testAccCheckMetalSameProject(t, &p2, &p3),
				),
			},
			{
				Config:      testAccMetalProjectConfig_basic(rInt),
				ExpectError: regexp.MustCompile("can not be removed"),
			},
		},
	})
}

func testAccMetalProjectCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_project" {
			continue
		}
		if _, _, err := client.Projects.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal Project still exists")
		}
	}

	return nil
}

func testAccMetalProjectExists(n string, project *packngo.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*Config).Client()

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

func testAccMetalProjectConfig_BT(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-project-%d"
	backend_transfer = true
}`, r)
}

func testAccMetalProjectConfig_basic(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-project-%d"
}`, r)
}

func testAccMetalProjectConfig_BGP(r int, pass string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-project-%d"
	bgp_config {
		deployment_type = "local"
		md5 = "%s"
		asn = 65000
	}
}`, r, pass)
}

func testAccMetalProjectConfig_organization(r string) string {
	return fmt.Sprintf(`
resource "equinix_metal_organization" "test" {
	name = "tfacc-project-%s"
}

resource "equinix_metal_project" "foobar" {
		name = "tfacc-project-%s"
		organization_id = "${equinix_metal_organization.test.id}"
}`, r, r)
}

func TestAccMetalProject_organization(t *testing.T) {
	var project packngo.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_organization(rn),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%s", rn)),
				),
			},
			{
				ResourceName:      "equinix_metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalProject_importBasic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_basic(rInt),
			},
			{
				ResourceName:      "equinix_metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
