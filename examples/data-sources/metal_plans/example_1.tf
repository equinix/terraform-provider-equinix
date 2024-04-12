# Following example will select device plans which are under 2.5$ per hour, are available in metro 'da' (Dallas)
# OR 'sv' (Sillicon Valley) and sort it by the hourly price ascending.
data "equinix_metal_plans" "example" {
    sort {
        attribute = "pricing_hour"
        direction = "asc"
    }
    filter {
        attribute = "pricing_hour"
        values    = [2.5]
        match_by  = "less_than"
    }
    filter {
        attribute = "available_in_metros"
        values    = ["da", "sv"]
    }
}

output "plans" {
    value = data.equinix_metal_plans.example.plans
}
