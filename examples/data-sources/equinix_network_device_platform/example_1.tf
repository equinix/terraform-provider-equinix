# Retrieve platform configuration of a large flavor for a CSR100V device type
# Platform has to support IPBASE software package
data "equinix_network_device_platform" "csrLarge" {
  device_type = "CSR1000V"
  flavor      = "large"
  packages    = ["IPBASE"]
}
