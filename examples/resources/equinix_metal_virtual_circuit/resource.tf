locals {
  project_id = "52000fb2-ee46-4673-93a8-de2c2bdba33c"
  conn_id = "73f12f29-3e19-43a0-8e90-ae81580db1e0"
}

data "equinix_metal_connection" test {
  connection_id = local.conn_id
}

resource "equinix_metal_vlan" "test" {
  project_id = local.project_id
  metro      = data.equinix_metal_connection.test.metro
}

resource "equinix_metal_virtual_circuit" "test" {
  connection_id = local.conn_id
  project_id = local.project_id
  port_id = data.equinix_metal_connection.test.ports[0].id
  vlan_id = equinix_metal_vlan.test.id
  nni_vlan = 1056
}
