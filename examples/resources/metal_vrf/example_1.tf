resource "equinix_metal_project" "example" {
    name = "example"
}

resource "equinix_metal_vrf" "example" {
    description = "VRF with ASN 65000 and a pool of address space that includes 192.168.100.0/25"
    name        = "example-vrf"
    metro       = "da"
    local_asn   = "65000"
    ip_ranges   = ["192.168.100.0/25", "192.168.200.0/25"]
    project_id  = equinix_metal_project.example.id
}
