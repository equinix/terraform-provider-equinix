# Create Metal Gateway for a VLAN with a private IPv4 block with 8 IP addresses

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "equinix_metal_gateway" "test" {
  project_id               = local.project_id
  vlan_id                  = equinix_metal_vlan.test.id
  private_ipv4_subnet_size = 8
}
