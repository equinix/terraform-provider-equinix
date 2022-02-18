---
subcategory: "Network Edge"
---

# equinix_network_device_type (Data Source)

Use this data source to get Equinix Network Edge device type details.

## Example Usage

```hcl
# Retrieve device type details of a Cisco router
# Device type has to be available in DC and SV metros
data "equinix_network_device_type" "ciscoRouter" {
  category    = "Router"
  vendor      = "Cisco"
  metro_codes = ["DC", "SV"]
}
```

## Argument Reference

* `name` - (Optional) Device type name
* `vendor` - (Optional) Device type vendor i.e. *"Cisco"*, *"Juniper Networks"*,
*"VERSA Networks"*
* `category` - (Optional) Device type category, one of:
  * Router
  * Firewall
  * SDWAN
* `metro_codes` - (Optional) List of metro codes where device type has to be available

## Attributes Reference

* `code` - Device type short code, unique identifier of a network device type
* `description` - Device type textual description
