// lookup by ID
data "equinix_metal_hardware_reservation" "example" {
  id = "4347e805-eb46-4699-9eb9-5c116e6a0172"
}

// lookup by device ID
data "equinix_metal_hardware_reservation" "example_by_device_id" {
  device_id = "ff85aa58-c106-4624-8f1c-7c64554047ea"
}
