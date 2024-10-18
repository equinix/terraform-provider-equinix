data "equinix_fabric_connection_route_filter" "attached_policy" {
  connection_id = "<connection_uuid>"
  route_filter_id = "<route_filter_uuid>"
}

output "connection_route_filter_id" {
  value = data.equinix_fabric_connection_route_filter.attached_policy.id
}

output "connection_route_filter_connection_id" {
  value = data.equinix_fabric_connection_route_filter.attached_policy.connection_id
}

output "connection_route_filter_direction" {
  value = data.equinix_fabric_connection_route_filter.attached_policy.direction
}

output "connection_route_filter_type" {
  value = data.equinix_fabric_connection_route_filter.attached_policy.type
}

output "connection_route_filter_attachment_status" {
  value = data.equinix_fabric_connection_route_filter.attached_policy.attachment_status
}
