---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "equinix_fabric_networks Data Source - terraform-provider-equinix"
subcategory: "Fabric"
description: |-
  Fabric V4 API compatible data resource that allow user to fetch Fabric Network for a given UUID
---

# equinix_fabric_networks (Data Source)

Fabric V4 API compatible data resource that allow user to fetch Fabric Network for a given UUID

Additional documentation:
* Getting Started: <https://docs.equinix.com/en-us/Content/Interconnection/Fabric/IMPLEMENTATION/fabric-networks-implement.htm>
* API: <https://developer.equinix.com/dev-docs/fabric/api-reference/fabric-v4-apis#fabric-networks>

## Example Usage

```hcl
data "equinix_fabric_networks" "test" {
    outer_operator = "AND"
    filter {
        property = "/type"
        operator = "="
        values 	 = ["IPWAN"]
    }
    filter {
        property = "/name"
        operator = "="
        values   = ["Tf_Network_PFCR"]
    }
    filter {
        group = "OR_group1"
        property = "/operation/equinixStatus"
        operator = "="
        values = ["PROVISIONED"]
    }
    filter {
        group = "OR_group1"
        property = "/operation/equinixStatus"
        operator = "LIKE"
        values = ["DEPROVISIONED"]
    }
    pagination {
        offset = 0
        limit = 5
    }
    sort {
        direction = "ASC"
        property = "/name"
    }
}

output "number_of_returned_networks" {
    value = length(data.equinix_fabric_networks.test.data)
}

output "first_network_name" {
    value = data.equinix_fabric_networks.test.data.0.name
}

output "first_network_connections_count" {
    value = data.equinix_fabric_networks.test.data.0.connections_count
}

output "first_network_scope" {
    value = data.equinix_fabric_networks.test.data.0.scope
}

output "first_network_type" {
    value = data.equinix_fabric_networks.test.data.0.type
}

output "first_network_location_region" {
    value = one(data.equinix_fabric_networks.test.data.0.location).region
}

output "first_network_project_id" {
    value = one(data.equinix_fabric_networks.test.data.0.project).project_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `filter` (Block List, Min: 1, Max: 10) Filters for the Data Source Search Request (see [below for nested schema](#nestedblock--filter))
- `outer_operator` (String) Determines if the filter list will be grouped by AND or by OR. One of [AND, OR]

### Optional

- `pagination` (Block Set, Max: 1) Pagination details for the Data Source Search Request (see [below for nested schema](#nestedblock--pagination))
- `sort` (Block List) Filters for the Data Source Search Request (see [below for nested schema](#nestedblock--sort))

### Read-Only

- `data` (List of Object) List of Cloud Routers (see [below for nested schema](#nestedatt--data))
- `id` (String) The ID of this resource.

<a id="nestedblock--filter"></a>
### Nested Schema for `filter`

Required:

- `operator` (String) Operators to use on your filtered field with the values given. One of [ =, !=, >, >=, <, <=, BETWEEN, NOT BETWEEN, LIKE, NOT LIKE, ILIKE, NOT ILIKE, IN, NOT IN]
- `property` (String) Possible field names to use on filters. One of [/name /uuid /scope /type /operation/equinixStatus /location/region /project/projectId /account/globalCustId /account/orgId /deletedDate /_*]
- `values` (List of String) The values that you want to apply the property+operator combination to in order to filter your data search

Optional:

- `group` (String) Optional custom id parameter to assign this filter to an inner AND or OR group. Group id must be prefixed with AND_ or OR_. Ensure intended grouped elements have the same given id. Ungrouped filters will be placed in the filter list group by themselves.


<a id="nestedblock--pagination"></a>
### Nested Schema for `pagination`

Optional:

- `limit` (Number) Number of elements to be requested per page. Number must be between 1 and 100. Default is 20
- `offset` (Number) The page offset for the pagination request. Index of the first element. Default is 0.


<a id="nestedblock--sort"></a>
### Nested Schema for `sort`

Optional:

- `direction` (String) The sorting direction. Can be one of: [DESC, ASC], Defaults to DESC
- `property` (String) The property name to use in sorting. One of [/name /uuid /scope /operation/equinixStatus /location/region /changeLog/createdDateTime /changeLog/updatedDateTime]. Defaults to /changeLog/updatedDateTime


<a id="nestedatt--data"></a>
### Nested Schema for `data`

Read-Only:

- `change` (Set of Object) (see [below for nested schema](#nestedobjatt--data--change))
- `change_log` (Set of Object) (see [below for nested schema](#nestedobjatt--data--change_log))
- `connections_count` (Number)
- `href` (String)
- `location` (Set of Object) (see [below for nested schema](#nestedobjatt--data--location))
- `name` (String)
- `notifications` (List of Object) (see [below for nested schema](#nestedobjatt--data--notifications))
- `operation` (Set of Object) (see [below for nested schema](#nestedobjatt--data--operation))
- `project` (Set of Object) (see [below for nested schema](#nestedobjatt--data--project))
- `scope` (String)
- `state` (String)
- `type` (String)
- `uuid` (String)

<a id="nestedobjatt--data--change"></a>
### Nested Schema for `data.change`

Read-Only:

- `href` (String)
- `type` (String)
- `uuid` (String)


<a id="nestedobjatt--data--change_log"></a>
### Nested Schema for `data.change_log`

Read-Only:

- `created_by` (String)
- `created_by_email` (String)
- `created_by_full_name` (String)
- `created_date_time` (String)
- `deleted_by` (String)
- `deleted_by_email` (String)
- `deleted_by_full_name` (String)
- `deleted_date_time` (String)
- `updated_by` (String)
- `updated_by_email` (String)
- `updated_by_full_name` (String)
- `updated_date_time` (String)


<a id="nestedobjatt--data--location"></a>
### Nested Schema for `data.location`

Read-Only:

- `ibx` (String)
- `metro_code` (String)
- `metro_name` (String)
- `region` (String)


<a id="nestedobjatt--data--notifications"></a>
### Nested Schema for `data.notifications`

Read-Only:

- `emails` (List of String)
- `send_interval` (String)
- `type` (String)


<a id="nestedobjatt--data--operation"></a>
### Nested Schema for `data.operation`

Read-Only:

- `equinix_status` (String)


<a id="nestedobjatt--data--project"></a>
### Nested Schema for `data.project`

Read-Only:

- `project_id` (String)