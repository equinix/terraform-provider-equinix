---
layout: "packet"
page_title: "Packet: packet_organization"
sidebar_current: "docs-packet-datasource-organization"
description: |-
  Provides a Packet Organization datasource. This can be used to read existing Organizations.
---

# packet_organization

Provides a Packet organization datasource.

## Example Usage

```hcl
# Fetch a organization data and show projects which belong to it
data "packet_organization" "test" {
  organization_id = local.org_id
}

output "projects_in_the_org" {
  value = data.packet_organization.test.project_ids
}
```

## Argument Reference

The following arguments are supported:

* `name` - The organization name
* `organization_id` - The UUID of the organization resource

Exactly one of `name` or `organization_id` must be given.

## Attributes Reference

The following attributes are exported:

* `project_ids` - UUIDs of project resources which belong to this organization
* `description` - Description string
* `website` - Website link
* `twitter` - Twitter handle
* `logo` - Logo URL
