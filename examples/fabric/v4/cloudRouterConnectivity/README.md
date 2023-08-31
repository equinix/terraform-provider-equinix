# Equinix Provider Examples

This directory contains a set of examples of using Equinix services with Terraform.
Each example has its own README file containing more details on what it does.

Equinix Provider examples are grouped into following directories:

* [connectivity](fabric/) - examples of **establishing connectivity with
service providers** that are part of Equinix Fabric community, including major
Cloud Service Providers like Google, Amazon or Microsoft
* [edge-networking](edge-networking/) - examples of running and connecting
**virtual network devices at the network and compute edge**

## Using examples

To run any example, clone the repository, **adjust variables**, initialize plugins
and run `terraform apply` within the example's own directory.

```sh
git clone https://github.com/equinix/terraform-provider-equinix
cd terraform-provider-equinix/examples/fabric/aws-connection
terraform init
terraform apply
```
