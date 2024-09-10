data "equinix_fabric_networks" "test" {
    outer_operator = "AND"
    filter {
        property = "/type"
        operator = "="
        values 	 = ["IPWAN"]
    }
    filter {
        property = "/name"
        operator = "="
        values   = ["Tf_Network_PFCR"]
    }
    filter {
        group = "OR_group1"
        property = "/operation/equinixStatus"
        operator = "="
        values = ["PROVISIONED"]
    }
    filter {
        group = "OR_group1"
        property = "/operation/equinixStatus"
        operator = "LIKE"
        values = ["DEPROVISIONED"]
    }
    pagination {
        offset = 0
        limit = 5
    }
    sort {
        direction = "ASC"
        property = "/name"
    }
}

output "number_of_returned_networks" {
    value = length(data.equinix_fabric_networks.test.data)
}

output "first_network_name" {
    value = data.equinix_fabric_networks.test.data.0.name
}

output "first_network_connections_count" {
    value = data.equinix_fabric_networks.test.data.0.connections_count
}

output "first_network_scope" {
    value = data.equinix_fabric_networks.test.data.0.scope
}

output "first_network_type" {
    value = data.equinix_fabric_networks.test.data.0.type
}

output "first_network_location_region" {
    value = one(data.equinix_fabric_networks.test.data.0.location).region
}

output "first_network_project_id" {
    value = one(data.equinix_fabric_networks.test.data.0.project).project_id
}
