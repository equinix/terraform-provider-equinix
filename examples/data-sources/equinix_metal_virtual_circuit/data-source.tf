data "equinix_metal_connection" "example_connection" {
  connection_id = "4347e805-eb46-4699-9eb9-5c116e6a017d"
}

data "equinix_metal_virtual_circuit" "example_vc" {
  virtual_circuit_id = data.equinix_metal_connection.example_connection.ports[1].virtual_circuit_ids[0]
}

