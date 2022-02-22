# Equinix Terraform Provider Migration Tool

[Equinix Metal](https://metal.equinix.com/) (formerly Packet), has been fully integrated into Platform Equinix and therefore the teraform provider changes too. Together with Equinix Fabric and Equinix Network Edge, from `v1.5.0` the Equinix Terraform Provider will be used to interact with the resources provided by Equinix Metal.

This tool will target a terraform working directory and transform all\* `metal` or `packet` names found in *.tf* and *.tfstate* files to the `equinix` provider name. It creates a backup of the target directory *\<working-directory\>.backup* as a sibling folder.

\**This tool will not transform variable names or comments even if they contain the words `metal` or `packet`.*

## Provider Setup and Config Verfification

The migration will transform the `metal` or `packet` provider block as well as the required_providers in the terraform block and, if included, it comments the attribute `version` to take the latest available of the `equinix` provider:

From:

```hcl
terraform {
  required_providers {
    packet = {
      source   = "packethost/packet"
      version = "3.2.1"
    }
  }
}

provider "packet" {
  auth_token = var.auth_token
}
```

To:

```hcl
terraform {
  required_providers {
    equinix = {
      source   = "equinix/equinix"
      #version = "3.2.1"
    }
  }
}

provider "equinix" {
  auth_token = var.auth_token
}
```

__NOTE__

If your code already includes both `equinix` provider and `metal` | `packet`, the resulting code will have two `equinix` provider blocks and they will also be duplicated in the required_providers definition. If this is your case, after migrate you must manually combine them in a single one with all the parameters required:

From:

```hcl
provider "equinix" {
  auth_token = var.auth_token
}
provider "equinix" {
  client_id = var.client_id
  client_secret = var.client_secret
}
```

To:

```hcl
provider "equinix" {
  auth_token = var.auth_token
  client_id = var.client_id
  client_secret = var.client_secret
}
```

If you have any other requirements in the provider definition that this tool does not address, you will need to manually modify them after running a migration.

## Remote State

The **equinix-terraform-tool** does not support [remote state](https://www.terraform.io/docs/state/remote.html). If you are using remote state, then the recommended approach is to copy the state file locally, run the **equinix-terraform-tool**, and then push the state file back to the remote location. See the documentation [here](https://www.terraform.io/docs/backends/config.html) for details about how to unconfigure and reconfigure your backend.

## Using the tool

To migrate your terraform project, follow these steps:  

From the project directory, run `terraform plan`, make sure there are no pending changes in your plan.

Execute the **equinix-terraform-tool** binary, passing the path to your project directory, example:  

`equinix-terraform-tool migrate -dir=<project-path>`

After migrating, run `terraform plan` again and verify there are no new pending modifications.

For Terraform v.10+, you will need to initialize terraform for the directory using `terraform init`

If the migration was not successful, manually restore the project files
from the .backup directory or run:  

`equinix-terraform-tool backup -dir=<project-path> -restore`

After you have verified the migration was successful, delete the
backup directory or run:

`equinix-terraform-tool backup -dir=<project-path> -purge`

## Credits

Based on [OCI Provider migration tool](https://registry.terraform.io/providers/hashicorp/oci/latest/docs/guides/version-2-upgrade#migration-tool) - *Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.*
