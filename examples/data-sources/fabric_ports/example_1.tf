data "equinix_fabric_ports" "ports_data_name" {
  filters {
    name = "<name_of_port||port_prefix>"
  }
}


output "id" {
  value = data.equinix_fabric_port.ports_data_name.data.0.id
}

output "name" {
  value = data.equinix_fabric_port.ports_data_name.data.0.name
}

output "state" {
  value = data.equinix_fabric_port.ports_data_name.data.0.state
}

output "account_name" {
  value = data.equinix_fabric_port.ports_data_name.data.0.account.0.account_name
}

output "type" {
  value = data.equinix_fabric_port.ports_data_name.data.0.type
}

output "bandwidth" {
  value = data.equinix_fabric_port.ports_data_name.data.0.bandwidth
}

output "used_bandwidth" {
  value = data.equinix_fabric_port.ports_data_name.data.0.used_bandwidth
}

output "encapsulation_type" {
  value = data.equinix_fabric_port.ports_data_name.data.0.encapsulation.0.type
}

output "ibx" {
  value = data.equinix_fabric_port.ports_data_name.data.0.location.0.ibx
}

output "metro_code" {
  value = data.equinix_fabric_port.ports_data_name.data.0.location.0.metro_code
}

output "metro_name" {
  value = data.equinix_fabric_port.ports_data_name.data.0.location.0.metro_name
}

output "region" {
  value = data.equinix_fabric_port.ports_data_name.data.0.location.0.region
}

output "device_redundancy_enabled" {
  value = data.equinix_fabric_port.ports_data_name.data.0.device.0.redundancy.0.enabled
}

output "device_redundancy_priority" {
  value = data.equinix_fabric_port.ports_data_name.data.0.device.0.redundancy.0.priority
}
