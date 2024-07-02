# Create Metal Gateway for a VLAN and reserved IP address block

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "equinix_metal_reserved_ip_block" "test" {
  project_id = local.project_id
  metro      = "sv"
  quantity   = 8
}

resource "equinix_metal_gateway" "test" {
  project_id        = local.project_id
  vlan_id           = equinix_metal_vlan.test.id
  ip_reservation_id = equinix_metal_reserved_ip_block.test.id
}
