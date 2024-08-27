# Fetch a facility by code and show its ID

data "equinix_metal_facility" "ny5" {
  code = "ny5"
}

output "id" {
  value = data.equinix_metal_facility.ny5.id
}
