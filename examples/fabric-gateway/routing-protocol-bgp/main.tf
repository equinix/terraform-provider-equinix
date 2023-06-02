provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_routingprotocol" "test"{
  connection_uuid = var.connection_uuid
  type = var.rp_type
  name = var.rp_name
  bgp_ipv4 {
  	customer_peer_ip = var.customer_peer_ipv4
  }
  bgp_ipv6{
  	customer_peer_ip = var.customer_peer_ipv6
  }
  customer_asn = var.customer_asn
}

output "rp_result" {
  value = equinix_fabric_routingprotocol.test.id
}

