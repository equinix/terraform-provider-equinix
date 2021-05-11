---
page_title: "Migrating devices from facilities to metros"
description: |-
  Migrating your device resources from facilities to metros
---

# Metros vs. Facilities

In April 2021, Equinix Metal rolled out a new location concept - [metros](https://feedback.equinixmetal.com/changelog/new-metros-feature-live). A metro is an Equinix-wide concept for data centers which are grouped together geographically. Data centers within a metro share capacity and networking features. You can read more about metros at https://metal.equinix.com/developers/docs/locations/metros/.

Until the metros introduction, resource deployment location used to be controlled by "facilty" - single data center with location code like "sv15" or "ny5". Metros group the facilities, and e.g. metro "sv" contains facility "sv15". If you specify a metro when creating a resource, it will (sometimes?) be deployed to one of the facilities in the metro group. You can then (sometimes?) find the deployed facility from a read-only attribute of the resource.


## Changing your Terraform templates to use metros instead of facilities

To take advantage of some of the features of the metro, you might want to change the configuration of your Terraform templates so that the devices have `metro` specified instead of `facility`. As both of the `metro` and `facility` are ForceNew paramters (a change will trigger re-creation of the resource), you should be cautious if you don't want to have the device destroyed. We updated the `metal_device` resource so that the change should be seamless, but please proceed with care.

If you only want to change the `metal_device` resource specfication from a facility to a metro containing the facility, e.g. from facility "sv15" to metro "sv", it's enough to remove the `facilities` attribute, and add the `metro` attribute. 

The `metal_device` resource has the facilities input attributes.

Given following configuration of a device deployed in the `sv15` facility.

```hcl-terraform
resource "metal_device" "node" {
  project_id       = local.project_id
  facilities       = ["sv15"]
  plan             = "c3.small.x86"
  operating_system = "ubuntu_16_04"
  hostname         = "test"
  billing_cycle    = "hourly"
}
```
