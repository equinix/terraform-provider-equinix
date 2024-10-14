resource "equinix_fabric_connection_route_filter" "policy_attachment" {
  connection_id = "<connection_uuid>"
  route_filter_id = "<route_filter_policy_uuid>"
  direction = "INBOUND"
}

output "connection_route_filter_id" {
  value = equinix_fabric_connection_route_filter.policy_attachment.id
}

output "connection_route_filter_connection_id" {
  value = equinix_fabric_connection_route_filter.policy_attachment.connection_id
}

output "connection_route_filter_direction" {
  value = equinix_fabric_connection_route_filter.policy_attachment.direction
}

output "connection_route_filter_type" {
  value = equinix_fabric_connection_route_filter.policy_attachment.type
}

output "connection_route_filter_attachment_status" {
  value = equinix_fabric_connection_route_filter.policy_attachment.attachment_status
}
