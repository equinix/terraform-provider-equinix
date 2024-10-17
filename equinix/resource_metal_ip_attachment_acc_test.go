package equinix

import (
	"fmt"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccMetalIPAttachment_basic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalIPAttachmentCheckDestroyed,
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
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-ip_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-ip-attachment-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
    project_id = equinix_metal_project.test.id
    facility   = equinix_metal_device.test.deployed_facility
    quantity   = 2
}

resource "equinix_metal_ip_attachment" "test" {
	device_id = equinix_metal_device.test.id
	cidr_notation = "${cidrhost(equinix_metal_reserved_ip_block.test.cidr_notation,0)}/32"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func TestAccMetalIPAttachment_metro(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalIPAttachmentCheckDestroyed,
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
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-ip_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-ip-attachment-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}

resource "equinix_metal_reserved_ip_block" "test" {
    project_id = equinix_metal_project.test.id
    metro      = equinix_metal_device.test.metro
    quantity   = 2
}


resource "equinix_metal_ip_attachment" "test" {
	device_id = equinix_metal_device.test.id
	cidr_notation = "${cidrhost(equinix_metal_reserved_ip_block.test.cidr_notation,0)}/32"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func testAccMetalIPAttachmentCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).Metal

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
