resource "equinix_metal_vlan" "test" {
  description = "VLAN in New York"
  metro       = "ny"
  project_id  = local.project_id
}

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
  type      = "hybrid"
}

resource "equinix_metal_port_vlan_attachment" "test" {
  device_id = equinix_metal_device_network_type.test.id
  port_name = "eth1"
  vlan_vnid = equinix_metal_vlan.test.vxlan
}

