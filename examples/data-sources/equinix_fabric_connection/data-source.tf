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
  value = [for account in data.equinix_fabric_connection.connection_data_name.account: account.account_number]
}

output "bandwidth" {
  value = data.equinix_fabric_connection.connection_data_name.bandwidth
}

output "project_id" {
  value = [for project in data.equinix_fabric_connection.connection_data_name.project: project.project_id]
}

output "redundancy_group" {
  value = [for redundancy in data.equinix_fabric_connection.connection_data_name.redundancy: redundancy.group]
}

output "redundancy_priority" {
  value = [for redundancy in data.equinix_fabric_connection.connection_data_name.redundancy: redundancy.priority]
}

output "state" {
  value = data.equinix_fabric_connection.connection_data_name.state
}

output "type" {
  value = data.equinix_fabric_connection.connection_data_name.type
}

# Same for z_side just use z_side instead of a_side
output "access_point_type" {
  value = [for aside in data.equinix_fabric_connection.connection_data_name.a_side:
    [for access in aside.access_point: access.type]]
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_type" {
  value = [for aside in data.equinix_fabric_connection.connection_data_name.a_side:
    [for access in aside.access_point:
      [for protocol in access.link_protocol: protocol.type]]]
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_vlan_tag" {
  value = [for aside in data.equinix_fabric_connection.connection_data_name.a_side:
    [for access in aside.access_point:
      [for protocol in access.link_protocol: protocol.vlan_tag]]]
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_vlan_c_tag" {
  value = [for aside in data.equinix_fabric_connection.connection_data_name.a_side:
    [for access in aside.access_point:
      [for protocol in access.link_protocol: protocol.vlan_c_tag]]]
}

# Same for z_side just use z_side instead of a_side
output "access_point_link_protocol_vlan_s_tag" {
  value = [for aside in data.equinix_fabric_connection.connection_data_name.a_side:
    [for access in aside.access_point:
      [for protocol in access.link_protocol: protocol.vlan_s_tag]]]
}

# Same for z_side just use z_side instead of a_side
output "access_point_provider_connection_id" {
  value = [for aside in data.equinix_fabric_connection.connection_data_name.a_side:
    [for access in aside.access_point: access.provider_connection_id]]
}

