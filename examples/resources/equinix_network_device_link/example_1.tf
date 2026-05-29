# Example of device link with HA device pair
# where each device is in a different metro
resource "equinix_network_device_link" "test" {
  name       = "test-DLG"
  project_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  device {
    id           = equinix_network_device.test.uuid
    interface_id = 6
  }
  device {
    id           = equinix_network_device.test.secondary_device[0].uuid
    interface_id = 7
  }
  metro_link {
    account_number  = equinix_network_device.test.account_number
    metro_code      = equinix_network_device.test.metro_code
    throughput      = "50"
    throughput_unit = "Mbps"
  }
  metro_link {
    account_number  = equinix_network_device.test.secondary_device[0].account_number
    metro_code      = equinix_network_device.test.secondary_device[0].metro_code
    throughput      = "50"
    throughput_unit = "Mbps"
  }
}

