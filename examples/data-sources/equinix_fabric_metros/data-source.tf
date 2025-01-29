data "equinix_fabric_metros" "metros" {
  pagination = {
    limit = 12,
    offset = 6
  }
}

output "number_of_returned_metros" {
  value = length(data.equinix_fabric_metros.metros.data)
}

output "first_metro_id" {
  value = data.equinix_fabric_metros.metros.data.0.id
}

output "first_metro_type" {
  value = data.equinix_fabric_metros.metros.data.0.type
}

output "first_metro_code" {
  value = data.equinix_fabric_metros.metros.data.0.code
}

output "first_metro_region" {
  value = data.equinix_fabric_metros.metros.data.0.region
}

output "first_metro_name" {
  value = data.equinix_fabric_metros.metros.data.0.name
}

output "first_metro_equinix_asn" {
  value = data.equinix_fabric_metros.metros.data.0.equinix_asn
}

output "first_metro_local_vc_bandwidth_max" {
  value = data.equinix_fabric_metros.metros.data.0.local_vc_bandwidth_max
}

output "first_metro_geo_coordinates" {
  value = data.equinix_fabric_metros.metros.data.0.geo_coordinates
}

output "first_metro_connected_metros" {
  value = data.equinix_fabric_metros.metros.data.0.connected_metros
}

output "first_metro_geo_scopes" {
  value = data.equinix_fabric_metros.metros.data.0.geo_scopes
}

