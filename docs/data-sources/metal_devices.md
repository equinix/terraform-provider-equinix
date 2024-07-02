---
subcategory: "Metal"
---

# equinix_metal_devices

The datasource can be used to find a list of devices which meet filter criteria.

If you need to fetch a single device by ID or by project ID and hostname, use the [equinix_metal_device](equinix_metal_device.md) datasource.

## Example Usage

```terraform
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
```

```terraform
# Following example takes advantage of the `search` field in the API request, and will select devices with
# string "database" in one of the searched attributes. See `search` in argument reference.
data "equinix_metal_devices" "example" {
    search = "database"
}

output "devices" {
    value = data.equinix_metal_devices.example.devices
}
```

## search vs filter

The difference between `search` and `filter` is that `search` is an API parameter, interpreted by the Equinix Metal service. The "filter" arguments will reduce the API list (or search) results by applying client-side filtering, within this provider.

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) ID of project containing the devices. Exactly one of `project_id` and `organization_id` must be set.
* `organization_id` - (Optional) ID of organization containing the devices.
* `search` - (Optional) - Search string to filter devices by hostname, description, short_id, reservation short_id, tags, plan name, plan slug, facility code, facility name, operating system name, operating system slug, IP addresses.
* `filter` - (Optional) One or more attribute/values pairs to filter. List of atributes to filter can be found in the [attribute reference](equinix_metal_device.md#attributes-reference) of the `equinix_metal_device` datasource.
  - `attribute` - (Required) The attribute used to filter. Filter attributes are case-sensitive
  - `values` - (Required) The filter values. Filter values are case-sensitive. If you specify multiple values for a filter, the values are joined with an OR by default, and the request returns all results that match any of the specified values
  - `match_by` - (Optional) The type of comparison to apply. One of: `in` , `re`, `substring`, `less_than`, `less_than_or_equal`, `greater_than`, `greater_than_or_equal`. Default is `in`.
  - `all` - (Optional) If is set to true, the values are joined with an AND, and the requests returns only the results that match all specified values. Default is `false`.

All fields in the `devices` block defined below can be used as attribute for both `sort` and `filter` blocks.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `devices` - list of resources with attributes like in the [equninix_metal_device datasources](equinix_metal_device.md).
