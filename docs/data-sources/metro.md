---
page_title: "Equinix Metal: metal_metro"
subcategory: ""
description: |-
  Provides an Equinix Metal metro datasource. This can be used to read metros.
---

# metal_metro

Provides an Equinix Metal metro datasource.

## Example Usage

```hcl
# Fetch a metro by code and show its ID

data "metal_metro" "sv" {
    code = "sv"
}

output "id" {
  value = data.metal_metro.sv.id
}
```

## Argument Reference

The following arguments are supported:

* `code` - The metro code

Metros can be looked up by `code`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the metro
* `code` - The code of the metro
* `country` - The country of the metro
* `name` - The name of the metro
