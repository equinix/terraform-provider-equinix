# Equinix Provider Fabric v3 Examples - connectivity

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

## Equinix Fabric v3 Terraform Modules

Terraform modules encapsulate groups of resources dedicated to one task, reducing
the amount of code you have to develop for similar infrastructure components.

Below table lists Terraform modules that can be used for convenient and
quick deployment of Equinix Fabric connections to most popular Service Providers.
They include all the necessary resources and configuration at both ends of the
connection.

Please check module links to for usage details and examples.

| Service Provider | Terraform module |
|------------------|------------------|
| Alibaba Cloud Express Connect | [Equinix Fabric L2 Connection To Alibaba Express Connect Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-alibaba/equinix/latest) |
| AWS Direct Connect / AWS Direct Connect - High Capacity | [Equinix Fabric L2 Connection To AWS Direct Connect Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-aws/equinix/latest) |
| Azure ExpressRoute | [Equinix Fabric L2 Connection To Microsoft Azure ExpressRoute Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-azure/equinix/latest) |
| Equinix Metal | [Equinix Fabric L2 Connection To Equinix Metal Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-metal/equinix/latest) |
| Google Cloud Partner Interconnect Zone 1 / Zone 2 | [Equinix Fabric L2 Connection To Google Cloud Interconnect Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-gcp/equinix/latest) |
| IBM Cloud Direct Link 2 | [Equinix Fabric L2 Connection To IBM Direct Link 2.0 Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-ibm/equinix/latest) |
| Oracle Cloud Infrastructure -OCI- FastConnect | [Equinix Fabric L2 Connection To Oracle Cloud Infrastructure FastConnect Terraform module](https://registry.terraform.io/modules/equinix-labs/fabric-connection-oci/equinix/latest) |

## Equinix Fabric Examples Without Modules

If you don't want to take advantage of the [Equinix Fabric Terraform Modules](#equinix-fabric-terraform-modules)
you can also find below some basic examples on how to establish port connectivity with
the most popular service providers.

* [alibaba-cloud](./alibaba-cloud-connection) - establishing layer 2 connection between
Equinix Fabric port and Alibaba Express Connect
* [aws-connection](./aws-connection) - establishing layer 2 connection between
Equinix Fabric port and AWS Direct Connect
* [azure-connection](./azure-connection) - establishing layer 2 connection between
Equinix Fabric port and Microsoft Azure ExpressRoute
* [equinix-metal](./equinix-metal-to-fabric-connection) - establishing layer 2 connection between
Equinix Fabric port and Equinix Metal
* [gcp-connection](./gcp-connection) - establishing layer 2 connection between
Equinix Fabric port and Google Cloud Partner Interconnect
* [ibm-cloud-connection](./ibm-cloud-connection) - establishing layer 2 connection
between Equinix Fabric port and IBM Direct Link
* [oracle-cloud-connection](./oracle-cloud-connection) - establishing layer2 connection
between Equinix Fabric port and Oracle Cloud FastConnect

## Equinix Fabric Examples To Connect To Your Own Assets

* [self-port-port-connection](./self-port-port-connection) - establishing layer2 connection
between two Equinix Fabric ports

## Equinix Fabric examples To Create a Seller Profile
Examples of using Equinix Fabric resources
to establish connectivity with most popular service providers

* [private-profile](./private-profile) - creating layer 2 private service  profile
* [public-profile](./public-profile) - creating layer 2 public service seller profile
