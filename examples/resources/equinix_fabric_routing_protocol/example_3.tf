resource "equinix_fabric_routing_protocol" "direct"{
  connection_uuid = <some_id>
  type = "DIRECT"
  name = "direct_rp"
  direct_ipv4 {
  	equinix_iface_ip = "190.1.1.1/30"
  }
  direct_ipv6{
  	equinix_iface_ip = "190::1:1/126"
  }
}

resource "equinix_fabric_routing_protocol" "bgp" {
  depends_on = [
    equinix_fabric_routing_protocol.direct
  ]
  connection_uuid = <same_connection_id_as_first_equinix_fabric_routing_protocol>
  type            = "BGP"
  name            = "bgp_rp"
  bgp_ipv4 {
    customer_peer_ip = "190.1.1.2"
    enabled          = true
  }
  bgp_ipv6 {
    customer_peer_ip = "190::1:2"
    enabled          = true
  }
  customer_asn = 4532
}
