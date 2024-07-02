resource "equinix_metal_vlan" "foovlan" {
  project_id = local.project_id
  metro = "sv"
  vxlan = 5
}

data "equinix_metal_vlan" "dsvlan" {
  vlan_id = equinix_metal_vlan.foovlan.id
}
