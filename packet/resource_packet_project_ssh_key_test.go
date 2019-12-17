package packet

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func packetProjectSSHKeyConfig_Basic(name, publicSshKey string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-project_ssh_key-%s"
}

resource "packet_project_ssh_key" "test" {
    name = "tfacc-project-key-test"
    public_key = "%s"
    project_id = "${packet_project.test.id}"
}

resource "packet_device" "test" {
    hostname            = "tfacc-device-key-test"
    plan                = "baremetal_0"
    facilities          = ["ewr1"]
    operating_system    = "ubuntu_16_04"
    billing_cycle       = "hourly"
    project_ssh_key_ids = ["${packet_project_ssh_key.test.id}"]
    project_id          = "${packet_project.test.id}"
}

`, name, publicSshKey)
}

func TestAccPacketProjectSSHKey_Basic(t *testing.T) {
	rs := acctest.RandString(10)
	var key packngo.SSHKey
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	cfg := packetProjectSSHKeyConfig_Basic(rs, publicKeyMaterial)
	log.Printf(cfg)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketProjectSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPacketSSHKeyExists("packet_project_ssh_key.test", &key),
					resource.TestCheckResourceAttr(
						"packet_project_ssh_key.test", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrPair(
						"packet_device.test", "ssh_key_ids.0",
						"packet_project_ssh_key.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"packet_project.test", "id",
						"packet_project_ssh_key.test", "project_id",
					),
				),
			},
		},
	})
}

func testAccCheckPacketProjectSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "packet_project_ssh_key" {
			continue
		}
		if _, _, err := client.SSHKeys.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("SSH key still exists")
		}
	}

	return nil
}
