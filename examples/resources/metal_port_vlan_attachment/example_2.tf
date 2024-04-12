resource "equinix_metal_device" "test" {
  hostname         = "test"
  plan             = "c3.small.x86"
  metro            = "ny"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
}

resource "equinix_metal_device_network_type" "test" {
  device_id = equinix_metal_device.test.id
  type      = "layer2-individual"
}

resource "equinix_metal_vlan" "test1" {
  description = "VLAN in New York"
  metro       = "ny"
  project_id  = local.project_id
}

resource "equinix_metal_vlan" "test2" {
  description = "VLAN in New Jersey"
  metro       = "ny"
  project_id  = local.project_id
}

resource "equinix_metal_port_vlan_attachment" "test1" {
  device_id = equinix_metal_device_network_type.test.id
  vlan_vnid = equinix_metal_vlan.test1.vxlan
  port_name = "eth1"
}

resource "equinix_metal_port_vlan_attachment" "test2" {
  device_id  = equinix_metal_device_network_type.test.id
  vlan_vnid  = equinix_metal_vlan.test2.vxlan
  port_name  = "eth1"
  native     = true
  depends_on = ["equinix_metal_port_vlan_attachment.test1"]
}
