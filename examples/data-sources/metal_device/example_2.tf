# Fetch a device data by ID and show its public IPv4
data "equinix_metal_device" "test" {
  device_id = "4c641195-25e5-4c3c-b2b7-4cd7a42c7b40"
}

output "ipv4" {
  value = data.equinix_metal_device.test.access_public_ipv4
}
