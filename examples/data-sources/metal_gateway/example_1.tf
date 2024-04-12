# Create Metal Gateway for a VLAN with a private IPv4 block with 8 IP addresses

resource "equinix_metal_vlan" "test" {
  description = "test VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

data "equinix_metal_gateway" "test" {
  gateway_id = local.gateway_id
}
