package project_ssh_key_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccDataSourceMetalProjectSSHKey_bySearch(t *testing.T) {
	datasourceName := "data.equinix_metal_project_ssh_key.foobar"
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccMetalProjectSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalProjectSSHKeyConfig_bySearch(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
			},
			{
				Config:      testAccDataSourceMetalProjectSSHKeyConfig_noKey(keyName, publicKeyMaterial),
				ExpectError: regexp.MustCompile("was not found"),
			},
			{
				// Exit the tests with an empty state and a valid config
				// following the previous error config. This is needed for the
				// destroy step to succeed.
				Config: `/* this config intentionally left blank */`,
			},
		},
	})
}

func TestAccDataSourceMetalProjectSSHKeyDataSource_byID(t *testing.T) {
	datasourceName := "data.equinix_metal_project_ssh_key.foobar"

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories:  acceptance.ProtoV5ProviderFactories,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccMetalProjectSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMetalProjectSSHKeyDataSourceConfig_byID(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
				// Why was follwing flag set? The plan is applied and then it's empty.
				// It's causing errors in acceptance tests. Was this because of some API bug?
				// ExpectNonEmptyPlan: true,
			},
			{
				Config:      testAccDataSourceMetalProjectSSHKeyConfig_noKey(keyName, publicKeyMaterial),
				ExpectError: regexp.MustCompile("was not found"),
			},
			{
				// Exit the tests with an empty state and a valid config
				// following the previous error config. This is needed for the
				// destroy step to succeed.
				Config: `/* this config intentionally left blank */`,
			},
		},
	})
}

func testAccDataSourceMetalProjectSSHKeyConfig_bySearch(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "%s"
}

resource "equinix_metal_project_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
	project_id = equinix_metal_project.test.id
}

data "equinix_metal_project_ssh_key" "foobar" {
	search = equinix_metal_project_ssh_key.foobar.name
	project_id = equinix_metal_project.test.id
}
`, keyName, keyName, publicSshKey)

	return config
}

func testAccDataSourceMetalProjectSSHKeyConfig_noKey(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "%s"
}

data "equinix_metal_project_ssh_key" "foobar" {
	search = "%s"
	project_id = equinix_metal_project.test.id
}`, keyName, keyName)
	return config
}

func testAccDataSourceMetalProjectSSHKeyDataSourceConfig_byID(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "%s"
}

data "equinix_metal_project_ssh_key" "foobar" {
	depends_on = [equinix_metal_project_ssh_key.foobar]
	id = equinix_metal_project_ssh_key.foobar.id
	project_id = equinix_metal_project.test.id
}

resource "equinix_metal_project_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
	project_id = equinix_metal_project.test.id
}`, keyName, keyName, publicSshKey)

	return config
}

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccDataSourceMetalProjectSSHKey_upgradeFromVersion(t *testing.T) {
	datasourceName := "data.equinix_metal_project_ssh_key.foobar"
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { acceptance.TestAccPreCheckMetal(t) },
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccMetalProjectSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.27.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: testAccDataSourceMetalProjectSSHKeyConfig_bySearch(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
			},
			{
				ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
				Config:                   testAccDataSourceMetalProjectSSHKeyConfig_bySearch(keyName, publicKeyMaterial),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}
