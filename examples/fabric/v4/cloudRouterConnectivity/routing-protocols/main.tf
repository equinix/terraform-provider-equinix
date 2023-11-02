provider "equinix" {
  client_id     = var.equinix_client_id
  client_secret = var.equinix_client_secret
}

resource "equinix_fabric_routing_protocols" "routing_protocols" {
  connection_uuid = var.connection_uuid

  direct_routing_protocol {
    name = var.direct_rp_name
    direct_ipv4 {
      equinix_iface_ip = var.equinix_ipv4_ip
    }
    direct_ipv6 {
      equinix_iface_ip = var.equinix_ipv6_ip
    }
  }

  bgp_routing_protocol {
    name = var.bgp_rp_name
    bgp_ipv4 {
      customer_peer_ip = var.customer_peer_ipv4
      enabled          = var.bgp_enabled_ipv4
    }
    bgp_ipv6 {
      customer_peer_ip = var.customer_peer_ipv6
      enabled          = var.bgp_enabled_ipv6
    }
    customer_asn = var.customer_asn
    equinix_asn  = var.equinix_asn
  }
}



output "rp_result" {
  value = equinix_fabric_routing_protocols.routing_protocols.id
}

