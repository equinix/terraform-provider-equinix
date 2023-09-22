# Fabric Terraform Example Scripts

## Setup

To use one of the Fabric examples, be sure to copy the `.tfvars.example`
to a `.tfvars` file in the folder of the script you want to run. Then you will
need to fill in the variable values with details specific to your account.

I.e.
`cp terraform.tfvars.example terraform.tfvars`

Additionally, if you do not want to use a `.tfvars` file, you can set each variable
with shell environment variables. You just need to prefix each variable name with TF_VAR.
Example: `equinix_client_id` can be set in your shell with `export TF_VAR_equinix_client_id=<id_value>`

You can also use a mix of variables defined in `.tfvars` and variables defined in your
shell environment. The `.tfvars` file will take priority over any variables defined in
your shell though. Variable definition priority is defined
[in this Hashicorp Guide](https://developer.hashicorp.com/terraform/language/values/variables#variable-definition-precedence).

## Recommended Usage

Place your secrets in your shell with `TF_VAR_equinix_client_id` and `TF_VAR_equinix_client_secret`
and only place the specific example scenario details in `.tfvars`. This provides flexibility to move
between examples and run them without needed to copy your secrets to many different places while
learning how to leverage Fabric terraform for your use cases. 

**Note: you'll have to delete any reference to
`equinix_client_id` and `equinix_client_secret` in the `.tfvars` file to let Terraform read your shell variables
during apply.**
