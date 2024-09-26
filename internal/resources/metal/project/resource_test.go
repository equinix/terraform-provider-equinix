package project_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/equinix"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/provider"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

func TestAccMetalProject_basic(t *testing.T) {
	var project metalv1.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
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
	providerConfig := testAccMetalProviderConfig(mockAPI.URL, "fake-for-mock-test", "fake-for-mock-test")
	projectConfig := testAccMetalProjectConfig_basic(rInt)

	resource.ParallelTest(t, resource.TestCase{
		ProtoV5ProviderFactories: mockProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      providerConfig + "\n" + projectConfig,
				ExpectError: regexp.MustCompile(`\bCould not create project: 422 Unprocessable Entity\b`),
			},
		},
	})
}

func mockProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	mockProviders := map[string]*schema.Provider{
		"equinix": equinix.Provider(),
	}
	mockFrameworkProvider := provider.CreateFrameworkProvider("version")
	mockProviderFactories := map[string]func() (tfprotov5.ProviderServer, error){
		"equinix": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				mockProviders["equinix"].GRPCProvider,
				providerserver.NewProtocol5(
					mockFrameworkProvider,
				),
			}
			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}
			return muxServer.ProviderServer(), nil
		},
	}
	return mockProviderFactories
}

func TestAccMetalProject_BGPBasic(t *testing.T) {
	var project metalv1.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectConfig_BGP(rInt, "2SFsdfsg43"),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.md5",
						"2SFsdfsg43"),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.asn",
						"65000"),
				),
			},
			{
				ResourceName:      "equinix_metal_project.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMetalProjectConfig_BGPWithoutMD5(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckNoResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.md5"),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "bgp_config.0.asn",
						"65000"),
				),
			},
		},
	})
}

func TestAccMetalProject_backendTransferUpdate(t *testing.T) {
	var project metalv1.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
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
	var project metalv1.Project
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
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

func testAccCheckMetalSameProject(t *testing.T, before, after *metalv1.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before.GetId() != after.GetId() {
			t.Fatalf("Expected project to be the same, but it was recreated: %s -> %s", before.GetId(), after.GetId())
		}
		return nil
	}
}

func TestAccMetalProject_BGPUpdate(t *testing.T) {
	var p1, p2, p3 metalv1.Project
	rInt := acctest.RandInt()
	res := "equinix_metal_project.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
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

func testAccMetalProjectExists(n string, project *metalv1.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.Config).NewMetalClientForTesting()

		foundProject, _, err := client.ProjectsApi.FindProjectById(context.Background(), rs.Primary.ID).Execute()
		if err != nil {
			return err
		}
		if foundProject.GetId() != rs.Primary.ID {
			return fmt.Errorf("Record not found: %v - %v", rs.Primary.ID, foundProject)
		}

		*project = *foundProject

		return nil
	}
}

func testAccMetalProviderConfig(
	endpoint string,
	token string,
	authToken string,
) string {
	return fmt.Sprintf(`
provider "equinix" {
    endpoint = "%s"
    token = "%s"
    auth_token = "%s"
}
`, endpoint, token, authToken)
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

func testAccMetalProjectConfig_BGPWithoutMD5(r int) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "foobar" {
    name = "tfacc-project-%d"
	bgp_config {
		deployment_type = "local"
		asn = 65000
	}
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
	var project metalv1.Project
	rn := acctest.RandStringFromCharSet(12, "abcdef0123456789")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
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
		PreCheck:                 func() { acceptance.TestAccPreCheck(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectCheckDestroyed,
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

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccMetalProject_basic_upgradeFromVersion(t *testing.T) {
	var project metalv1.Project
	rInt := acctest.RandInt()
	cfg := testAccMetalProjectConfig_basic(rInt)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t); acceptance.TestAccPreCheckProviderConfigured(t) },
		CheckDestroy: testAccMetalProjectCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.30.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccMetalProjectExists("equinix_metal_project.foobar", &project),
					resource.TestCheckResourceAttr(
						"equinix_metal_project.foobar", "name", fmt.Sprintf("tfacc-project-%d", rInt)),
				),
			},
			{
				ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
				Config:                   cfg,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
