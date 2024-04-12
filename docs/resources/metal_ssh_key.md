---
subcategory: "Metal"
---

# equinix_metal_ssh_key (Resource)

Provides a resource to manage User SSH keys on your Equinix Metal user account. If you create a new device in a project, all the keys of the project's collaborators will be injected to the device.

The link between User SSH key and device is implicit. If you want to make sure that a key will be copied to a device, you must ensure that the device resource `depends_on` the key resource.

## Example Usage

```terraform
# Create a new SSH key
resource "equinix_metal_ssh_key" "key1" {
  name       = "terraform-1"
  public_key = file("/home/terraform/.ssh/id_rsa.pub")
}

# Create new device with "key1" included. The device resource "depends_on" the
# key, in order to make sure the key is created before the device.
resource "equinix_metal_device" "test" {
  hostname         = "test-device"
  plan             = "c3.small.x86"
  metro            = "sv"
  operating_system = "ubuntu_20_04"
  billing_cycle    = "hourly"
  project_id       = local.project_id
  depends_on       = ["equinix_metal_ssh_key.key1"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSH key for identification
* `public_key` - (Required) The public key. If this is a file, it can be read using the file interpolation function

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the key.
* `name` - The name of the SSH key.
* `public_key` - The text of the public key.
* `fingerprint` - The fingerprint of the SSH key.
* `owner_id` - The UUID of the Equinix Metal API User who owns this key.
* `created` - The timestamp for when the SSH key was created.
* `updated` - The timestamp for the last time the SSH key was updated.

## Import

This resource can be imported using an existing SSH Key ID:

```sh
terraform import equinix_metal_ssh_key {existing_sshkey_id}
```
