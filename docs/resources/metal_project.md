---
subcategory: "Metal"
---

# equinix_metal_project (Resource)

Provides an Equinix Metal project resource to allow you manage devices in your projects.

-> **NOTE:** Keep in mind that Equinix Metal invoicing is per project, so creating many `equinix_metal_project` resources will affect the rendered invoice. If you want to keep your Equinix Metal bill simple and easy to review, please re-use your existing projects.

## Example Usage

### Create a new project

```terraform
resource "equinix_metal_project" "tf_project_1" {
  name = "Terraform Fun"
}
```

### Example with BGP config

```terraform
# Create a new Project
resource "equinix_metal_project" "tf_project_1" {
  name = "tftest"
  bgp_config {
    deployment_type = "local"
    md5             = "C179c28c41a85b"
    asn             = 65000
  }
}
```

### Enabling BGP in an existing project

If you want to enable BGP in an existing Equinix Metal project, you should first create a resource in your TF config for the existing projects. Set your BGP configuration.

```terraform
resource "equinix_metal_project" "existing_project" {
  name = "The name of the project (if different, will rewrite)"
  bgp_config {
    deployment_type = "local"
    md5             = "C179c28c41a85b"
    asn             = 65000
  }
}
```

Then, find out the UUID of the existing project, and import it to your TF state.

```sh
terraform import equinix_metal_project.existing_project e188d7db-46a7-46cb-8969-e63ec22695d5
```

Your existing project is now loaded in your local TF state, and linked to the resource with given name.

After running `terraform apply`, the project will be updated with configuration provided in the TF template.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project. The maximum length is 80 characters
* `organization_id` - (Required) The UUID of organization under which you want to create the project. If you leave it out, the project will be created under your the default organization of your account.
* `payment_method_id` - The UUID of payment method for this project. The payment method and the project need to belong to the same organization (passed with `organization_id`, or default).
* `backend_transfer` - Enable or disable [Backend Transfer](https://metal.equinix.com/developers/docs/networking/backend-transfer/), default is `false`.
* `bgp_config` - Optional BGP settings. Refer to [Equinix Metal guide for BGP](https://metal.equinix.com/developers/docs/networking/local-global-bgp/).

-> **NOTE:** Once you set the BGP config in a project, it can't be removed (due to a limitation in the Equinix Metal API). It can be updated.

The `bgp_config` block supports:

* `asn` - (Required) Autonomous System Number for local BGP deployment.
* `deployment_type` - (Required) `local` or `global`, the `local` is likely to be usable immediately, the `global` will need to be reviewed by Equinix Metal engineers.
* `md5` - (Optional) Password for BGP session in plaintext (not a checksum).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the project.
* `created` - The timestamp for when the project was created.
* `updated` - The timestamp for the last time the project was updated.

The `bgp_config` block additionally exports:

* `status` - status of BGP configuration in the project.
* `max_prefix` - The maximum number of route filters allowed per server.

## Import

This resource can be imported using an existing project ID:

```sh
terraform import equinix_metal_project {existing_project_id}
```
