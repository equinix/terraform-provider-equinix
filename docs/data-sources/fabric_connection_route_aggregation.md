---
subcategory: "Fabric"
---

# equinix_fabric_connection_route_aggregation (Data Source)

Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Connection Route Aggregation by UUID
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations

## Example Usage

```terraform
data "equinix_fabric_connection_route_aggregation" "attached_policy" {
  route_aggregation_id = "<route_aggregation_id>"
  connection_id = "<connection_id>"
}

output "connection_route_Aggregation_id" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.id
}

output "connection_route_aggregation_connection_id" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.connection_id
}

output "connection_route_aggregation_type" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.type
}

output "connection_route_aggregation_attachment_status" {
  value = data.equinix_fabric_connection_route_aggregation.attached_policy.attachment_status
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `connection_id` (String) The uuid of the connection this data source should retrieve
- `route_aggregation_id` (String) The uuid of the route aggregation this data source should retrieve

### Read-Only

- `attachment_status` (String) Status of the Route Aggregation Policy attachment lifecycle
- `href` (String) URI to the attached Route Aggregation Policy on the Connection
- `id` (String) The unique identifier of the resource
- `type` (String) Route Aggregation Type. One of ["BGP_IPv4_PREFIX_AGGREGATION"]
- `uuid` (String) Equinix Assigned ID for Route Aggregation Policy
