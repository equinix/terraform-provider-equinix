data "equinix_fabric_routing_protocol" "routing_protocol_data_name" {
  connection_uuid = "<uuid_of_connection_routing_protocol_is_applied_to>"
  uuid = "<uuid_of_routing_protocol>"
}

output "id" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.id
}

output "name" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.name
}

output "type" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.type
}

output "direct_ipv4" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.direct_ipv4.0.equinix_iface_ip
}

output "direct_ipv6" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.direct_ipv6.0.equinix_iface_ip
}

output "bgp_ipv4_customer_peer_ip" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.bgp_ipv4.0.customer_peer_ip
}

output "bgp_ipv4_equinix_peer_ip" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.bgp_ipv4.0.equinix_peer_ip
}

output "bgp_ipv4_enabled" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.bgp_ipv4.0.enabled
}

output "bgp_ipv6_customer_peer_ip" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.bgp_ipv6.0.customer_peer_ip
}

output "bgp_ipv6_equinix_peer_ip" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.bgp_ipv6.0.equinix_peer_ip
}

output "bgp_ipv6_enabled" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.bgp_ipv6.0.enabled
}

output "customer_asn" {
  value = data.equinix_fabric_routing_protocol.routing_protocol_data_name.customer_asn
}
