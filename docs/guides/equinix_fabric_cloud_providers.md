---
page_title: "Connecting to the cloud with Equinix Fabric via Terraform"
---

# Connecting to the cloud with Equinix Fabric via Terraform

-> **NOTE:** See the [Equinix Fabric](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/Fabric-landing-main.htm)
documentation for more details.

Equinix Fabric™ is a software-defined interconnection service that allows any business to connect
its own distributed infrastructure to any other company's infrastructure or service provider on
Platform Equinix® across a globally connected network. This guide focuses mainly on establishing an
interconnection to a cloud service provider and how to take advantage of the
[Equinix Fabric connection terraform modules](https://registry.terraform.io/search/modules?namespace=equinix-labs&q=fabric-connection),
for further details and other options you can check on [References](#references) section below.

## Getting started - Enabling an interconnection

Whether you are setting up a Multi-Cloud or Hybrid Cloud architecture, the main resource you will
need to define is an [equinix_ecx_l2_connection](../resources/equinix_ecx_l2_connection.md) which
will let Equinix to enable an interconnection on your behalf to the specified cloud provider.
However, there are other required resources that must be configured on both Equinix and your cloud
provider to have the interconnection up and running. Below we describe these general steps using
Azure ExpressRoute as an example.

**1.** Enabling interconnection in the cloud provider - Usally this implies creating a cloud router
and an interconnection asset (e.g. Google Cloud Router and VLAN attachment, Oracle FastConnect,
etc.). For this example it is required an [azurerm_express_route_circuit](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/express_route_circuit)
resource in which the service provider (i.e. `Equinix`) must be specified in order to generate a
valid pairing key.

```hcl-terraform
provider "azurerm" {
  features {}
}

resource "azurerm_express_route_circuit" "example" {
  name                  = "my-circuit"
  resource_group_name   = "my-resource-group"
  location              = "Germany West Central"
  service_provider_name = "Equinix"
  peering_location      = "Frankfurt"
  bandwidth_in_mbps     = 100
  sku {
    tier   = "Premium"
    family = "UnlimitedData"
  }
  allow_classic_operations = false
}
```

**2.** Request an Equinix Fabric connection - From an
[Equinix Fabric Port](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/ports/Fabric-port-details.htm) /
[Network Edge virtual router](https://docs.equinix.com/en-us/Content/Interconnection/NE/landing-pages/NE-landing-main.htm) /
[Equinix Service Token](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm)
to the cloud virtual interconnection asset.

In this example, we will establish a redundant connection from a Network Edge device (See [Edge networking examples](https://github.com/equinix/terraform-provider-equinix/tree/master/examples/edge-networking) for usage details of
[equinix_network_device](https://registry.terraform.io/providers/equinix/equinix/latest/docs/resources/equinix_network_device)
resource). Note that the value of `authorization_key` must be the `azurerm_express_route_circuit.example.service_key`
(i.e. the pairing key mentioned above).

```hcl-terraform
provider "equinix" {}

data "equinix_ecx_l2_sellerprofile" "azure" {
  name                     = "Azure ExpressRoute"
  organization_global_name = "Microsoft"
}

resource "equinix_ecx_l2_connection" "example" {
  name                = "my-connection-to-azure-pri"
  profile_uuid        = data.equinix_ecx_l2_sellerprofile.azure.uuid
  speed               = azurerm_express_route_circuit.example.bandwidth_in_mbps
  speed_unit          = "MB"
  notifications       = ["example@equinix.com"]
  device_uuid         = equinix_network_device.example.id
  device_interface_id = 5
  seller_metro_code   = "FR" // Frankfurt
  authorization_key   = azurerm_express_route_circuit.example.service_key
  named_tag           = "Private" // One of Private or Microsoft

  secondary_connection {
    name                = "my-connection-to-azure-sec"
    device_uuid         = equinix_network_device.example.id
    device_interface_id = 6
  }
}
```

**3.** Configure BGP in cloud side - Known as circuit peering or virtual interface, all cloud
providers offer a resource to add a BGP peer. Some commonly required details that need to be
provided are:

- Customer router IP (your destination router peer IP to which the cloud provider should send
traffic).
- Cloud router IP (virtual router peer IP to send traffic to the cloud).
- Customer BGP ASN (the Border Gateway Protocol Autonomous System Number of your on-premises peer
router).

For this example, we need to create an `azurerm_express_route_circuit_peering` where the details
for both the primary and secondary connections must be defined.

```hcl-terraform
resource "azurerm_express_route_circuit_peering" "example" {
  express_route_circuit_name = azurerm_express_route_circuit.example.name
  resource_group_name        = "my-resource-group"

  peering_type                  = "AzurePrivatePeering"
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  vlan_id                       = 300
}
```

**4.** Configure BGP in Equinix side - Finally the customer side must be configured with the same
information. Since in this example the origin of the connection is a Network Edge device you can
take advantage of the `equinix_network_bgp` resource to configure the BGP peering in the virtual
router.

```hcl-terraform
resource "equinix_network_bgp" "primary" {
  connection_id      = equinix_ecx_l2_connection.example.uuid
  local_ip_address   = "${cidrhost(azurerm_express_route_circuit_peering.example.primary_peer_address_prefix, 1)}/30"
  local_asn          = azurerm_express_route_circuit_peering.example.peer_asn
  remote_ip_address  = cidrhost(azurerm_express_route_circuit_peering.example.primary_peer_address_prefix, 2)
  remote_asn         = azurerm_express_route_circuit_peering.example.azure_asn
}

resource "equinix_network_bgp" "secondary" {
  connection_id      = equinix_ecx_l2_connection.example.uuid
  local_ip_address   = "${cidrhost(azurerm_express_route_circuit_peering.example.secondary_peer_address_prefix, 1)}/30"
  local_asn          = azurerm_express_route_circuit_peering.example.peer_asn
  remote_ip_address  = cidrhost(azurerm_express_route_circuit_peering.example.secondary_peer_address_prefix, 2)
  remote_asn         = azurerm_express_route_circuit_peering.example.azure_asn
}
```

## Terraform modules - The easiest way

Although the configuration will look similar on the Equinix side, all the resources required to
complete the configuration will depend on each cloud provider. This requires you to have prior
knowledge on how to configure the interconnection on each platform. Alternatively, you can take
advantage of the [Equinix Fabric connection terraform modules](https://registry.terraform.io/search/modules?namespace=equinix-labs&q=fabric-connection).

The terraform modules containerize multiple resources that are used together in a configuration.
With the [Equinix Fabric modules](https://registry.terraform.io/search/modules?namespace=equinix-labs&q=fabric-connection)
you can configure all described above just defining a single resource. In addition, in the
terraform registry you will get information about the specific parameters required for your cloud
provider's configuration, as well as some examples with the most frequent use cases.

Below code is all you need to fully replace the example above:

```hcl-platform
# main.tf
provider "equinix" {}

provider "azurerm" {
  features {}
}

module "equinix-fabric-connection-azure" {
  source = "github.com/equinix-labs/terraform-equinix-fabric-connection-azure"
  
  # required variables
  fabric_notification_users = ["example@equinix.com"]

  # optional variables
  network_edge_device_id           = var.primary_device_id
  network_edge_secondary_device_id = var.secondary_device_id
  network_edge_configure_bgp       = true

  fabric_speed = 100

  az_region = "Germany West Central"

  az_exproute_peering_customer_asn      = 100
  az_exproute_peering_primary_address   = "169.0.0.0/30"
  az_exproute_peering_secondary_address = "169.0.0.4/30"
}
```

There are (April 2022) modules available to interconnect with
[AWS Direct Connect](https://registry.terraform.io/modules/equinix-labs/fabric-connection-aws/equinix/latest),
[Azure ExpressRoute](https://registry.terraform.io/modules/equinix-labs/fabric-connection-azure/equinix/latest) and
[Google Cloud Interconnect](https://registry.terraform.io/modules/equinix-labs/fabric-connection-gcp/equinix/latest).
Some others will be available soon. If you don't find one for your cloud provider you can still
take advantage of the [Equinix Fabric connection base module](https://registry.terraform.io/modules/equinix-labs/fabric-connection/equinix/latest)
or open a ticket in the [github repository](https://github.com/equinix/terraform-provider-equinix)
with your request.

## References

- See [Equinix connectivity examples](https://github.com/equinix/terraform-provider-equinix/tree/master/examples/connectivity)
for providers for which there is no module available yet.
- See the [API how to guides](https://developer.equinix.com/docs/how-guide-v3-apis) for further
details on each cloud service provider requirements.
- Check the [available providers](https://www.equinix.com/interconnection-services/equinix-fabric/provider-availability)
on Platform Equinix® to find your required service provider.
