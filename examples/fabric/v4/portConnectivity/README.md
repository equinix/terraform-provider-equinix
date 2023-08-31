# Equinix Provider Examples - connectivity

## Equinix Fabric Terraform Modules

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
you can also find below some basic examples on how to establish connectivity with
the most popular service providers.

* [alibaba-cloud](./alibaba) - establishing layer 2 connection between
  Equinix Fabric port and Alibaba Express Connect
* [aws-connection](./aws) - establishing layer 2 connection between
  Equinix Fabric port and AWS Direct Connect
* [azure-connection](./azure) - establishing layer 2 connection between
  Equinix Fabric port and Microsoft Azure ExpressRoute
* [gcp-connection](./google) - establishing layer 2 connection between
  Equinix Fabric port and Google Cloud Partner Interconnect
* [oracle-cloud-connection](./oracle) - establishing layer2 connection
  between Equinix Fabric port and Oracle Cloud FastConnect

## Equinix Fabric Examples To Connect To Your Own Assets

* [self-port-port-connection](./port2portself) - establishing layer2 connection
  between two Equinix Fabric ports

## Equinix Fabric examples To Connect to a Seller Profile
Examples of using Equinix Fabric resources
to establish connectivity with most popular service providers

* [private-profile](./port2serviceprofileprivate) - creating layer 2 private service profile connection
* [public-profile](./port2serviceprofilepublic) - creating layer 2 public service seller profile connection