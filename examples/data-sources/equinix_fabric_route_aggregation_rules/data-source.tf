data "equinix_fabric_route_aggregation_rules" "ra_rules" {
  route_aggregation_id = "<route_aggregation_id>"
  pagination = {
    limit = 2
    offset = 1
  }
}

output "route_aggregation_rule_name" {
  value = data.equinix_fabric_route_aggregation_rules.ra_rules.data.0.name
}

output "route_aggregation_rule_description" {
  value = data.equinix_fabric_route_aggregation_rules.ra_rules.data.0.description
}

output "route_aggregation_rule_prefix" {
  value = data.equinix_fabric_route_aggregation_rules.ra_rules.data.0.prefix
}

output "route_aggregation_rule_state" {
  value = data.equinix_fabric_route_aggregation_rules.ra_rules.data.0.state
}
