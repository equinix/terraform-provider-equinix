# Equinix Provider Examples

This directory contains a set of examples of using Equinix services with Terraform.
Each example has its own README file containing more details on what it does.

To run any example, clone the repository, **adjust variables**, initialize plugins
and run `terraform apply` within the example's own directory.

```sh
git clone https://github.com/equinix/terraform-provider-equinix
cd terraform-provider-equinix/examples/aws-connection
terraform init
terraform apply
```
