package metal

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/packethost/packngo"
)

func metalProjectSSHKeyConfig_Basic(name, publicSshKey string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-project_ssh_key-%s"
}

resource "metal_project_ssh_key" "test" {
    name = "tfacc-project-key-test"
    public_key = "%s"
    project_id = "${metal_project.test.id}"
}

resource "metal_device" "test" {
    hostname            = "tfacc-device-key-test"
    plan                = "baremetal_0"
    facilities          = ["ewr1"]
    operating_system    = "ubuntu_16_04"
    billing_cycle       = "hourly"
    project_ssh_key_ids = ["${metal_project_ssh_key.test.id}"]
    project_id          = "${metal_project.test.id}"
}

`, name, publicSshKey)
}

func TestAccMetalProjectSSHKey_Basic(t *testing.T) {
	rs := acctest.RandString(10)
	var key packngo.SSHKey
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	cfg := metalProjectSSHKeyConfig_Basic(rs, publicKeyMaterial)
	log.Printf(cfg)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalProjectSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMetalSSHKeyExists("metal_project_ssh_key.test", &key),
					resource.TestCheckResourceAttr(
						"metal_project_ssh_key.test", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrPair(
						"metal_device.test", "ssh_key_ids.0",
						"metal_project_ssh_key.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"metal_project.test", "id",
						"metal_project_ssh_key.test", "project_id",
					),
				),
			},
		},
	})
}

func testAccCheckMetalProjectSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metal_project_ssh_key" {
			continue
		}
		if _, _, err := client.SSHKeys.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("SSH key still exists")
		}
	}

	return nil
}
