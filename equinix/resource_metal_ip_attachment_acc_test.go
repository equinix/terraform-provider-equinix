package equinix

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMetalIPAttachment_basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalIPAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalIPAttachmentConfig_basic(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_ip_attachment.test", "public", "true"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_ip_attachment.test", "device_id",
						"equinix_metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "equinix_metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalIPAttachmentConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-ip_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-ip-attachment-test"
  plan             = "c3.small.x86"
  facilities       = ["sv15"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
}

resource "equinix_metal_reserved_ip_block" "test" {
    project_id = equinix_metal_project.test.id
    facility   = "sv15"
	quantity   = 2
}


resource "equinix_metal_ip_attachment" "test" {
	device_id = equinix_metal_device.test.id
	cidr_notation = "${cidrhost(equinix_metal_reserved_ip_block.test.cidr_notation,0)}/32"
}`, name)
}

func TestAccMetalIPAttachment_metro(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccMetalIPAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalIPAttachmentConfig_metro(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_ip_attachment.test", "public", "true"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_ip_attachment.test", "device_id",
						"equinix_metal_device.test", "id"),
				),
			},
			{
				ResourceName:      "equinix_metal_ip_attachment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMetalIPAttachmentConfig_metro(name string) string {
	return fmt.Sprintf(`
resource "equinix_metal_project" "test" {
    name = "tfacc-ip_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-ip-attachment-test"
  plan             = "c3.medium.x86"
  metro            = "sv"
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
}

resource "equinix_metal_reserved_ip_block" "test" {
    project_id = equinix_metal_project.test.id
    metro      = "sv"
	quantity   = 2
}


resource "equinix_metal_ip_attachment" "test" {
	device_id = equinix_metal_device.test.id
	cidr_notation = "${cidrhost(equinix_metal_reserved_ip_block.test.cidr_notation,0)}/32"
}`, name)
}

func testAccMetalIPAttachmentCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config).metal

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "equinix_metal_ip_attachment" {
			continue
		}
		if _, _, err := client.ProjectIPs.Get(rs.Primary.ID, nil); err == nil {
			return fmt.Errorf("Metal IP attachment still exists")
		}
	}

	return nil
}
