# Retrieve device type details of a Cisco router
# Device type has to be available in DC and SV metros
data "equinix_network_device_type" "ciscoRouter" {
  category    = "Router"
  vendor      = "Cisco"
  metro_codes = ["DC", "SV"]
}
