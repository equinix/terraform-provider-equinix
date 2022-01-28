---
page_title: "Equinix Metal: metal_project"
subcategory: ""
description: |-
  Provides an Equinix Metal Project datasource.
---

# metal\_project

Use this datasource to retrieve attributes of the Project API resource.

## Example Usage

```hcl
# Get Project by name and print UUIDs of its users
data "metal_project" "tf_project_1" {
  name = "Terraform Fun"
}

output "users_of_Terraform_Fun" {
  value = data.metal_project.tf_project_1.user_ids
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name which is used to look up the project
* `project_id` - The UUID by which to look up the project

## Attributes Reference

The following attributes are exported:

* `payment_method_id` - The UUID of payment method for this project
* `organization_id` - The UUID of this project's parent organization
* `backend_transfer` - Whether Backend Transfer is enabled for this project
* `created` - The timestamp for when the project was created
* `updated` - The timestamp for the last time the project was updated
* `user_ids` - List of UUIDs of user accounts which belong to this project
* `bgp_config` - Optional BGP settings. Refer to [Equinix Metal guide for BGP](https://metal.equinix.com/developers/docs/networking/local-global-bgp/).

The `bgp_config` block contains:

* `asn` - Autonomous System Number for local BGP deployment
* `md5` - Password for BGP session in plaintext (not a checksum)
* `deployment_type` - `private` or `public`, the `private` is likely to be usable immediately, the `public` will need to be review by Equinix Metal engineers
* `status` - status of BGP configuration in the project
* `max_prefix` - The maximum number of route filters allowed per server
