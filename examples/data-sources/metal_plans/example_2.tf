# Following example will select device plans with class containing string 'large', are available in metro 'da' (Dallas)
# AND 'sv' (Sillicon Valley), are elegible for spot_market deployments.
data "equinix_metal_plans" "example" {
    filter {
        attribute = "class"
        values    = ["large"]
        match_by  = "substring"
    }
    filter {
        attribute = "deployment_types"
        values    = ["spot_market"]
    }
    filter {
        attribute = "available_in_metros"
        values    = ["da", "sv"]
        all       = true
    }
}

output "plans" {
    value = data.equinix_metal_plans.example.plans
}
