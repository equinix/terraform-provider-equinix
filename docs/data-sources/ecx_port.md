---
layout: "equinix"
page_title: "Equinix: equinix_ecx_port"
sidebar_current: "docs-equinix-datasource-ecx-port"
description: |-
  Provides a Equinix ecx_port datasource. This can be used to read existing ecx_ports.
---

# equinix_ecx_port

Data source `equinix_ecx_port` is used to fetch attributes of ECX port (like UUID) with given port name.

The data source relies on the Equinix Cloud Exchange Fabric API. The parameters
and attributes map to the fields described at
<https://developer.equinix.com/docs/ecx-layer-2-buyer-apis-v3#get-user-port>.

## Example Usage

```hcl
data "equinix_ecx_port" "tf-pri-dot1q" {
  name = "sit-001-CX-NY5-NL-Dot1q-BO-10G-PRI-JP-157"
}

output "id" {
  value = data.equinix_ecx_port.tf-pri-dot1q.id
}
```

## Argument Reference

- `name` - _(Required)_ Name of the port

## Attributes Reference

The following attributes are exported:

- `uuid` - Unique identifier of the port
- `status` - Status of the connection
- `region` - Region in which the port resides
- `ibx` - Equinix IBX where the port resides.
- `metro_code` - The metro code of the metro where the port resides
- `priority` - The priority of the device (primary / secondary) where the port
  resides
- `encapsulation` - The VLAN encapsulation of the port (Dot1q or QinQ)
- `buyout` - Indicates whether the port supports unlimited connections. If
  "false", the port is a standard port with limited connections. If "true", the
  port is an "unlimited connections" port that allows multiple connections at no
  additional charge.
- `bandwidth` - Port Bandwidth in bytes.
- `status` - Port status that indicates whether a port has been assigned or is
  ready for connection.
