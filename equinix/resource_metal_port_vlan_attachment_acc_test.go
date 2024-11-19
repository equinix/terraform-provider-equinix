package equinix

import (
	"fmt"
	"path"
	"strings"
	"testing"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccMetalPortVlanAttachmentConfig_L2Bonded_1(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-port-vlan-attachment-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func testAccMetalPortVlanAttachmentConfig_L2Bonded_2(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_vlan" "test1" {
  description = "tfacc-vlan test VLAN 1"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_vlan" "test2" {
  description = "tfacc-vlan test VLAN 2"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type = "layer2-bonded"
}

resource "equinix_metal_port_vlan_attachment" "test1" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test1.vxlan
  port_name = "bond0"
}

resource "equinix_metal_port_vlan_attachment" "test2" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test2.vxlan
  port_name = "bond0"
}

`, testAccMetalPortVlanAttachmentConfig_L2Bonded_1(name))
}

func TestAccMetalPortVlanAttachment_L2Bonded(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalPortVlanAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalPortVlanAttachmentConfig_L2Bonded_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("equinix_metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccMetalPortVlanAttachmentConfig_L2Bonded_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test1", "port_name", "bond0"),
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test2", "port_name", "bond0"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test1", "device_id",
						"equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr("equinix_metal_device_network_type.test", "type", "layer2-bonded"),
				),
			},
		},
	})
}

func testAccMetalPortVlanAttachmentConfig_L2Individual_1(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-vlan-l2i-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}
`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func testAccMetalPortVlanAttachmentConfig_L2Individual_2(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_vlan" "test1" {
  description = "tfacc-vlan test VLAN 1"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_vlan" "test2" {
  description = "tfacc-vlan test VLAN 1"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type = "layer2-individual"
}

resource "equinix_metal_port_vlan_attachment" "test1" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test1.vxlan
  port_name = "eth1"
}

resource "equinix_metal_port_vlan_attachment" "test2" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test2.vxlan
  port_name = "eth1"
}

`, testAccMetalPortVlanAttachmentConfig_L2Individual_1(name))
}

func TestAccMetalPortVlanAttachment_L2Individual(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalPortVlanAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalPortVlanAttachmentConfig_L2Individual_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccMetalPortVlanAttachmentConfig_L2Individual_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test1", "device_id",
						"equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_device_network_type.test", "type", "layer2-individual"),
				),
			},
		},
	})
}

func testAccMetalPortVlanAttachmentConfig_Hybrid_1(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-hybrid-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func testAccMetalPortVlanAttachmentConfig_Hybrid_2(name string) string {
	return fmt.Sprintf(`
%s 

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type = "hybrid"
}

resource "equinix_metal_vlan" "test" {
  description = "tfacc-vlan test vlan"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_port_vlan_attachment" "test" {
  device_id  = equinix_metal_device_network_type.test.id
  vlan_vnid  = equinix_metal_vlan.test.vxlan
  port_name  = "eth1"
  force_bond = false
}`, testAccMetalPortVlanAttachmentConfig_Hybrid_1(name))
}

func TestAccMetalPortVlanAttachment_hybridBasic(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalPortVlanAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalPortVlanAttachmentConfig_Hybrid_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccMetalPortVlanAttachmentConfig_Hybrid_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test", "device_id",
						"equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_device_network_type.test", "type", "hybrid"),
				),
			},
		},
	})
}

func testAccMetalPortVlanAttachmentConfig_HybridMultipleVlans_1(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
  name = "tfacc-port_vlan_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-hmv-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func testAccMetalPortVlanAttachmentConfig_HybridMultipleVlans_2(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_vlan" "test" {
  count       = 3
  description = "tfacc-vlan test VLAN"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type      = "hybrid"
}

resource "equinix_metal_port_vlan_attachment" "test" {
  count     = length(equinix_metal_vlan.test)
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test[count.index].vxlan
  port_name = "eth1"
}`, testAccMetalPortVlanAttachmentConfig_HybridMultipleVlans_1(name))
}

func TestAccMetalPortVlanAttachment_hybridMultipleVlans(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalPortVlanAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalPortVlanAttachmentConfig_HybridMultipleVlans_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccMetalPortVlanAttachmentConfig_HybridMultipleVlans_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test.0", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test.0", "device_id", "equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test.1", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test.1", "device_id", "equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test.2", "port_name", "eth1"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test.2", "device_id", "equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_device_network_type.test", "type", "hybrid"),
				),
			},
		},
	})
}

func testAccMetalPortVlanAttachmentCheckDestroyed(s *terraform.State) error {
	client := testAccProvider.Meta().(*config.Config).Metal

	device_id := ""
	vlan_id := ""
	port_id := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "equinix_metal_device" {
			device_id = rs.Primary.ID
		}
		if rs.Type == "equinix_metal_port_vlan_attachment" {
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

func testAccMetalPortVlanAttachmentConfig_L2Native_1(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_project" "test" {
    name = "tfacc-port_vlan_attachment-%s"
}

resource "equinix_metal_device" "test" {
  hostname         = "tfacc-device-l2n-test"
  plan             = local.plan
  metro            = local.metro
  operating_system = local.os
  billing_cycle    = "hourly"
  project_id       = equinix_metal_project.test.id
  termination_time = "%s"
}`, confAccMetalDevice_base(preferable_plans, preferable_metros, preferable_os), name, testDeviceTerminationTime())
}

func testAccMetalPortVlanAttachmentConfig_L2Native_2(name string) string {
	return fmt.Sprintf(`
%s

resource "equinix_metal_vlan" "test1" {
  description = "tfacc-vlan test VLAN 1"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_vlan" "test2" {
  description = "tfacc-vlan test VLAN 2"
  metro       = equinix_metal_device.test.metro
  project_id  = equinix_metal_project.test.id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type = "layer2-individual"
}

resource "equinix_metal_port_vlan_attachment" "test1" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test1.vxlan
  port_name = "eth1"
}

resource "equinix_metal_port_vlan_attachment" "test2" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test2.vxlan
  native    = true
  port_name = "eth1"
  depends_on = ["equinix_metal_port_vlan_attachment.test1"]
}

`, testAccMetalPortVlanAttachmentConfig_L2Native_1(name))
}

func TestAccMetalPortVlanAttachment_L2Native(t *testing.T) {
	rs := acctest.RandString(10)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        testExternalProviders,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccMetalPortVlanAttachmentCheckDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccMetalPortVlanAttachmentConfig_L2Native_1(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_device.test", "network_type", "layer3"),
				),
			},
			{
				Config: testAccMetalPortVlanAttachmentConfig_L2Native_2(rs),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test1", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test2", "port_name", "eth1"),
					resource.TestCheckResourceAttr(
						"equinix_metal_port_vlan_attachment.test2", "native", "true"),
					resource.TestCheckResourceAttrPair(
						"equinix_metal_port_vlan_attachment.test1", "device_id",
						"equinix_metal_device.test", "id"),
					resource.TestCheckResourceAttr(
						"equinix_metal_device_network_type.test", "type", "layer2-individual"),
				),
			},
		},
	})
}
