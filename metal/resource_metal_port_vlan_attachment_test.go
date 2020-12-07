package metal

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

func testAccCheckMetalPortVlanAttachmentConfig_L2Bonded_1(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}
`, name)
}

func testAccCheckMetalPortVlanAttachmentConfig_L2Bonded_2(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_vlan" "test1" {
  description = "test VLAN 1"
  facility    = "nrt1"
  project_id  = "${metal_project.test.id}"
}

resource "metal_vlan" "test2" {
  description = "test VLAN 2"
  facility    = "nrt1"
  project_id  = "${metal_project.test.id}"
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type = "layer2-bonded"
}

resource "metal_port_vlan_attachment" "test1" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test1.vxlan}"
  port_name = "bond0"
}

resource "metal_port_vlan_attachment" "test2" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test2.vxlan}"
  port_name = "bond0"
}

`, testAccCheckMetalPortVlanAttachmentConfig_L2Bonded_1(name))
}

func TestAccMetalPortVlanAttachment_L2Bonded(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_L2Bonded_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_L2Bonded_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test1", "port_name", "bond0"),
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test2", "port_name", "bond0"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test1", "device_id",
						"metal_device.test", "id"),
					resource.TestCheckResourceAttr("metal_device_network_type.test", "type", "layer2-bonded"),
				),
			},
		},
	})
}

func testAccCheckMetalPortVlanAttachmentConfig_L2Individual_1(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-vlan-l2i-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}
`, name)
}

func testAccCheckMetalPortVlanAttachmentConfig_L2Individual_2(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_vlan" "test1" {
  description = "test VLAN 1"
  facility    = "nrt1"
  project_id  = "${metal_project.test.id}"
}

resource "metal_vlan" "test2" {
  description = "test VLAN 2"
  facility    = "nrt1"
  project_id  = "${metal_project.test.id}"
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type = "layer2-individual"
}

resource "metal_port_vlan_attachment" "test1" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test1.vxlan}"
  port_name = "eth1"
}

resource "metal_port_vlan_attachment" "test2" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test2.vxlan}"
  port_name = "eth1"
}

`, testAccCheckMetalPortVlanAttachmentConfig_L2Individual_1(name))
}

func TestAccMetalPortVlanAttachment_L2Individual(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_L2Individual_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_L2Individual_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test1", "device_id",
						"metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_device_network_type.test", "type", "layer2-individual"),
				),
			},
		},
	})
}

func testAccCheckMetalPortVlanAttachmentConfig_Hybrid_1(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-hybrid-test"
  plan             = "n2.xlarge.x86"
  facilities       = ["dfw2"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}`, name)
}

func testAccCheckMetalPortVlanAttachmentConfig_Hybrid_2(name string) string {
	return fmt.Sprintf(`
%s 

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type = "hybrid"
}

resource "metal_vlan" "test" {
  description = "test vlan"
  facility    = "dfw2"
  project_id  = "${metal_project.test.id}"
}

resource "metal_port_vlan_attachment" "test" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test.vxlan}"
  port_name = "eth1"
  force_bond = false
}`, testAccCheckMetalPortVlanAttachmentConfig_Hybrid_1(name))
}

func TestAccMetalPortVlanAttachment_HybridBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_Hybrid_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_Hybrid_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test", "device_id",
						"metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_device_network_type.test", "type", "hybrid"),
				),
			},
		},
	})
}

func testAccCheckMetalPortVlanAttachmentConfig_HybridMultipleVlans_1(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
  name = "tfacc-port_vlan_attachment-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-hmv-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = metal_project.test.id
}`, name)
}

func testAccCheckMetalPortVlanAttachmentConfig_HybridMultipleVlans_2(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_vlan" "test" {
  count       = 3
  description = "test VLAN"
  facility    = "nrt1"
  project_id  = metal_project.test.id
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type = "hybrid"
}

resource "metal_port_vlan_attachment" "test" {
  count     = length(metal_vlan.test)
  device_id = metal_device_network_type.test.id
  vlan_vnid = metal_vlan.test[count.index].vxlan
  port_name = "eth1"
}`, testAccCheckMetalPortVlanAttachmentConfig_HybridMultipleVlans_1(name))
}

func TestAccMetalPortVlanAttachment_HybridMultipleVlans(t *testing.T) {
	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_HybridMultipleVlans_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_HybridMultipleVlans_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test.0", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test.0", "device_id", "metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test.1", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test.1", "device_id", "metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test.2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test.2", "device_id", "metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_device_network_type.test", "type", "hybrid"),
				),
			},
		},
	})
}

func testAccCheckMetalPortVlanAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*packngo.Client)

	device_id := ""
	vlan_id := ""
	port_id := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "metal_device" {
			device_id = rs.Primary.ID
		}
		if rs.Type == "metal_port_vlan_attachment" {
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

func testAccCheckMetalPortVlanAttachmentConfig_L2Native_1(name string) string {
	return fmt.Sprintf(`
resource "metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "metal_device" "test" {
  hostname         = "tfacc-device-l2n-test"
  plan             = "s1.large.x86"
  facilities       = ["nrt1"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${metal_project.test.id}"
}`, name)
}

func testAccCheckMetalPortVlanAttachmentConfig_L2Native_2(name string) string {
	return fmt.Sprintf(`
%s

resource "metal_vlan" "test1" {
  description = "test VLAN 1"
  facility    = "nrt1"
  project_id  = "${metal_project.test.id}"
}

resource "metal_vlan" "test2" {
  description = "test VLAN 2"
  facility    = "nrt1"
  project_id  = "${metal_project.test.id}"
}

resource "metal_device_network_type" "test" {
  device_id = metal_device.test.id
  type = "layer2-individual"
}

resource "metal_port_vlan_attachment" "test1" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test1.vxlan}"
  port_name = "eth1"
}

resource "metal_port_vlan_attachment" "test2" {
  device_id = metal_device_network_type.test.id
  vlan_vnid = "${metal_vlan.test2.vxlan}"
  native    = true
  port_name = "eth1"
  depends_on = ["metal_port_vlan_attachment.test1"]
}

`, testAccCheckMetalPortVlanAttachmentConfig_L2Native_1(name))
}

func TestAccMetalPortVlanAttachment_L2Native(t *testing.T) {

	rs := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetalPortVlanAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_L2Native_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccCheckMetalPortVlanAttachmentConfig_L2Native_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"metal_port_vlan_attachment.test2", "native", "true"),
					resource.TestCheckResourceAttrPair(
						"metal_port_vlan_attachment.test1", "device_id",
						"metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"metal_device_network_type.test", "type", "layer2-individual"),
				),
			},
		},
	})
}
