---
subcategory: "Metal"
---

# equinix_metal_plans

Provides an Equinix Metal plans datasource. This can be used to find plans that meet a filter criteria.

## Example Usage

```hcl
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
```

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `sort` - (Optional) One or more attribute/direction pairs on which to sort results. If multiple
sorts are provided, they will be applied in order
  - `attribute` - (Required) The attribute used to sort the results. Sort attributes are case-sensitive
  - `direction` - (Optional) Sort results in ascending or descending order. Strings are sorted in alphabetical order. One of: asc, desc
* `filter` - (Optional) One or more attribute/values pairs to filter off of
  - `attribute` - (Required) The attribute used to filter. Filter attributes are case-sensitive
  - `values` - (Required) The filter values. Filter values are case-sensitive. If you specify multiple values for a filter, the values are joined with an OR by default, and the request returns all results that match any of the specified values
  - `match_by` - (Optional) The type of comparison to apply. One of: `in` , `re`, `substring`, `less_than`, `less_than_or_equal`, `greater_than`, `greater_than_or_equal`. Default is `in`.
  - `all` - (Optional) If is set to true, the values are joined with an AND, and the requests returns only the results that match all specified values. Default is `false`.

All fields in the `plans` block defined below can be used as attribute for both `sort` and `filter` blocks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plans` - The ID of the facility
  - `id` - id of the plan
  - `name` - name of the plan
  - `slug`- plan slug
  - `description`- description of the plan
  - `line`- plan line, e.g. baremetal
  - `legacy`- flag showing if it's a legacy plan
  - `class`- plan class
  - `pricing_hour`- plan hourly price
  - `pricing_month`- plan monthly price
  - `deployment_types`- list of deployment types, e.g. on_demand, spot_market
  - `available_in`- list of facilities where the plan is available
  - `available_in_metros`- list of facilities where the plan is available