# Retrieve data for an existing Equinix Network Edge device with UUID "f0b5c553-cdeb-4bc3-95b8-23db9ccfd5ee"
data "equinix_network_device" "by_uuid" {
  uuid = "f0b5c553-cdeb-4bc3-95b8-23db9ccfd5ee"
}

# Retrieve data for an existing Equinix Network Edge device named "Arcus-Gateway-A1"
data "equinix_network_device" "by_name" {
  name = "Arcus-Gateway-A1"
}
