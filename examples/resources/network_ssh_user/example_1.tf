# Create SSH user with password auth method and associate it with
# two virtual network devices

resource "equinix_network_ssh_user" "john" {
  username = "john"
  password = "secret"
  device_ids = [
    equinix_network_device.csr1000v-ha.uuid,
    equinix_network_device.csr1000v-ha.redundant_uuid
  ]
}
