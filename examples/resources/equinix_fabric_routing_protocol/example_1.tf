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
