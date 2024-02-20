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
need to define is an [equinix_fabric_connection](../resources/equinix_fabric_connection.md) which
will let Equinix enable an interconnection on your behalf to the specified cloud provider.
However, there are other required resources that must be configured on both Equinix and your cloud
provider to have the interconnection up and running. Below we describe these general steps using
Azure ExpressRoute as an example.

**1.** Enabling interconnection in the cloud provider - Usually this implies creating a cloud router
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
[Network Edge Device](https://docs.equinix.com/en-us/Content/Interconnection/NE/landing-pages/NE-landing-main.htm) /
[Equinix Service Token](https://docs.equinix.com/en-us/Content/Interconnection/Fabric/service%20tokens/Fabric-Service-Tokens.htm) /
[Equinix Fabric Cloud Router](https://docs.equinix.com/en-us/Content/Interconnection/FCR/FCR-intro.htm)
to the cloud virtual interconnection asset.

In this example, we will establish a connection from a Cloud Router
resource). Note that the value of `authorization_key` must be the `azurerm_express_route_circuit.example.service_key`
(i.e. the pairing key mentioned above).

```hcl-terraform
provider "equinix" {}

data "equinix_fabric_service_profile" "azure" {
  name                     = "Azure ExpressRoute"
  filter {
    property = "/name"
    operator = "="
    values   = ["Azure ExpressRoute"]
  }
}

resource "equinix_fabric_connection" "fcr2azure"{
  name = "ConnectionName"
  type = "IP_VC"
  notifications{
    type = "ALL"
    emails = ["example@equinix.com","test1@equinix.com"]
  }
  bandwidth = azurerm_express_route_circuit.example.bandwidth_in_mbps
  order {
  purchase_order_number = "1-323292"
  }
  a_side {
    access_point {
      type = "CLOUD_ROUTER"
      router {
        uuid = "<cloud_router_uuid>"
      }
    }
  }
  z_side {
    access_point {
      type = "SP"
      authentication_key = azurerm_express_route_circuit.example.service_key
      peering_type = "PRIVATE"
      profile {
      type = "L2_PROFILE"
        uuid = data.equinix_fabric_service_profile.azure.data.0.id
      }
      location {
        metro_code = "SV"
      }
    }
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
  bandwidth_in_mpbs             = 50
}
```

**4.** Configure BGP in Equinix side - Finally the customer side must be configured with the same
information. Since in this example the origin of the connection is a Fabric Cloud Router you can
take advantage of the `equinix_fabric_routing_protocol` resource to configure the BGP peering in the cloud
router.

```hcl-terraform
resource "equinix_fabric_routing_protocol" "direct"{
  connection_uuid = equinix_fabric_connection.fcr2azure.id
  type = "DIRECT"
  name = "direct_rp"
  direct_ipv4 {
    equinix_iface_ip = "190.1.1.1/30"
  }
  direct_ipv6{
    equinix_iface_ip = "190::1:1/126"
  }
}

resource "equinix_fabric_routing_protocol" "bgp" {
  depends_on = [
    equinix_fabric_routing_protocol.direct
  ]
  connection_uuid = equinix_fabric_connection.fcr2azure.id
  type            = "BGP"
  name            = "bgp_rp"
  bgp_ipv4 {
    customer_peer_ip = "190.1.1.2"
    enabled          = true
  }
  bgp_ipv6 {
    customer_peer_ip = "190::1:2"
    enabled          = true
  }
  customer_asn = 4532
}
```

## Terraform Modules - The easiest way

Although the configuration will look similar on the Equinix side, all the resources required to
complete the configuration will depend on each cloud provider. This requires you to have prior
knowledge on how to configure the interconnection on each platform. Alternatively, you can take
advantage of the [Equinix Fabric Terraform Modules](https://registry.terraform.io/modules/equinix/fabric/equinix/latest).

The terraform modules containerize multiple resources that are used together in a configuration.
With the [Equinix Fabric Terraform Cloud Router 2 Azure Connection Example Module](https://registry.terraform.io/modules/equinix/fabric/equinix/latest/examples/cloud-router-2-azure-connection)
and the [Equinix Fabric Terraform Routing Protocols Module](https://registry.terraform.io/modules/equinix/fabric/equinix/latest/examples/routing-protocols)
you can configure all described above by just defining three resources.

Below code is all you need to fully replace the example above:

```hcl-terraform
# main.tf
provider "equinix" {}

provider "azurerm" {}

module "equinix-fabric-connection-azure" {
  source = "equinix/fabric/equinix//examples/cloud-router-2-azure-connection"


  # Connection Details
  connection_name                 = "fcr_2_azure"
  connection_type                 = "IP_VC"
  notifications_type              = "ALL"
  notifications_emails            = ["example@equinix.com","test1@equinix.com"]
  purchase_order_number           = "1-323292"
  bandwidth                       = 50
  aside_ap_type                   = "CLOUD_ROUTER"
  aside_fcr_uuid                  = "<Primary Fabric Cloud router UUID>"
  zside_ap_type                   = "SP"
  zside_ap_profile_type           = "L2_PROFILE"
  zside_location                  = "SV"
  zside_peering_type              = "PRIVATE"
  zside_fabric_sp_name            = "Azure ExpressRoute"
  
  # Azure details
  azure_client_id                 = "<Azure Client Id>"
  azure_client_secret             = "<Azure Client Secret Value>"
  azure_tenant_id                 = "<Azure Tenant Id>"
  azure_subscription_id           = "<Azure Subscription Id>"
  azure_resource_name             = "my-resource-group"
  azure_location                  = "West US 2"
  azure_service_key_name          = "Test_Azure_Key"
  azure_service_provider_name     = "<Service Provider Name>"
  azure_peering_location          = "Silicon Valley Test"
  azure_tier                      = "Standard"
  azure_family                    = "UnlimitedData"
  azure_environment               = "PROD"
}

resource "azurerm_express_route_circuit_peering" "example" {
  express_route_circuit_name = module.equinix-fabric-connection-azure.azurerm_express_route_circuit_name
  resource_group_name        = "my-resource-group"

  peering_type                  = "AzurePrivatePeering"
  peer_asn                      = 100
  primary_peer_address_prefix   = "123.0.0.0/30"
  secondary_peer_address_prefix = "123.0.0.4/30"
  vlan_id                       = 300
  bandwidth_in_mpbs             = 50
}

module "routing_protocols" {
  source = "equinix/fabric/equinix//modules/routing-protocols"

  connection_uuid = module.equinix-fabric-connection-azure.module_output

  # Direct RP Details
  direct_rp_name         = "direct_rp"
  direct_equinix_ipv4_ip = "190.1.1.1/30"
  direct_equinix_ipv6_ip = "190::1:1/126"

  # BGP RP Details
  bgp_rp_name            = "bgp_rp"
  bgp_customer_asn       = 4532
  bgp_customer_peer_ipv4 = "190.1.1.2"
  bgp_enabled_ipv4       = true
  bgp_customer_peer_ipv6 = var.bgp_customer_peer_ipv6
  bgp_enabled_ipv6       = "190::1:2"
}
```

- See [Equinix Connectivity Examples on Fabric Modules Page](https://registry.terraform.io/modules/equinix/fabric/equinix/latest)
  for information on more Fabric use cases that are supported by Fabric Terraform Modules

If you don't find one for your cloud provider you can open a ticket in the [github repository](https://github.com/equinix/terraform-provider-equinix/issues)
with your request.

## References

- See the [API how to guides](https://developer.equinix.com/docs?page=/dev-docs/fabric/overview) for further
details on each cloud service provider requirements.
- Check the [available providers](https://www.equinix.com/interconnection-services/equinix-fabric/provider-availability)
on Platform Equinix® to find your required service provider.
