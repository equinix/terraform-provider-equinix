---
page_title: "Migrating from the Equinix Metal provider"
---

# Migrating from equinix/metal to equinix/equinix

[Equinix Metal](https://metal.equinix.com/) (formerly Packet), has been fully integrated into Platform Equinix and therefore the terraform provider changes too. This (terraform-provider-equinix, provider equinix/equinix) is the current provider of the various services available on Platform Equinix that can be managed using Terraform.

If you've been using terraform-provider-metal, and you want to use a newer provider version to manage resources in Equinix Metal, you will need to change the references in you HCL files. You can just change the names of the resources, e.g. from `metal_device` to `equinix_metal_device`. That should work, but it will cause the `metal_device` to be destroyed and new `equinix_metal_device` to be created instead. Re-creation of the resources might be undesirable, and this guide shows how to migrate to `equinix_metal_` resources without the re-creation.

Before starting to migrate your Terraform templates, please upgrade

* equinix/metal provider to the latest version (3.2.1)
* Terraform to version at least v0.13

## Fast migration with replace-provider and sed

Just like the Terraform HCL templates, the Terraform state is a file containing resource names and their attributes in structured text. We can attempt the migration as a text substitution task, basically replacing `metal_` with `equinix_metal_` wherever possible, and fixing the provider source reference.

It's a good idea to make a backup of the whole Terraform directory before doing this.

Considering we have infrastructure created from following template:

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
  project_id  = metal_project.example.id
  facility    = "sv15"
  description = "example"
}
```

We can first change the provider in Terraform state file (`terraform.tfstate`) with `terraform state` subcommand `replace-provider`:

```shell
terraform state replace-provider equinix/metal equinix/equinix
```

Then we replace the provider reference in the HCL templates. Do this for every file where you have the reference:

```shell
sed -i 's|equinix/metal|equinix/equinix|g' main.tf
```

Then we simply replace all strings `metal_` with `equinix_metal_` in the Terraform HCL files.

```shell
sed -i 's/metal_/equinix_metal_/g' main.tf
```

..this is a bit dangerous, so check your `git diff` after. It should replace all the `metal_` prefixes and also the key from the `required_providers` block.

Then replace `metal_` with `equinix_metal_` in the terraform state file:

```shell
sed -i 's/metal_/equinix_metal_/g' terraform.tfstate
```

The example template would now look as:

```hcl-terraform
terraform {
  required_providers {
    equinix = {
      source = "equinix/equinix"
    }
  }
}

resource "equinix_metal_project" "example" {
  name = "example"
}

resource "equinix_metal_vlan" "example" {
  project_id  = equinix_metal_project.example.id
  facility    = "sv15"
  description = "example"
}
```

We then need to install the `equinix/equinix` provider by running `terraform init`. After that, our templates should be in check with the Terraform state and with the upstream resources in Equinix Metal. You can verify the result by running `terraform plan`.

If the plan is not empty, it means that some resources can not be simply read fom upstream, or that attributes have changed between your version of the `equinix/metal` provider and the current version of the `equinix/equinix` provider.

## Migrating one resource at a time

We can use `terraform state` and `terraform import` to achieve transition without destroying existing resources.

### Existing infrastructure

We assume to have infrastructure created with provider `equinix/metal` with a device and an IP reservation. The HCL looks like:

```hcl
terraform {
  required_providers {
    metal = {
      source = "equinix/metal"
      version = "3.2.1"
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

### Resource UUIDs

In order to transition to provider `equinix/equinix`, we need to find out UUIDs of all the resources we want to migrate. In this case `metal_reserved_ip_block.example` and `metal_device.example`. We can use `terraform state` to find out the UUIDs.

For the reserved IP block:

```shell
$ terraform state show metal_reserved_ip_block.example

# metal_reserved_ip_block.example:
resource "metal_reserved_ip_block" "example" {
  [...]
    id = "e689072f-aa6e-4d51-8e37-c2fbe18b4ff0"
  [...]
}
```

For the device:

```shell
$ terraform state show metal_device.example

# metal_device.example
resource "metal_device" "example" {
  [...]
    id = "8eb3bc10-0e1a-476a-aec2-6dc699df9c1c"
  [...]

```

### Migrated template

Once we find out the UUIDs of resources to migrate, in the HCL template, we need to change:

* the required_providers block to require `equinix/equinix`
* the names of the resources to corresponding resources from provider `equinix/equinix`: `sed 's/metal_/equinix_metal_'`
* all the references from `metal_` resources to `equinix_metal_` resources

The modified template will then look as:

```hcl
terraform {
  required_providers {
    equinix = {
      source = "equinix/equinix"
    }
  }
}

resource "equinix_metal_reserved_ip_block" "example" {
  project_id = local.project_id
  facility   = "sv15"
  quantity   = 2
}

resource "equinix_metal_device" "example" {
  project_id       = local.project_id
  facilities       = ["sv15"]
  plan             = "c3.medium.x86"
  operating_system = "ubuntu_20_04"
  hostname         = "test"
  billing_cycle    = "hourly"

  ip_address {
    type            = "public_ipv4"
    cidr            = 31
    reservation_ids = [equinix_metal_reserved_ip_block.example.id]
  }

  ip_address {
    type = "private_ipv4"
  }
}
```

### Migrating Terraform state

Once we changed the template accordingly, we can remove the old `metal_` resources from Terraform state and import the new ones as `equinix_metal_` resources by their UUIDs.

From checking the state before, we remember that UUID of the metal_device.example is 8eb3bc10-0e1a-476a-aec2-6dc699df9c1c, and UUID of the metal_reserved_ip_block.example is e689072f-aa6e-4d51-8e37-c2fbe18b4ff0.

 In the terraform state and import commands, we use the resource type and name, separated by dot:

```shell
terraform state rm metal_reserved_ip_block.example
terraform import equinix_metal_reserved_ip_block.example e689072f-aa6e-4d51-8e37-c2fbe18b4ff0
terraform state rm metal_device.example
terraform import equinix_metal_device.example 8eb3bc10-0e1a-476a-aec2-6dc699df9c1c
```

We then need to install the `equinix/equinix` provider by running `terraform init`. After that, our templates should be in check with the Terraform state and with the upstream resources in Equinix Metal. We can verify the migration by running `terraform plan`, it should show that infrastructure is up to date.

## Resolving migration issues

When we run `terraform plan` to verify that migration was successful, terraform might warn that some resource attributes from templates are not aligned with imported state. It's because not all of the resource attribute can be computed, for example the `ip_address` blocks in `metal_device` are user-defined and will result to a non-empty diff against downloaded imported state.

In case of the `ip_address`, a consequent `terraform apply` will update the local state without changing the upstream resource, but if an attribute causes an upstream update, you will need to resolve it manually, either changing your template, or letting Terraform change the resource upstream.
