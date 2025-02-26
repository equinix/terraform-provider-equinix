resource "equinix_fabric_connection_route_aggregation" "policy_attachment" {
  route_aggregation_id = "<route_aggregation_id>"
  connection_id = "<connection_id>"
}

output "connection_route_Aggregation_id" {
  value = equinix_fabric_connection_route_aggregation.policy_attachment.id
}

output "connection_route_aggregation_connection_id" {
  value = equinix_fabric_connection_route_aggregation.policy_attachment.connection_id
}

output "connection_route_aggregation_type" {
  value = equinix_fabric_connection_route_aggregation.policy_attachment.type
}

output "connection_route_aggregation_attachment_status" {
  value = equinix_fabric_connection_route_aggregation.policy_attachment.attachment_status
}
