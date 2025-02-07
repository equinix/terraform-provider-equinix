data "equinix_fabric_route_aggregation" "ra_policy" {
  uuid = "<uuid_of_route_aggregation>"
}

output "id" {
  value = data.equinix_fabric_route_aggregation.ra_policy.id
}

output "type" {
  value = data.equinix_fabric_route_aggregation.ra_policy.type
}

output "state" {
  value = data.equinix_fabric_route_aggregation.ra_policy.state
}

output "connections_count" {
  value = data.equinix_fabric_route_aggregation.ra_policy.connections_count
}

output "rules_count" {
  value = data.equinix_fabric_route_aggregation.ra_policy.rules_count
}
