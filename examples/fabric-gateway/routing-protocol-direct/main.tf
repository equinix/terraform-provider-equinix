provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_routingprotocol" "test"{
  connection_uuid = var.connection_uuid
  type = var.rp_type
  name = var.rp_name
  direct_ipv4 {
  	equinix_iface_ip = var.equinix_ipv4_ip
  }
  direct_ipv6{
  	equinix_iface_ip = var.equinix_ipv6_ip
  }
}

output "rp_result" {
  value = equinix_fabric_routingprotocol.test.id
}

