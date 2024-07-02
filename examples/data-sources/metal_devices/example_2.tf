# Following example takes advantage of the `search` field in the API request, and will select devices with
# string "database" in one of the searched attributes. See `search` in argument reference.
data "equinix_metal_devices" "example" {
    search = "database"
}

output "devices" {
    value = data.equinix_metal_devices.example.devices
}
