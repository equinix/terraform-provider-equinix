data "equinix_fabric_connection" "connection_data_name" {
  uuid = "<uuid_of_connection>"
}

output "id" {
  value = data.equinix_fabric_connection.connection_data_name.id
}

output "name" {
  value = data.equinix_fabric_connection.connection_data_name.bandwidth
}

output "account_number" {
  value = data.equinix_fabric_connection.connection_data_name.account.0.account_number
}

output "bandwidth" {
  value = data.equinix_fabric_connection.connection_data_name.bandwidth
}

output "project_id" {
  value = data.equinix_fabric_connection.connection_data_name.project.0.project_id
}

output "redundancy_group" {
  value = data.equinix_fabric_connection.connection_data_name.redundancy.0.group
}

output "redundancy_priority" {
  value = data.equinix_fabric_connection.connection_data_name.redundancy.0.priority
}

output "state" {
  value = data.equinix_fabric_connection.connection_data_name.state
}

output "type" {
  value = data.equinix_fabric_connection.connection_data_name.type
}

# Same for z_side just use z_side instead of a_side
output "access_point_type" {
  value = data.equinix_fabric_connection.connection_data_name.a_side.0.access_point.0.type
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_type" {
  value = data.equinix_fabric_connection.connection_data_name.a_side.0.access_point.0.link_protocol.0.type
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_vlan_tag" {
  value = data.equinix_fabric_connection.connection_data_name.a_side.0.access_point.0.link_protocol.0.vlan_tag
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_vlan_c_tag" {
  value = data.equinix_fabric_connection.connection_data_name.a_side.0.access_point.0.link_protocol.0.vlan_c_tag
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_vlan_s_tag" {
  value = data.equinix_fabric_connection.connection_data_name.a_side.0.access_point.0.link_protocol.0.vlan_s_tag
}

# Same for z_side just use z_side instead of a_side
output "access_point_provider_connection_id" {
  value = data.equinix_fabric_connection.connection_data_name.a_side.0.access_point.0.provider_connection_id
}

