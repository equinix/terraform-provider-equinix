---
subcategory: "Metal"
---

# equinix_metal_project_ssh_key (Resource)

Provides an Equinix Metal project SSH key resource to manage project-specific SSH keys. Project SSH keys will only be populated onto servers that belong to that project, in contrast to User SSH Keys.

## Example Usage

```terraform
locals {
  project_id = "<UUID_of_your_project>"
}

resource "equinix_metal_project_ssh_key" "test" {
  name       = "test"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDM/unxJeFqxsTJcu6mhqsMHSaVlpu+Jj/P+44zrm6X/MAoHSX3X9oLgujEjjZ74yLfdfe0bJrbL2YgJzNaEkIQQ1VPMHB5EhTKUBGnzlPP0hHTnxsjAm9qDHgUPgvgFDQSAMzdJRJ0Cexo16Ph9VxCoLh3dxiE7s2gaM2FdVg7P8aSxKypsxAhYV3D0AwqzoOyT6WWhBoQ0xZ85XevOTnJCpImSemEGs6nVGEsWcEc1d1YvdxFjAK4SdsKUMkj4Dsy/leKsdi/DEAf356vbMT1UHsXXvy5TlHu/Pa6qF53v32Enz+nhKy7/8W2Yt2yWx8HnQcT2rug9lvCXagJO6oauqRTO77C4QZn13ZLMZgLT66S/tNh2EX0gi6vmIs5dth8uF+K6nxIyKJXbcA4ASg7F1OJrHKFZdTc5v1cPeq6PcbqGgc+8SrPYQmzvQqLoMBuxyos2hUkYOmw3aeWJj9nFa8Wu5WaN89mUeOqSkU4S5cgUzWUOmKey56B/j/s1sVys9rMhZapVs0wL4L9GBBM48N5jAQZnnpo85A8KsZq5ME22bTLqnxsDXqDYZvS7PSI6Dxi7eleOFE/NYYDkrgDLHTQri8ucDMVeVWHgoMY2bPXdn7KKy5jW5jKsf8EPARXg77A4gRYmgKrcwIKqJEUPqyxJBe0CPoGTqgXPRsUiQ== tomk@hp2"
  project_id = local.project_id
}

resource "equinix_metal_device" "test" {
  hostname            = "test"
  plan                = "c3.medium.x86"
  metro               = "ny"
  operating_system    = "ubuntu_20_04"
  billing_cycle       = "hourly"
  project_ssh_key_ids = [equinix_metal_project_ssh_key.test.id]
  project_id          = local.project_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSH key for identification.
* `public_key` - (Required) The public key. If this is a file, it can be read using the file interpolation function.
* `project_id` - (Required) The ID of parent project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the key.
* `name` - The name of the SSH key.
* `owner_id` - The ID of parent project (same as project_id).
* `fingerprint` - The fingerprint of the SSH key.
* `created` - The timestamp for when the SSH key was created.
* `updated` - The timestamp for the last time the SSH key was updated.
