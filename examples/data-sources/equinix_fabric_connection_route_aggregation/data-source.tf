data "equinix_fabric_connection_route_aggregation" "attached_policy" {
  route_aggregation_id = "<route_aggregation_id>"
  connection_id = "<connection_id>"
}

output "connection_route_Aggregation_id" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.id
}

output "connection_route_aggregation_connection_id" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.connection_id
}

output "connection_route_aggregation_type" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.type
}

output "connection_route_aggregation_attachment_status" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.attachment_status
}