resource "equinix_metal_vlan" "example" {
  description = "proj-vrf-bgp-neighbor-example VLAN in SV"
  metro       = "sv"
  project_id  = local.project_id
}

resource "equinix_metal_vrf" "example" {
  description = "proj-vrf-bgp-neighbor-example VRF in SV"
  name        = "tfacc-vrf-example"
  metro       = "sv"
  local_asn   = "65000"
  ip_ranges   = ["2001:d78::/59"]
  project_id  = local.project_id
}

resource "equinix_metal_reserved_ip_block" "example" {
  project_id = local.project_id
  type       = "vrf"
  vrf_id     = equinix_metal_vrf.example.id
  network    = "2001:d78::"
  metro      = "sv"
  cidr       = 64
}

resource "equinix_metal_gateway" "example" {
  project_id        = local.project_id
  vlan_id           = equinix_metal_vlan.example.id
  ip_reservation_id = equinix_metal_reserved_ip_block.example.id
}

resource "equinix_metal_vrf_bgp_dynamic_neighbor" "example" {
  gateway_id = equinix_metal_gateway.example.id
  range      = "2001:d78:0:0:4000::/66"
  asn        = "56789"
}
