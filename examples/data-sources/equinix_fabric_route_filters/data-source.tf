data "equinix_fabric_route_filters" "rf_policies" {
  filter {
    property = "/type"
    operator = "="
    values 	 = ["BGP_IPv4_PREFIX_FILTER"]
  }
  filter {
    property = "/state"
    operator = "="
    values   = ["PROVISIONED"]
  }
  filter {
    property = "/project/projectId"
    operator = "="
    values = ["<project_id>"]
  }
  pagination {
    offset = 0
    limit = 5
    total = 25
  }
  sort {
    direction = "ASC"
    property = "/name"
  }
}

output "first_rf_uuid" {
  value = data.equinix_fabric_route_filters.rf_policies.data.0.uuid
}

output "type" {
  value = data.equinix_fabric_route_filters.rf_policies.data.0.type
}

output "state" {
  value = data.equinix_fabric_route_filters.rf_policies.data.0.state
}

output "not_matched_rule_action" {
  value = data.equinix_fabric_route_filters.rf_policies.data.0.not_matched_rule_action
}

output "connections_count" {
  value = data.equinix_fabric_route_filters.rf_policies.data.0.connections_count
}

output "rules_count" {
  value = data.equinix_fabric_route_filters.rf_policies.data.0.rules_count
}
