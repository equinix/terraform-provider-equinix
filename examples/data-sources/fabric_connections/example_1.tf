data "equinix_fabric_connections" "test" {
    outer_operator = "AND"
    filter {
        property = "/name"
        operator = "LIKE"
        values 	 = ["PNFV"]
    }
    filter {
        property = "/aSide/accessPoint/location/metroCode"
        operator = "="
        values   = ["SY"]
    }
    filter {
        group = "OR_group1"
        property = "/redundancy/priority"
        operator = "="
        values = ["PRIMARY"]
    }
    filter {
        group = "OR_group1"
        property = "/redundancy/priority"
        operator = "="
        values = ["SECONDARY"]
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

output "number_of_returned_connections" {
    value = length(data.equinix_fabric_connections.test.data)
}

output "first_connection_name" {
    value = data.equinix_fabric_connections.test.data.0.name
}

output "first_connection_uuid" {
    value = data.equinix_fabric_connections.test.data.0.uuid
}

output "first_connection_bandwidth" {
    value = data.equinix_fabric_connections.test.data.0.bandwidth
}

output "first_connection_type" {
    value = data.equinix_fabric_connections.test.data.0.type
}

output "first_connection_redundancy_priority" {
    value = one(data.equinix_fabric_connections.test.data.0.redundancy).priority
}

output "first_connection_purchase_order_number" {
    value = one(data.equinix_fabric_connections.test.data.0.order).purchase_order_number
}

output "first_connection_aSide_type" {
    value = one(one(data.equinix_fabric_connections.test.data.0.a_side).access_point).type
}

output "first_connection_aSide_link_protocol_type" {
    value = one(one(one(data.equinix_fabric_connections.test.data.0.a_side).access_point).link_protocol).type
}

output "first_connection_aSide_link_protocol_vlan_tag" {
    value = one(one(one(data.equinix_fabric_connections.test.data.0.a_side).access_point).link_protocol).vlan_tag
}

output "first_connection_aSide_location_metro_code" {
    value = one(one(one(data.equinix_fabric_connections.test.data.0.a_side).access_point).location).metro_code
}

output "first_connection_zSide_type" {
    value = one(one(data.equinix_fabric_connections.test.data.0.z_side).access_point).type
}

output "first_connection_zSide_link_protocol_type" {
    value = one(one(one(data.equinix_fabric_connections.test.data.0.z_side).access_point).link_protocol).type
}

output "first_connection_zSide_link_protocol_vlan_tag" {
    value = one(one(one(data.equinix_fabric_connections.test.data.0.z_side).access_point).link_protocol).vlan_tag
}

output "first_connection_zSide_location_metro_code" {
    value = one(one(one(data.equinix_fabric_connections.test.data.0.z_side).access_point).location).metro_code
}
