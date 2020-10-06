---
layout: "equinix"
page_title: "Equinix: ne_account"
sidebar_current: "docs-equinix-datasource-ne-account"
description: |-
 Get information on Network Edge billing account
---

# Data Source: ne_account

Use this data source to get number and identifier of Network Edge
billing account in a given metro location.

Billing account reference is required to create Network Edge virtual device
in corresponding metro location.

## Example Usage

```hcl
# Retrieve details of an account in Active status in DC metro
data "equinix_ne_account" "dc" {
  metro_code = "DC"
  status     = "Active"
}
```

## Argument Reference

* `metro_code` - (Required) Account metro code
* `name` - (Optional) Account name for filtering
* `status` - (Optional) Account status for filtering. Possible values are "Active",
"Processing", "Submitted", "Staged"

## Attributes Reference

* `number` - Account number
* `ucm_id` - Account unique identifier
