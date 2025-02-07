data "equinix_fabric_route_aggregations" "ra_policy" {
  filter =  {
    property = "/project/projectId"
    operator = "="
    values    = ["<route_aggregation_project_id>"]
  }
  pagination = {
    limit = 2
    offset = 1
  }
}

output "first_route_aggregation_name" {
  value = data.equinix_fabric_route_aggregations.ra_policy.data.0.name
}

output "first_route_aggregation_description" {
  value = data.equinix_fabric_route_aggregations.ra_policy.data.0.description
}

output "first_route_aggregation_connections_count" {
  value = data.equinix_fabric_route_aggregations.ra_policy.data.0.connections_count
}

output "first_route_aggregation_rules_count" {
  value = data.equinix_fabric_route_aggregations.ra_policy.data.0.rules_count
}
