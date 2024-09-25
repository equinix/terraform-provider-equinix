data "equinix_metal_vlan" "dsvlan" {
  project_id = local.project_id
  vxlan      = 5
  metro      = "sv"
}
