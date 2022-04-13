---
subcategory: "Network Edge"
---

# equinix_network_account (Data Source)

Use this data source to get number and identifier of Equinix Network Edge
billing account in a given metro location.

Billing account reference is required to create Network Edge virtual device
in corresponding metro location.

## Example Usage

```hcl
# Retrieve details of an account in Active status in DC metro
data "equinix_network_account" "dc" {
  metro_code = "DC"
  status     = "Active"
}

output "number" {
  value = data.equinix_network_account.dc.number
}
```

## Argument Reference

The following arguments are supported:

* `metro_code` - (Required) Account location metro code.
* `name` - (Optional) Account name for filtering.
* `status` - (Optional) Account status for filtering. Possible values are: `Active`, `Processing`,
`Submitted`, `Staged`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `number` - Account unique number.
* `ucm_id` - Account unique identifier.
