---
subcategory: "Fabric"
---

# equinix_ecx_port (Data Source)

Use this data source to get details of Equinix Fabric port with a given name.

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

* `name` - (Required) Name of the port

## Attributes Reference

The following attributes are exported:

* `uuid` - Unique identifier of the port
* `status` - Status of the port
* `region` - Port location region
* `ibx` - Port location Equinix Business Exchange (IBX)
* `metro_code` - Port location metro code
* `priority` - The priority of the device (primary / secondary) where the port
  resides
* `encapsulation` - The VLAN encapsulation of the port (Dot1q or QinQ)
* `buyout` - Boolean value that indicates whether the port supports unlimited connections.
  If "false", the port is a standard port with limited connections. If "true", the
  port is an "unlimited connections" port that allows multiple connections at no
  additional charge.
* `bandwidth` - Port Bandwidth in bytes.
* `status` - Port status that indicates whether a port has been assigned or is
  ready for connection.
