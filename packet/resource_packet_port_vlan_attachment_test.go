package packet

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckPacketPortVlanAttachmentConfig_L2Bonded(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "TerraformTestProject-%s"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "s1.large.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  network_type     = "layer2-bonded"
}

resource "packet_vlan" "test1" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test1.vxlan}"
  port_name = "bond0"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test2.vxlan}"
  port_name = "bond0"
}

`, name)
}

func TestAccPacketPortVlanAttachment_L2Bonded(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Bonded(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test1", "port_name", "bond0"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test2", "port_name", "bond0"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test1", "device_id",
						"packet_device.test", "id"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentConfig_L2Individual(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "TerraformTestProject-%s"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "s1.large.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  network_type     = "layer2-individual"
}

resource "packet_vlan" "test1" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test1.vxlan}"
  port_name = "eth1"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test2.vxlan}"
  port_name = "eth1"
}

`, name)
}

func TestAccPacketPortVlanAttachment_L2Individual(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Individual(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test1", "device_id",
						"packet_device.test", "id"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentConfig_Hybrid(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "TerraformTestProject-%s"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "s1.large.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  network_type     = "hybrid"
}

resource "packet_vlan" "test" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_port_vlan_attachment" "test" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test.vxlan}"
  port_name = "eth1"
  force_bond = false
}`, name)
}

func TestAccPacketPortVlanAttachment_Hybrid(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_Hybrid(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test", "device_id",
						"packet_device.test", "id"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
  name = "TerraformTestProject-%s"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "s1.large.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = packet_project.test.id
  network_type     = "hybrid"
}

resource "packet_vlan" "test" {
  count       = 3
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = packet_project.test.id
}

resource "packet_port_vlan_attachment" "test" {
  count     = length(packet_vlan.test)
  device_id = packet_device.test.id
  vlan_vnid = packet_vlan.test[count.index].vxlan
  port_name = "eth1"
}`, name)
}

func TestAccPacketPortVlanAttachment_HybridMultipleVlans(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test.0", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test.0", "device_id", "packet_device.test", "id"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test.1", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test.1", "device_id", "packet_device.test", "id"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test.2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test.2", "device_id", "packet_device.test", "id"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	device_id := ""
	vlan_id := ""
	port_id := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "packet_device" {
			device_id = rs.Primary.ID
		}
		if rs.Type == "packet_port_vlan_attachment" {
			port_vlan := strings.Split(rs.Primary.ID, ":")
			vlan_id = port_vlan[0]
			port_id = port_vlan[1]

		}
	}
	d, _, err := client.Devices.Get(device_id, nil)
	if err != nil {
		// if device doesn't exists, its port can't be attached
		return nil
	}
	for _, p := range d.NetworkPorts {
		if p.ID == port_id {
			if len(p.AttachedVirtualNetworks) == 1 {
				if path.Base(p.AttachedVirtualNetworks[0].Href) == vlan_id {
					return fmt.Errorf("Vlan is still attached to the device")
				}
			}
		}
	}

	return nil
}

func testAccCheckPacketPortVlanAttachmentConfig_L2Native(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "TerraformTestProject-%s"
}

resource "packet_device" "test" {
  hostname         = "test"
  plan             = "s1.large.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
  network_type     = "layer2-individual"
}

resource "packet_vlan" "test1" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "VLAN in New Jersey"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test1.vxlan}"
  port_name = "eth1"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = "${packet_device.test.id}"
  vlan_vnid = "${packet_vlan.test2.vxlan}"
  native    = true
  port_name = "eth1"
  depends_on = ["packet_port_vlan_attachment.test1"]
}

`, name)
}

func TestAccPacketPortVlanAttachment_L2Native(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Native(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test2", "native", "true"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test1", "device_id",
						"packet_device.test", "id"),
				),
			},
		},
	})
}
