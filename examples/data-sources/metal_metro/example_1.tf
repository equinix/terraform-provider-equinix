# Fetch a metro by code and show its ID

data "equinix_metal_metro" "sv" {
  code = "sv"
}

output "id" {
  value = data.equinix_metal_metro.sv.id
}
