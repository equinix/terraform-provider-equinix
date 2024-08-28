data "equinix_fabric_network" "network_data_name" {
  uuid = "<uuid_of_network>"
}

output "id" {
  value = data.equinix_fabric_network.network_data_name.id
}

output "name" {
  value = data.equinix_fabric_network.network_data_name.name
}

output "scope" {
  value = data.equinix_fabric_network.network_data_name.scope
}

output "type" {
  value = data.equinix_fabric_network.network_data_name.type
}

output "region" {
  value = data.equinix_fabric_network.network_data_name.location.0.region
}
