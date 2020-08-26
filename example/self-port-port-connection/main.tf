provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

data "equinix_ecx_port" "a-port" {
  name = var.equinix_aside_port_name
}

data "equinix_ecx_port" "z-port" {
  name = var.equinix_zside_port_name
}

resource "equinix_ecx_l2_connection" "port-port" {
  name            = "tf-port-port"
  speed           = 50
  speed_unit      = "MB"
  notifications   = ["example@equinix.com"]
  port_uuid       = data.equinix_ecx_port.a-port.uuid
  vlan_stag       = 1010
  zside_port_uuid = data.equinix_ecx_port.z-port.uuid
  zside_vlan_stag = 4004
}
