---
subcategory: "Network Edge"
---

# equinix_network_file (Resource)

Resource `equinix_network_file` allows creation and management of Equinix Network Edge files.

## Example Usage

```hcl
variable "filepath" { default = "fileFolder/fileName.txt" }

resource "equinix_network_file" "test-file" {
  file_name = "fileName.txt"
  content = file("${path.module}/${var.filepath}")
  metro_code = "SV"
  device_type_code = "AVIATRIX_EDGE"
  process_type = "CLOUD_INIT"
  self_managed = true
  byol = true
}
```

## Argument Reference

The following arguments are supported:

* `file_name` - (Required) File name.
* `content` - (Required) Uploaded file content, expected to be a UTF-8 encoded string.
* `metro_code` - (Required) File upload location metro code. It should match the device location metro code.
* `type_code` - (Required) Device type code.
* `process_type` - (Required) File process type (LICENSE or CLOUD_INIT).
* `self_managed` - (Required) Boolean value that determines device management mode, i.e.,
  `self-managed` or `Equinix-managed`.
* `byol` - (Required) Boolean value that determines device licensing mode, i.e.,
  `bring your own license` or `subscription`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `uuid` - Unique identifier of file resource.
* `status` - File upload status.

## Import

This resource can be imported using an existing ID:

```sh
terraform import equinix_network_file.example {existing_id}
```

The `content`, `self_managed` and `byol` fields can not be imported.
