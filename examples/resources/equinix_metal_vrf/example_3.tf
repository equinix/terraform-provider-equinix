data "equinix_metal_connection" "example" {
    connection_id = var.metal_dedicated_connection_id
}

resource "equinix_metal_virtual_circuit" "example" {
    name          = "example-vc"
    description   = "Virtual Circuit"
    connection_id = data.equinix_metal_connection.example.id
    project_id    = equinix_metal_project.example.id
    port_id       = data.equinix_metal_connection.example.ports[0].id
    nni_vlan      = 1024
    vrf_id        = equinix_metal_vrf.example.id
    peer_asn      = 65530
    subnet        = "192.168.100.16/31"
    metal_ip      = "192.168.100.16"
    customer_ip   = "192.168.100.17"
}
