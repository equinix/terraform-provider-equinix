# Fetch a device data by hostname and show it's ID

data "equinix_metal_device" "test" {
  project_id = local.project_id
  hostname   = "mydevice"
}

output "id" {
  value = data.equinix_metal_device.test.id
}
