package packet

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPacketProjectSSHKeyDataSource_BySearch(t *testing.T) {
	t.Parallel()

	datasourceName := "data.packet_project_ssh_key.foobar"
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")

	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckPacketSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPacketProjectSSHKeyDataSourceConfig_bySearch(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
			},
			{
				Config:      testAccPacketProjectSSHKeyDataSourceConfig_noKey(keyName, publicKeyMaterial),
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

func TestAccPacketProjectSSHKeyDataSource_ByID(t *testing.T) {
	t.Parallel()

	datasourceName := "data.packet_project_ssh_key.foobar"

	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	keyName := acctest.RandomWithPrefix("tfacc-project-key")

	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckPacketSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPacketProjectSSHKeyDataSourceConfig_byID(keyName, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						datasourceName, "name", keyName),
					resource.TestCheckResourceAttr(
						datasourceName, "public_key", publicKeyMaterial),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      testAccPacketProjectSSHKeyDataSourceConfig_noKey(keyName, publicKeyMaterial),
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

func testAccPacketProjectSSHKeyDataSourceConfig_bySearch(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "packet_project" "test" {
    name = "%s"
}

resource "packet_project_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
	project_id = packet_project.test.id
}

data "packet_project_ssh_key" "foobar" {
	search = packet_project_ssh_key.foobar.name
	project_id = packet_project.test.id
}
`, keyName, keyName, publicSshKey)

	return config
}

func testAccPacketProjectSSHKeyDataSourceConfig_noKey(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "packet_project" "test" {
    name = "%s"
}

data "packet_project_ssh_key" "foobar" {
	search = "%s"
	project_id = packet_project.test.id
}`, keyName, keyName)
	return config
}

func testAccPacketProjectSSHKeyDataSourceConfig_byID(keyName, publicSshKey string) string {
	config := fmt.Sprintf(`
resource "packet_project" "test" {
    name = "%s"
}

data "packet_project_ssh_key" "foobar" {
	depends_on = [packet_project_ssh_key.foobar]
	id = packet_project_ssh_key.foobar.id
	project_id = packet_project.test.id
}

resource "packet_project_ssh_key" "foobar" {
	name = "%s"
	public_key = "%s"
	project_id = packet_project.test.id
}`, keyName, keyName, publicSshKey)

	return config
}
