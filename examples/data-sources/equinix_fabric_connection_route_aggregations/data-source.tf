data "equinix_fabric_connection_route_aggregations" "attached_policies" {
  connection_id = "connection_id"
}

output "connection_first_route_Aggregation_uuid" {
  value = data.equinix_fabric_connection_route_aggregations.attached_policies.data.0.uuid
}

output "connection_first_route_aggregation_type" {
  value = data.equinix_fabric_connection_route_aggregations.attached_policies.data.0.type
}

output "connection_first_route_aggregation_attachment_status" {
  value = data.equinix_fabric_connection_route_aggregations.attached_policies.data.0.attachment_status
}
