data "equinix_fabric_route_aggregation_rule" "ra_rule" {
  route_aggregation_id = "<route_aggregation_id>"
  route_aggregation_rule_id = "<route_aggregation_rule_id>"
}

output "route_aggregation_rule_name" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.name
}

output "route_aggregation_rule_description" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.description
}

output "route_aggregation_rule_type" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.type
}

output "route_aggregation_rule_prefix" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.prefix
}

output "route_aggregation_rule_state" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.state
}