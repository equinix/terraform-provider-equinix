package packet

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func packetProjectSSHKeyConfig_Basic(publicSshKey string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
	name = "test"
}

resource "packet_project_ssh_key" "test" {
    name = "test"
    public_key = "%s"
    project_id = "${packet_project.test.id}"
}

resource "packet_device" "test" {
    hostname            = "test"
    plan                = "baremetal_0"
    facility            = "ewr1"
    operating_system    = "ubuntu_16_04"
    billing_cycle       = "hourly"
	project_ssh_key_ids = ["${packet_project_ssh_key.test.id}"]
    project_id          = "${packet_project.test.id}"
}

`, publicSshKey)
}

func TestAccPacketProjectSSHKey_Basic(t *testing.T) {
	var key packngo.SSHKey
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	cfg := packetProjectSSHKeyConfig_Basic(publicKeyMaterial)
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
