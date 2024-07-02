---
subcategory: "Metal"
---

# equinix_metal_project_ssh_key (Data Source)

Use this datasource to retrieve attributes of a Project SSH Key API resource.

## Example Usage

```terraform
# Get Project SSH Key by name
data "equinix_metal_project_ssh_key" "my_key" {
  search     = "username@hostname"
  project_id = local.project_id
}
```

## Argument Reference

The following arguments are supported:

* `search` - (Optional) The name, fingerprint, or public_key of the SSH Key to search for in the Equinix Metal project.
* `id` - (Optional) The id of the SSH Key to search for in the Equinix Metal project.
* `project_id` - (Optional) The Equinix Metal project id of the Equinix Metal SSH Key.

-> **NOTE:** One of either `search` or `id` must be provided along with `project_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the key.
* `name` - The name of the SSH key.
* `public_key` - The text of the public key.
* `project_id` - The ID of parent project.
* `owner_id` - The ID of parent project (same as project_id).
* `fingerprint` - The fingerprint of the SSH key.
* `created` - The timestamp for when the SSH key was created.
* `updated` - The timestamp for the last time the SSH key was updated.
