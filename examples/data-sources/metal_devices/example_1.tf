# Following example will select c3.small.x86 devices which are deplyed in metro 'da' (Dallas)
# OR 'sv' (Sillicon Valley).
data "equinix_metal_devices" "example" {
    project_id = local.project_id
    filter {
        attribute = "plan"
        values    = ["c3.small.x86"]
    }
    filter {
        attribute = "metro"
        values    = ["da", "sv"]
    }
}

output "devices" {
    organization_id = local.org_id
    value = data.equinix_metal_devices.example.devices
}
