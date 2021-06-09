package metal

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetalProjectSSHKeyDataSource_BySearch(t *testing.T) {
	t.Parallel()

	datasourceName := "data.metal_project_ssh_key.foobar"
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")

	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectSSHKeyDataSourceConfig_bySearch(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
			},
			{
				Config:      testAccMetalProjectSSHKeyDataSourceConfig_noKey(keyName, publicKeyMaterial),
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

func TestAccMetalProjectSSHKeyDataSource_ByID(t *testing.T) {
	t.Parallel()

	datasourceName := "data.metal_project_ssh_key.foobar"

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckMetalSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalProjectSSHKeyDataSourceConfig_byID(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      testAccMetalProjectSSHKeyDataSourceConfig_noKey(keyName, publicKeyMaterial),
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

func testAccMetalProjectSSHKeyDataSourceConfig_bySearch(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "metal_project" "test" {
    name = "%s"
}

resource "metal_project_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
	project_id = metal_project.test.id
}

data "metal_project_ssh_key" "foobar" {
	search = metal_project_ssh_key.foobar.name
	project_id = metal_project.test.id
}
`, keyName, keyName, publicSshKey)

	return config
}

func testAccMetalProjectSSHKeyDataSourceConfig_noKey(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "metal_project" "test" {
    name = "%s"
}

data "metal_project_ssh_key" "foobar" {
	search = "%s"
	project_id = metal_project.test.id
}`, keyName, keyName)
	return config
}

func testAccMetalProjectSSHKeyDataSourceConfig_byID(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "metal_project" "test" {
    name = "%s"
}

data "metal_project_ssh_key" "foobar" {
	depends_on = [metal_project_ssh_key.foobar]
	id = metal_project_ssh_key.foobar.id
	project_id = metal_project.test.id
}

resource "metal_project_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
	project_id = metal_project.test.id
}`, keyName, keyName, publicSshKey)

	return config
}
