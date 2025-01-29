data "equinix_fabric_metro" "metro" {
  metro_code =  "<metro_code>"
}

output "id" {
  value = data.equinix_fabric_metro.metro.id
}

output "type" {
  value = data.equinix_fabric_metro.metro.type
}

output "metro_code" {
  value = data.equinix_fabric_metro.metro.metro_code
}

output "region" {
  value = data.equinix_fabric_metro.metro.region
}

output "name" {
  value = data.equinix_fabric_metro.metro.name
}

output "equinix_asn" {
  value = data.equinix_fabric_metro.metro.equinix_asn
}

output "local_vc_bandwidth_max" {
  value = data.equinix_fabric_metro.metro.local_vc_bandwidth_max
}

output "geo_coordinates" {
  value = data.equinix_fabric_metro.metro.geo_coordinates
}

output "connected_metros" {
  value = data.equinix_fabric_metro.metro.connected_metros
}

output "geoScopes" {
  value = data.equinix_fabric_metro.metro.geo_scopes
}

