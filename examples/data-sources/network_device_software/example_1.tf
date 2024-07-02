# Retrieve details for CSR1000V device software with latest path of 16.09 version
# that supports IPBASE package
data "equinix_network_device_software" "csrLatest1609" {
  device_type   = "CSR1000V"
  version_regex = "^16.09.+"
  packages      = ["IPBASE"]
  most_recent   = true
}
