data "equinix_connection_route_filters" "attached_policies" {
  connection_id = "<connection_uuid>"
}

output "connection_first_route_filter_uuid" {
  value = data.equinix_fabric_connection_route_filter.attached_policies.0.uuid
}

output "connection_first_route_filter_connection_id" {
  value = data.equinix_fabric_connection_route_filter.attached_policies.0.connection_id
}

output "connection_first_route_filter_direction" {
  value = data.equinix_fabric_connection_route_filter.attached_policies.0.direction
}

output "connection_first_route_filter_type" {
  value = data.equinix_fabric_connection_route_filter.attached_policies.0.type
}

output "connection_first_route_filter_attachment_status" {
  value = data.equinix_fabric_connection_route_filter.attached_policies.0.attachment_status
}
