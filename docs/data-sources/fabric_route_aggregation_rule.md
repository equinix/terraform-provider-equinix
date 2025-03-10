---
subcategory: "Fabric"
---

# equinix_fabric_route_aggregation_rule (Data Source)

Fabric V4 API compatible data resource that allow user to fetch Equinix Fabric Route Aggregation Rule by UUID
Additional Documentation:
* API: https://developer.equinix.com/catalog/fabricv4#tag/Route-Aggregations

## Example Usage

```terraform
data "equinix_fabric_route_aggregation_rule" "ra_rule" {
  route_aggregation_id = "<route_aggregation_id>"
  route_aggregation_rule_id = "<route_aggregation_rule_id>"
}

output "route_aggregation_rule_name" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.name
}

output "route_aggregation_rule_description" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.description
}

output "route_aggregation_rule_type" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.type
}

output "route_aggregation_rule_prefix" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.prefix
}

output "route_aggregation_rule_state" {
  value = data.equinix_fabric_route_aggregation_rule.ra_rule.state
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `route_aggregation_id` (String) The uuid of the route aggregation this data source should retrieve
- `route_aggregation_rule_id` (String) The uuid of the route aggregation rule this data source should retrieve

### Optional

- `description` (String) Customer-provided route aggregation rule description

### Read-Only

- `change` (Attributes) Current state of latest route aggregation rule change (see [below for nested schema](#nestedatt--change))
- `change_log` (Attributes) Details of the last change on the stream resource (see [below for nested schema](#nestedatt--change_log))
- `href` (String) Equinix auto generated URI to the route aggregation rule resource
- `id` (String) The unique identifier of the resource
- `name` (String) Customer provided name of the route aggregation rule
- `prefix` (String) Customer-provided route aggregation rule prefix
- `state` (String) Value representing provisioning status for the route aggregation rule resource
- `type` (String) Equinix defined Route Aggregation Type; BGP_IPv4_PREFIX_AGGREGATION, BGP_IPv6_PREFIX_AGGREGATION
- `uuid` (String) Equinix-assigned unique id for the route aggregation rule resource

<a id="nestedatt--change"></a>
### Nested Schema for `change`

Required:

- `type` (String) Equinix defined Route Aggregation Change Type
- `uuid` (String) Equinix-assigned unique id for a change

Read-Only:

- `href` (String) Equinix auto generated URI to the route aggregation change


<a id="nestedatt--change_log"></a>
### Nested Schema for `change_log`

Read-Only:

- `created_by` (String) User name of creator of the stream resource
- `created_by_email` (String) Email of creator of the stream resource
- `created_by_full_name` (String) Legal name of creator of the stream resource
- `created_date_time` (String) Creation time of the stream resource
- `deleted_by` (String) User name of deleter of the stream resource
- `deleted_by_email` (String) Email of deleter of the stream resource
- `deleted_by_full_name` (String) Legal name of deleter of the stream resource
- `deleted_date_time` (String) Deletion time of the stream resource
- `updated_by` (String) User name of last updater of the stream resource
- `updated_by_email` (String) Email of last updater of the stream resource
- `updated_by_full_name` (String) Legal name of last updater of the stream resource
- `updated_date_time` (String) Last update time of the stream resource
