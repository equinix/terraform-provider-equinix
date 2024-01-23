package equinix_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/packethost/packngo"
)

func init() {
	resource.AddTestSweepers("equinix_metal_project", &resource.Sweeper{
		Name:         "equinix_metal_project",
		Dependencies: []string{"equinix_metal_vlan"},
		F:            testSweepProjects,
	})
}

func testSweepProjects(region string) error {
	log.Printf("[DEBUG] Sweeping projects")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping projects: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping projects: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if isSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	for _, pid := range pids {
		log.Printf("Removing project %s", pid)
		_, err := metal.Projects.Delete(pid)
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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

// TODO(displague) How do we test this without TF_ACC set?
func TestAccMetalProject_errorHandling(t *testing.T) {
	rInt := acctest.RandInt()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	mockAPI := httptest.NewServer(http.HandlerFunc(handler))
	mockEquinix := equinix.Provider()
	mockEquinix.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := config.Config{
			BaseURL:   mockAPI.URL,
			Token:     "fake-for-mock-test",
			AuthToken: "fake-for-mock-test",
		}
		err := config.Load(ctx)
		return &config, diag.FromErr(err)
	}

	mockProviders := map[string]*schema.Provider{
		"equinix": mockEquinix,
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

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("X-Request-Id", "needed for equinix_errors.FriendlyError")
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	mockAPI := httptest.NewServer(http.HandlerFunc(handler))
	mockEquinix := equinix.Provider()
	mockEquinix.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		config := config.Config{
			BaseURL:   mockAPI.URL,
			Token:     "fake-for-mock-test",
			AuthToken: "fake-for-mock-test",
		}
		err := config.Load(ctx)
		return &config, diag.FromErr(err)
	}

	mockProviders := map[string]*schema.Provider{
		"equinix": mockEquinix,
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

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

		client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

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
	address {
		address = "tfacc org street"
		city = "london"
		zip_code = "12345"
		country = "GB"
	}
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders: acceptance.TestExternalProviders,
		Providers:         acceptance.TestAccProviders,
		CheckDestroy:      testAccMetalProjectCheckDestroyed,
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
