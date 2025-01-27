data "equinix_fabric_metros" "metros" {
  pagination = {}
}

output "id" {
  value = data.equinix_fabric_metros.metros.id
}

output "type" {
  value = data.equinix_fabric_metros.metros.data.0.type
}

output "metro_code" {
  value = data.equinix_fabric_metros.metros.data.0.code
}

output "region" {
  value = data.equinix_fabric_metros.metros.data.0.region
}

output "name" {
  value = data.equinix_fabric_metros.metros.data.0.name
}

output "equinix_asn" {
  value = data.equinix_fabric_metros.metros.data.0.equinix_asn
}

output "local_vc_bandwidth_max" {
  value = data.equinix_fabric_metros.metros.data.0.local_vc_bandwidth_max
}

output "geo_coordinates" {
  value = data.equinix_fabric_metros.metros.data.0.geo_coordinates
}

output "connected_metros" {
  value = data.equinix_fabric_metros.metros.data.0.connected_metros
}

output "geo_scopes" {
  value = data.equinix_fabric_metros.metros.data.0.geo_scopes
}

