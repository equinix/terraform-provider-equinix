---
page_title: "Migrating devices from facilities to metros"
---

# Metros vs. Facilities

In April 2021, [Equinix Metal](https://metal.equinix.com/) rolled out a new location concept - [metros](https://feedback.equinixmetal.com/changelog/new-metros-feature-live). A metro is an Equinix-wide concept for data centers which are grouped together geographically. Data centers within a metro share capacity and networking features. You can read more about metros at https://metal.equinix.com/developers/docs/locations/metros/. In April 2023, [facility deployment was deprecated](https://feedback.equinixmetal.com/changelog/bye-facilities-hello-again-metros) with metro provisioning offered as an all-around superior deployment strategy. See the [facilities guide in the Metal API docs](https://deploy.equinix.com/developers/docs/metal/locations/facilities/) for more details about the API change.

Before the introduction of metros, resources were deployed to a single `facility` location.  When provisioning `equinix_metal_device` resources, the facility could be chosen by Equinix Metal with a user-supplied list of `facilities`, or a wildcard `any` facility.  The individual facility locations use a code like "sv15" or "ny5". Metros group facilities. For example, metro "sv" contains the "sv15" facility, among others. If you specify a metro when creating a resource, it will be deployed to one of the facilities in the metro group. You can then find the deployed facility using a read-only attribute of the resource (e.g. `deployed_facility` for `equinix_metal_device` resources).

## Changing your Terraform templates to use metros instead of facilities

To take advantage of some of the features of the metro, you might want to change the configuration of your Terraform templates so that the devices have `metro` specified instead of `facilities`. As both of the `metro` and `facilities` are ForceNew parameters (a change will trigger re-creation of the resource), you should be cautious if you don't want to have the device destroyed.

We updated the `equinix_metal_device` resource so that the change should be seamless, but please proceed with care. The `metro` parameter is also a computed attribute, and if you use newer provider version than 3.2.1, the `metro` attribute is actually present in your resource. You then only need to add it explicitly to your configuration.

The `facilities` parameter is only used for facility selection when creating the device resource. The actual facility where the device is deployed is in the `deployed_facility` Computed attribute.

If you only want to change the `equinix_metal_device` resource specfication from facility-based to metro-based, e.g. from facilities ["sv15"] to metro "sv", it's enough to remove the `facilities` attribute, and add the `metro` attribute. 

For example, given following configuration of a device deployed in the `sv15` facility:

```hcl-terraform
resource "equinix_metal_device" "node" {
  project_id       = local.project_id
  facilities       = ["sv15"]
  plan             = "c3.small.x86"
  operating_system = "ubuntu_16_04"
  hostname         = "test"
  billing_cycle    = "hourly"
}
```

.. you can remove `faclities` and add `metro`, changing the configuration to:


```hcl-terraform
resource "equinix_metal_device" "node" {
  project_id       = local.project_id
  metro            = "sv"
  plan             = "c3.small.x86"
  operating_system = "ubuntu_16_04"
  hostname         = "test"
  billing_cycle    = "hourly"
}
```

To test that the change didn't taint the state, and that the device will not be re-created, you can check if `terraform plan` reports any differences. The terraform state should be up to date as long as the facility in which the device was deployed was contained within the metro.

If the plan diff is not empty, you might have used a metro not containing the facility to which the device was deployed. This might happen if you've used more facilities in the `facilities` list, or you have used the special "any" facility.

You can find out the deployed facility, and the containing metro by examining the terraform state of the `equinix_metal_device` resource:

```
$ terraform state show equinix_metal_device.node | grep deployed
    deployed_facility                = "sv15"
```

```
$ terraform state show equinix_metal_device.node | grep metro
    metro                            = "sv"
```

You should then set the existing metro in your Terraform templates.
