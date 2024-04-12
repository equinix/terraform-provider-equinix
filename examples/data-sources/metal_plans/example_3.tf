# Following example will use equinix_metal_plans to select the cheapest plan available in metro 'sv' (Sillicon Valley)
data "equinix_metal_plans" "example" {
    sort {
        attribute = "pricing_hour"
        direction = "asc"
    }
    filter {
        attribute = "name"
        values    = ["c3.small.x86", "c3.medium.x86", "m3.large.x86"]
    }
    filter {
        attribute = "available_in_metros"
        values    = ["sv"]
    }
}

# This equinix_metal_device will use the first returned plan and the first metro in which that plan is available
# It will ignore future changes on plan and metro
resource "equinix_metal_device" "example" {
  hostname         = "example"
  plan             = data.equinix_metal_plans.example.plans[0].name
  metro            = data.equinix_metal_plans.example.plans[0].available_in_metros[0]
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = var.project_id

  lifecycle {
    ignore_changes = [
        plan,
        metro,
    ]
  }
}
