---
page_title: "Equinix Metal: connection"
subcategory: ""
description: |-
  Retrieve Equinix Fabric Connection
---

# metal\_connection

Use this data source to retrieve a connection resource from [Equinix Fabric - software-defined interconnections](https://metal.equinix.com/developers/docs/networking/fabric/)

## Example Usage

```hcl
data "metal_connection" "example" {
  connection_id     = "4347e805-eb46-4699-9eb9-5c116e6a017d"
}
```

## Argument Reference

* `connection_id` - (Required) ID of the connection resource

## Attributes Reference

* `description` - Description of the connection resource
* `name` - Name of the connection resource
* `facility` - Slug of a facility to which the connection belongs
* `metro` - Slug of a metro to which the connection belongs
* `organization_id` - ID of organization to which the connection belongs
* `status` - Status of the connection
* `token` - Fabric Token for the [Equinix Fabric Portal](https://ecxfabric.equinix.com/dashboard)
* `type` - Connection type, dedicated or shared
* `mode` - Mode for connections in IBX facilities with the dedicated type - standard or tunnel
* `redundancy` - Connection redundancy, reduntant or primary
* `speed` - Connection speed in bits per second
* `ports` - List of connection ports - primary (`ports[0]`) and secondary (`ports[1]`)
  * `name` - Port name
  * `id` - Port UUID
  * `role` - Port role - primary or secondary
  * `speed` - Port speed in bits per second
  * `status` - Port status 
  * `link_status` - Port link status
  * `virtual_circuit_ids` - List of IDs of virtual cicruits attached to this port
