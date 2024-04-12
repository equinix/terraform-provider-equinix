resource "equinix_metal_vlan" "example1" {
    project_id      = local.my_project_id
    metro           = "SV"
}

resource "equinix_metal_vlan" "example2" {
    project_id      = local.my_project_id
    metro           = "SV"
}

resource "equinix_metal_connection" "example" {
    name            = "tf-port-to-metal-legacy"
    project_id      = local.my_project_id
    metro           = "SV"
    redundancy      = "redundant"
    type            = "shared"
    contact_email   = "username@example.com"
    vlans              = [
      equinix_metal_vlan.example1.vxlan,
      equinix_metal_vlan.example2.vxlan
    ]
}

data "equinix_fabric_port" "example" {
  name = "CX-FR5-NL-Dot1q-BO-1G-PRI"
}

resource "equinix_fabric_connection" "example" {
  name                = "tf-port-to-metal-legacy"
  speed               = "200"
  speed_unit          = "MB"
  notifications       = ["example@equinix.com"]
  port_uuid           = data.equinix_fabric_port.example.id
  vlan_stag           = 1020
  authorization_key   = equinix_metal_connection.example.token
}
