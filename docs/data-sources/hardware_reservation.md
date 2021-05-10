---
page_title: "Equinix Metal: metal_hardware_reservation"
subcategory: ""
description: |-
  Retrieve Equinix Metal Hardware Reservation
---

# metal_hardware_reservation

Use this data source to retrieve a [hardware reservation resource from Equinix Metal](https://metal.equinix.com/developers/docs/deploy/reserved/).

## Example Usage

```hcl
data "hardware_reservation" "example" {
  id     = "4347e805-eb46-4699-9eb9-5c116e6a0172"
}
```

## Argument Reference

* `id` - (Required) ID of the hardware reservation

## Attributes Reference

* `id` - ID of the hardware reservation to look up
* `short_id` - Reservation short ID
* `project_id` - UUID of project this reservation is scoped to
* `device_id` - UUID of device occupying the reservation
* `plan` - Plan type for the reservation
* `facility` - Plan type for the reservation
* `provisionable` - Flag indicating whether the reservation can be currently used to create a device
* `spare` - Flag indicating whether the reservation is spare (@displague help),
* `switch_uuid` - UUID of switch (@displague help)
* `intervals` - (@displague help)
* `current_period` - (@displague help)
