# Get Project by name and print UUIDs of its users

data "equinix_metal_device_bgp_neighbors" "test" {
  device_id = "4c641195-25e5-4c3c-b2b7-4cd7a42c7b40"
}

output "bgp_neighbors_listing" {
  value = data.equinix_metal_device_bgp_neighbors.test.bgp_neighbors
}
