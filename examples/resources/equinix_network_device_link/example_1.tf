# Example of device link with HA device pair
# where each device is in different metro
resource "equinix_network_device_link" "test" {
  name   = "test-link"
  subnet = "192.168.40.64/27"
  project_id  = "a86d7112-d740-4758-9c9c-31e66373746b"
  device {
    id           = equinix_network_device.test.uuid
    asn          = equinix_network_device.test.asn > 0 ? equinix_network_device.test.asn : 22111
    interface_id = 6
  }
  device {
    id           = equinix_network_device.test.secondary_device[0].uuid
    asn          = equinix_network_device.test.secondary_device[0].asn > 0 ? equinix_network_device.test.secondary_device[0].asn : 22333
    interface_id = 7
  }
  link {
    account_number  = equinix_network_device.test.account_number
    src_metro_code  = equinix_network_device.test.metro_code
    dst_metro_code  = equinix_network_device.test.secondary_device[0].metro_code
    throughput      = "50"
    throughput_unit = "Mbps"
  }
}

