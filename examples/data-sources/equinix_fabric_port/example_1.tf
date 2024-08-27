data "equinix_fabric_port" "port_data_name" {
  uuid = "<uuid_of_port>"
}

output "id" {
  value = data.equinix_fabric_port.port_data_name.id
}

output "name" {
  value = data.equinix_fabric_port.port_data_name.name
}

output "state" {
  value = data.equinix_fabric_port.port_data_name.state
}

output "account_name" {
  value = data.equinix_fabric_port.port_data_name.account.0.account_name
}

output "type" {
  value = data.equinix_fabric_port.port_data_name.type
}

output "bandwidth" {
  value = data.equinix_fabric_port.port_data_name.bandwidth
}

output "used_bandwidth" {
  value = data.equinix_fabric_port.port_data_name.used_bandwidth
}

output "encapsulation_type" {
  value = data.equinix_fabric_port.port_data_name.encapsulation.0.type
}

output "ibx" {
  value = data.equinix_fabric_port.port_data_name.location.0.ibx
}

output "metro_code" {
  value = data.equinix_fabric_port.port_data_name.location.0.metro_code
}

output "metro_name" {
  value = data.equinix_fabric_port.port_data_name.location.0.metro_name
}

output "region" {
  value = data.equinix_fabric_port.port_data_name.location.0.region
}

output "device_redundancy_enabled" {
  value = data.equinix_fabric_port.port_data_name.device.0.redundancy.0.enabled
}

output "device_redundancy_priority" {
  value = data.equinix_fabric_port.port_data_name.device.0.redundancy.0.priority
}
