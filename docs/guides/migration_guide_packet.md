---
page_title: "Migrating from the Packet provider"
description: |-
  Migrating your templates from packet_ resources to metal_ resources with minimal disruption.
---

# Migrating from packethost/packet to equinix/metal

Packet is now Equinix Metal, and the name of the Terraform provider changed too. This (terraform-provider-metal, provider equinix/metal) is the current provider for Equinix Metal.

If you've been using terraform-provider-packet, and you want to use a newer provider version to manage resources in Equinix Metal, you will need to change the references in you HCL files. You can just change the names of the resources, e.g. from `packet_device` to `metal_device`. That should work, but it will cause the `packet_device` to be destroyed and new `metal_device` to be created instead. Re-creation of the resources might be undesirable, and this guide shows how to migrate to metal_ resources without the re-creation.

Before starting to migrate your Terraform templates, please upgrade
* packethost/packet provider to the latest version (3.2.1)
* Terraform to version at least v0.13


## Fast migration with replace-provider and sed

Just like the Terraform HCL templates, the Terraform state is a file containing resource names and their attributes in structured text. We can attempt the migration as a text substitution task, basically replacing `packet_` with `metal_` wherever possible, and fixing the provider source reference.

It's a good idea to make a backup of the whole Terraform directory before doing this.

Considering we have infrastructure created from following template:

```hcl-terraform
terraform {
  required_providers {
    packet = {
      source = "packethost/packet"
    }
  }
}

resource "packet_project" "example" {
  name = "example"
}

resource "packet_vlan" "example" {
  project_id       = packet_project.example.id
  facility         = "sv15"
  description      = "example"
}
```

We can first change the provider in Terraform state file (`terraform.tfstate`) with `terraform state` subcommand `replace-provider`:

```
$ terraform state replace-provider packethost/packet equinix/metal
```

Then we replace the provider reference in the HCL templates. Do this for every file where you have the reference:

```
$ sed -i 's|packethost/packet|equinix/metal|g' main.tf
``` 

Then we simply replace all strings `packet` with `metal` in the Terraform HCL files.

```
$ sed -i 's/packet/metal/g' main.tf
```

..this is a bit dangerous, so check your `git diff` after. It should replace all the `packet_` prefices and also the key from the `required_providers` block.

Then replace `packet_` with `metal_` in the terraform state file:

```
$ sed -i 's/packet_/metal_/g' terraform.tfstate
```

The example template would now look as:

```hcl-terraform
terraform {
  required_providers {
    metal = {
      source = "equinix/metal"
    }
  }
}

resource "metal_project" "example" {
  name = "example"
}

resource "metal_vlan" "example" {
  project_id       = metal_project.example.id
  facility         = "sv15"
  description      = "example"
}
```

We then need to install the equinix/metal provider by running `$ terraform init`. After that, our templates should be in check with the Terraform state and with the upstream resources in Equinix Metal. You can verify the result by running `$ terraform plan`.

If the plan is not empty, it means that some resources can not be simply read fom upstream, or that attributes have changed between your version of the packethost/packet provider and the current version of the equinix/metal provider. 

## Migrating one resource at a time

We can use `$ terraform state` and `$ terraform import` to achieve transition without destroying existing resources.

### Existing infrastructure

We assume to have infrastructure created with provider packethost/packet with a device and an IP reservation. The HCL looks like:

```hcl-terraform
terraform {
  required_providers {
    packet = {
      source = "packethost/packet"
      version = "3.2.1"
    }
  }
}

resource "packet_reserved_ip_block" "example" {
  project_id = local.project_id
  facility   = "sv15"
  quantity   = 2
}

resource "packet_device" "example" {
  project_id       = local.project_id
  facilities       = ["sv15"]
  plan             = "c3.medium.x86"
  operating_system = "ubuntu_20_04"
  hostname         = "test"
  billing_cycle    = "hourly"

  ip_address {
    type            = "public_ipv4"
    cidr            = 31
    reservation_ids = [packet_reserved_ip_block.example.id]
  }

  ip_address {
    type = "private_ipv4"
  }
}
```

### Resource UUIDs

In order to transition to provider equinix/metal, we need to find out UUIDs of all the resources we want to migrate. In this case `packet_reserved_ip_block.example` and `packet_device.example`. We can use `terraform state` to find out the UUIDs.

For the reserved IP block:

```
$ terraform state show packet_reserved_ip_block.example
# packet_reserved_ip_block.example:
resource "packet_reserved_ip_block" "example" {
	[...]
    id = "e689072f-aa6e-4d51-8e37-c2fbe18b4ff0"
	[...]
}
```

For the device:

```
$ terraform state show packet_device.example
# packet_device.example
resource "packet_device" "example" {
	[...]
    id = "8eb3bc10-0e1a-476a-aec2-6dc699df9c1c"
	[...]

```

### Migrated template

Once we find out the UUIDs of resources to migrate, in the HCL template, we need to change:
 
* the required_providers block to require equinix/metal
* the names of the resources to corresponding resoruces from provider equinix/metal (sed 's/packet_/metal_')
* all the references from packet_ resources to metal_ resources

The modified template will then look as:

```hcl-terraform
terraform {
  required_providers {
    metal = {
      source = "equinix/metal"
      version = "2.0.1"
    }
  }
}

resource "metal_reserved_ip_block" "example" {
  project_id = local.project_id
  facility   = "sv15"
  quantity   = 2
}

resource "metal_device" "example" {
  project_id       = local.project_id
  facilities       = ["sv15"]
  plan             = "c3.medium.x86"
  operating_system = "ubuntu_20_04"
  hostname         = "test"
  billing_cycle    = "hourly"

  ip_address {
    type            = "public_ipv4"
    cidr            = 31
    reservation_ids = [metal_reserved_ip_block.example.id]
  }

  ip_address {
    type = "private_ipv4"
  }
}
```

### Migrating Terraform state

Once we changed the template accordingly, we can remove the old packet_ resources from Terraform state and import the new ones as metal_ resources by their UUIDs.

From checking the state before, we remember that UUID of the packet_device.example is 8eb3bc10-0e1a-476a-aec2-6dc699df9c1c, and UUID of the packet_reserved_ip_block.example is e689072f-aa6e-4d51-8e37-c2fbe18b4ff0.

 In the terraform state and import commands, we use the resource type and name, separated by dot:

```
$ terraform state rm packet_reserved_ip_block.example
$ terraform import metal_reserved_ip_block.example e689072f-aa6e-4d51-8e37-c2fbe18b4ff0
$ terraform state rm packet_device.example
$ terraform import metal_device.example 8eb3bc10-0e1a-476a-aec2-6dc699df9c1c
```

We then need to install the equinix/metal provider by running `$ terraform init`. After that, our templates should be in check with the Terraform state and with the upstream resources in Equinix Metal. We can verify the migration by running `$ terraform plan`, it should show that infrastructure is up to date.

## Resolving migration issues

When we run `$ terraform plan` to verify that migration was successful, terraform might warn that some resource attributes from templates are not aligned with imported state. It's because not all of the resource attribute can be computed, for example the `ip_address` blocks in packet_device are user-defined and will result to a non-empty diff against downloaded imported state.

In case of the `ip_address`, a consequent `$ terraform apply` will update the local state without changing the upstream resource, but if an attribute causes an upstream update, you will need to resolve it manually, either changing your template, or letting Terraform change the resource upstream.


