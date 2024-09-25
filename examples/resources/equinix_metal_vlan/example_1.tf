# Create a new VLAN in metro "esv"
resource "equinix_metal_vlan" "vlan1" {
  description = "VLAN in New Jersey"
  metro       = "sv"
  project_id  = local.project_id
  vxlan       = 1040
}
