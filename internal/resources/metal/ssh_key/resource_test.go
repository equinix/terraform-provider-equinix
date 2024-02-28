package ssh_key_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetalSSHKey_basic(t *testing.T) {
	var key metalv1.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalSSHKeyConfig_basic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckMetalSSHKeyExists("equinix_metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_ssh_key.foobar", "owner_id"),
				),
			},
		},
	})
}

func TestAccMetalSSHKey_projectBasic(t *testing.T) {
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_projectBasic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project.test", "id",
						"equinix_metal_project_ssh_key.foobar", "project_id",
					),
				),
			},
		},
	})
}

func TestAccMetalSSHKey_update(t *testing.T) {
	var key metalv1.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalSSHKeyConfig_basic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckMetalSSHKeyExists("equinix_metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "public_key", publicKeyMaterial),
				),
			},
			{
				Config: testAccMetalSSHKeyConfig_basic(rInt+1, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckMetalSSHKeyExists("equinix_metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt+1)),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "public_key", publicKeyMaterial),
				),
			},
		},
	})
}

func TestAccMetalSSHKey_projectImportBasic(t *testing.T) {
	sshKey, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalSSHKeyConfig_projectBasic(acctest.RandInt(), sshKey),
			},
			{
				ResourceName:      "equinix_metal_project_ssh_key.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccMetalSSHKey_importBasic(t *testing.T) {
	sshKey, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalSSHKeyConfig_basic(acctest.RandInt(), sshKey),
			},
			{
				ResourceName:            "equinix_metal_ssh_key.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"updated"},
			},
		},
	})
}

func testAccMetalSSHKeyCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_ssh_key" {
			continue
		}
		if _, _, err := client.SSHKeys.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal SSH key still exists")
		}
	}

	return nil
}

func testAccMetalSSHKeyConfig_basic(rInt int, publicSshKey string) string {
	return fmt.Sprintf(`
resource "equinix_metal_ssh_key" "foobar" {
    name = "tfacc-user-key-%d"
    public_key = "%s"
}`, rInt, publicSshKey)
}

func testAccCheckMetalSSHKeyConfig_projectBasic(rInt int, publicSshKey string) string {
	return fmt.Sprintf(`

resource "equinix_metal_project" "test" {
    name = "tfacc-project-key-test-%d"
}

resource "equinix_metal_project_ssh_key" "foobar" {
    name = "tfacc-project-key-%d"
    public_key = "%s"
	project_id = equinix_metal_project.test.id
}`, rInt, rInt, publicSshKey)
}

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccMetalSSHKey_upgradeFromVersion(t *testing.T) {
	var key metalv1.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	cfg := testAccMetalSSHKeyConfig_basic(rInt, publicKeyMaterial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t) },
		CheckDestroy: testAccMetalSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.24.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
				},
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckMetalSSHKeyExists("equinix_metal_ssh_key.foobar", &key),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "name", fmt.Sprintf("tfacc-user-key-%d", rInt)),
					resource.TestCheckResourceAttr(
						"equinix_metal_ssh_key.foobar", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrSet(
						"equinix_metal_ssh_key.foobar", "owner_id"),
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
