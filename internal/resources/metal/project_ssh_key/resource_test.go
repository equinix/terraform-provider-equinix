package project_ssh_key_test

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccMetalProjectSSHKeyConfig_basic(name, publicSshKey string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-project_ssh_key-%s"
}

resource "equinix_metal_project_ssh_key" "test" {
    name = "tfacc-project-key-test"
    public_key = "%s"
    project_id = "${equinix_metal_project.test.id}"
}

resource "equinix_metal_device" "test" {
    hostname            = "tfacc-device-key-test"
    plan                = local.plan
    metro               = local.metro
    operating_system    = local.os
    billing_cycle       = "hourly"
    project_ssh_key_ids = ["${equinix_metal_project_ssh_key.test.id}"]
    project_id          = "${equinix_metal_project.test.id}"
    termination_time    = "%s"
}

`, acceptance.ConfAccMetalDevice_base(
		acceptance.Preferable_plans,
		acceptance.Preferable_metros,
		acceptance.Preferable_os),
		name,
		publicSshKey,
		acceptance.TestDeviceTerminationTime())
}

func TestAccMetalProjectSSHKey_basic(t *testing.T) {
	rs := acctest.RandString(10)
	var key metalv1.SSHKey
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	cfg := testAccMetalProjectSSHKeyConfig_basic(rs, publicKeyMaterial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acceptance.TestAccPreCheckMetal(t) },
		ExternalProviders:        acceptance.TestExternalProviders,
		ProtoV5ProviderFactories: acceptance.ProtoV5ProviderFactories,
		CheckDestroy:             testAccMetalProjectSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckMetalSSHKeyExists("equinix_metal_project_ssh_key.test", &key),
					resource.TestCheckResourceAttr(
						"equinix_metal_project_ssh_key.test", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "ssh_key_ids.0",
						"equinix_metal_project_ssh_key.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project.test", "id",
						"equinix_metal_project_ssh_key.test", "project_id",
					),
				),
			},
		},
	})
}

func testAccMetalProjectSSHKeyCheckDestroyed(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.Config).Metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_project_ssh_key" {
			continue
		}
		if _, _, err := client.SSHKeys.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal SSH key still exists")
		}
	}

	return nil
}

// Test to verify that switching from SDKv2 to the Framework has not affected provider's behavior
// TODO (ocobles): once migrated, this test may be removed
func TestAccMetalProjectSSHKey_upgradeFromVersion(t *testing.T) {
	rs := acctest.RandString(10)
	var key metalv1.SSHKey
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	cfg := testAccMetalProjectSSHKeyConfig_basic(rs, publicKeyMaterial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.TestAccPreCheckMetal(t) },
		CheckDestroy: testAccMetalProjectSSHKeyCheckDestroyed,
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"equinix": {
						VersionConstraint: "1.24.0", // latest version with resource defined on SDKv2
						Source:            "equinix/equinix",
					},
					"random": {
						Source: "hashicorp/random",
					},
				},
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckMetalSSHKeyExists("equinix_metal_project_ssh_key.test", &key),
					resource.TestCheckResourceAttr(
						"equinix_metal_project_ssh_key.test", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_device.test", "ssh_key_ids.0",
						"equinix_metal_project_ssh_key.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_project.test", "id",
						"equinix_metal_project_ssh_key.test", "project_id",
					),
				),
			},
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"random": {
						Source: "hashicorp/random",
					},
				},
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
