package packet

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/packethost/packngo"
)

func testAccCheckPacketPortVlanAttachmentConfig_L2Bonded_1(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}
`, name)
}

func testAccCheckPacketPortVlanAttachmentConfig_L2Bonded_2(name string) string {
	return fmt.Sprintf(`
%s

resource "packet_vlan" "test1" {
  description = "test VLAN 1"
  facility    = "nrt1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "test VLAN 2"
  facility    = "nrt1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_device_network_type" "test" {
  device_id = packet_device.test.id
  type = "layer2-bonded"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test1.vxlan}"
  port_name = "bond0"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test2.vxlan}"
  port_name = "bond0"
}

`, testAccCheckPacketPortVlanAttachmentConfig_L2Bonded_1(name))
}

func TestAccPacketPortVlanAttachment_L2Bonded(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Bonded_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("packet_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Bonded_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test1", "port_name", "bond0"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test2", "port_name", "bond0"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test1", "device_id",
						"packet_device.test", "id"),
					resource.TestCheckResourceAttr("packet_device_network_type.test", "type", "layer2-bonded"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentConfig_L2Individual_1(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-vlan-l2i-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}
`, name)
}

func testAccCheckPacketPortVlanAttachmentConfig_L2Individual_2(name string) string {
	return fmt.Sprintf(`
%s

resource "packet_vlan" "test1" {
  description = "test VLAN 1"
  facility    = "nrt1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "test VLAN 2"
  facility    = "nrt1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_device_network_type" "test" {
  device_id = packet_device.test.id
  type = "layer2-individual"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test1.vxlan}"
  port_name = "eth1"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test2.vxlan}"
  port_name = "eth1"
}

`, testAccCheckPacketPortVlanAttachmentConfig_L2Individual_1(name))
}

func TestAccPacketPortVlanAttachment_L2Individual(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Individual_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Individual_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test1", "device_id",
						"packet_device.test", "id"),
					resource.TestCheckResourceAttr(
						"packet_device_network_type.test", "type", "layer2-individual"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentConfig_Hybrid_1(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-device-hybrid-test"
  plan             = "n2.xlarge.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}`, name)
}

func testAccCheckPacketPortVlanAttachmentConfig_Hybrid_2(name string) string {
	return fmt.Sprintf(`
%s 

resource "packet_device_network_type" "test" {
  device_id = packet_device.test.id
  type = "hybrid"
}

resource "packet_vlan" "test" {
  description = "test vlan"
  facility    = "dfw2"
  project_id  = "${packet_project.test.id}"
}

resource "packet_port_vlan_attachment" "test" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test.vxlan}"
  port_name = "eth1"
  force_bond = false
}`, testAccCheckPacketPortVlanAttachmentConfig_Hybrid_1(name))
}

func TestAccPacketPortVlanAttachment_HybridBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_Hybrid_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_Hybrid_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_port_vlan_attachment.test", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"packet_port_vlan_attachment.test", "device_id",
						"packet_device.test", "id"),
					resource.TestCheckResourceAttr(
						"packet_device_network_type.test", "type", "hybrid"),
				),
			},
		},
	})
}

func testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans_1(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
  name = "tfacc-port_vlan_attachment-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-device-hmv-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = packet_project.test.id
}`, name)
}

func testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans_2(name string) string {
	return fmt.Sprintf(`
%s

resource "packet_vlan" "test" {
  count       = 3
  description = "test VLAN"
  facility    = "nrt1"
  project_id  = packet_project.test.id
}

resource "packet_device_network_type" "test" {
  device_id = packet_device.test.id
  type = "hybrid"
}

resource "packet_port_vlan_attachment" "test" {
  count     = length(packet_vlan.test)
  device_id = packet_device_network_type.test.id
  vlan_vnid = packet_vlan.test[count.index].vxlan
  port_name = "eth1"
}`, testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans_1(name))
}

func TestAccPacketPortVlanAttachment_HybridMultipleVlans(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_HybridMultipleVlans_2(rs),
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
					resource.TestCheckResourceAttr(
						"packet_device_network_type.test", "type", "hybrid"),
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

func testAccCheckPacketPortVlanAttachmentConfig_L2Native_1(name string) string {
	return fmt.Sprintf(`
resource "packet_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "packet_device" "test" {
  hostname         = "tfacc-device-l2n-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.test.id}"
}`, name)
}

func testAccCheckPacketPortVlanAttachmentConfig_L2Native_2(name string) string {
	return fmt.Sprintf(`
%s

resource "packet_vlan" "test1" {
  description = "test VLAN 1"
  facility    = "nrt1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_vlan" "test2" {
  description = "test VLAN 2"
  facility    = "nrt1"
  project_id  = "${packet_project.test.id}"
}

resource "packet_device_network_type" "test" {
  device_id = packet_device.test.id
  type = "layer2-individual"
}

resource "packet_port_vlan_attachment" "test1" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test1.vxlan}"
  port_name = "eth1"
}

resource "packet_port_vlan_attachment" "test2" {
  device_id = packet_device_network_type.test.id
  vlan_vnid = "${packet_vlan.test2.vxlan}"
  native    = true
  port_name = "eth1"
  depends_on = ["packet_port_vlan_attachment.test1"]
}

`, testAccCheckPacketPortVlanAttachmentConfig_L2Native_1(name))
}

func TestAccPacketPortVlanAttachment_L2Native(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPacketPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Native_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"packet_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckPacketPortVlanAttachmentConfig_L2Native_2(rs),
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
					resource.TestCheckResourceAttr(
						"packet_device_network_type.test", "type", "layer2-individual"),
				),
			},
		},
	})
}
