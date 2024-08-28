resource "equinix_metal_reserved_ip_block" "example" {
    description = "Reserved IP block (192.168.100.0/29) taken from on of the ranges in the VRF's pool of address space."
    project_id  = equinix_metal_project.example.id
    metro       = equinix_metal_vrf.example.metro
    type        = "vrf"
    vrf_id      = equinix_metal_vrf.example.id
    cidr        = 29
    network     = "192.168.100.0"
}

resource "equinix_metal_vlan" "example" {
    description = "A VLAN for Layer2 and Hybrid Metal devices"
    metro       = equinix_metal_vrf.example.metro
    project_id  = equinix_metal_project.example.id
}

resource "equinix_metal_gateway" "example" {
    project_id        = equinix_metal_project.example.id
    vlan_id           = equinix_metal_vlan.example.id
    ip_reservation_id = equinix_metal_reserved_ip_block.example.id
}
