---
subcategory: "Metal"
---

# equinix_metal_hardware_reservation (Data Source)

Use this data source to retrieve a [hardware reservation resource from Equinix Metal](https://metal.equinix.com/developers/docs/deploy/reserved/).

You can look up hardware reservation by its ID or by ID of device which occupies it.

## Example Usage

```terraform
// lookup by ID
data "equinix_metal_hardware_reservation" "example" {
  id = "4347e805-eb46-4699-9eb9-5c116e6a0172"
}

// lookup by device ID
data "equinix_metal_hardware_reservation" "example_by_device_id" {
  device_id = "ff85aa58-c106-4624-8f1c-7c64554047ea"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) ID of the hardware reservation.
* `device_id` - (Optional) UUID of device occupying the reservation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the hardware reservation to look up.
* `short_id` - Reservation short ID.
* `project_id` - UUID of project this reservation is scoped to.
* `device_id` - UUID of device occupying the reservation.
* `plan` - Plan type for the reservation.
* `facility` - (**Deprecated**) Facility for the reservation. Use metro instead; read the [facility to metro migration guide](https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices)
* `provisionable` - Flag indicating whether the reserved server is provisionable or not. Spare devices can't be provisioned unless they are activated first.
* `spare` - Flag indicating whether the Hardware Reservation is a spare. Spare Hardware Reservations are used when a Hardware Reservations requires service from Metal Equinix.
* `switch_uuid` - Switch short ID, can be used to determine if two devices are connected to the same switch.
