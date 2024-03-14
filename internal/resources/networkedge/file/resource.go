package file

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/equinix/ne-go"
	"github.com/equinix/rest-go"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// neDeviceSchemaNames is a light copy of the map of Network Edge device schema field names
var neDeviceSchemaNames = map[string]string{
	"IsBYOL":        "byol",
	"IsSelfManaged": "self_managed",
}

var networkFileSchemaNames = map[string]string{
	"UUID":           "uuid",
	"FileName":       "file_name",
	"Content":        "content",
	"MetroCode":      "metro_code",
	"DeviceTypeCode": "device_type_code",
	"ProcessType":    "process_type",
	"IsSelfManaged":  "self_managed",
	"IsBYOL":         "byol",
	"Status":         "status",
}

var networkFileDescriptions = map[string]string{
	"UUID":           "Unique identifier of file resource",
	"FileName":       "File name",
	"Content":        "Uploaded file content, expected to be a UTF-8 encoded string",
	"MetroCode":      "File upload location metro code",
	"DeviceTypeCode": "Device type code",
	"ProcessType":    "File process type (LICENSE or CLOUD_INIT)",
	"IsSelfManaged":  "Boolean value that determines device management mode: self-managed or equinix-managed",
	"IsBYOL":         "Boolean value that determines device licensing mode: bring your own license or subscription",
	"Status":         "File upload status",
}

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkFileCreate,
		ReadContext:   resourceNetworkFileRead,
		DeleteContext: resourceNetworkFileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema:      createNetworkFileSchema(),
		Description: "Resource allows creation and management of Equinix Network Edge device files",
	}
}

func createNetworkFileSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		networkFileSchemaNames["UUID"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkFileDescriptions["UUID"],
		},
		networkFileSchemaNames["FileName"]: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: networkFileDescriptions["FileName"],
		},
		networkFileSchemaNames["Content"]: {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Sensitive:   true,
			Description: networkFileDescriptions["Content"],
		},
		networkFileSchemaNames["MetroCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: equinix_validation.StringIsMetroCode,
			Description:  networkFileDescriptions["MetroCode"],
		},
		networkFileSchemaNames["DeviceTypeCode"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  networkFileDescriptions["DeviceTypeCode"],
		},
		networkFileSchemaNames["ProcessType"]: {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"LICENSE", "CLOUD_INIT"}, false),
			Description:  networkFileDescriptions["ProcessType"],
		},
		networkFileSchemaNames["IsSelfManaged"]: {
			Type:        schema.TypeBool,
			Required:    true,
			ForceNew:    true,
			Description: networkFileDescriptions["IsSelfManaged"],
		},
		networkFileSchemaNames["IsBYOL"]: {
			Type:        schema.TypeBool,
			Required:    true,
			ForceNew:    true,
			Description: networkFileDescriptions["IsBYOL"],
		},
		networkFileSchemaNames["Status"]: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: networkFileDescriptions["Status"],
		},
	}
}

func resourceNetworkFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	fileRequest := createFileRequest(d)
	uuid, err := client.UploadFile(fileRequest["MetroCode"], fileRequest["DeviceTypeCode"], fileRequest["ProcessType"],
		fileRequest["DeviceManagementType"], fileRequest["LicenseMode"], fileRequest["FileName"],
		strings.NewReader(fileRequest["Content"]))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ne.StringValue(uuid))
	diags = append(diags, resourceNetworkFileRead(ctx, d, m)...)
	return diags
}

func resourceNetworkFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*config.Config).Ne
	m.(*config.Config).AddModuleToNEUserAgent(&client, d)
	var diags diag.Diagnostics
	file, err := client.GetFile(d.Id())
	if err != nil {
		if restErr, ok := err.(rest.Error); ok {
			if restErr.HTTPCode == http.StatusNotFound {
				d.SetId("")
				return diags
			}
		}
		return diag.FromErr(err)
	}
	if err := updateFileResource(file, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceNetworkFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func createFileRequest(d *schema.ResourceData) map[string]string {
	fileRequest := make(map[string]string)
	if v, ok := d.GetOk(networkFileSchemaNames["FileName"]); ok {
		fileRequest["FileName"] = v.(string)
	}
	if v, ok := d.GetOk(networkFileSchemaNames["Content"]); ok {
		fileRequest["Content"] = v.(string)
	}
	if v, ok := d.GetOk(networkFileSchemaNames["MetroCode"]); ok {
		fileRequest["MetroCode"] = v.(string)
	}
	if v, ok := d.GetOk(networkFileSchemaNames["DeviceTypeCode"]); ok {
		fileRequest["DeviceTypeCode"] = v.(string)
	}
	if v, ok := d.GetOk(networkFileSchemaNames["ProcessType"]); ok {
		fileRequest["ProcessType"] = v.(string)
	}
	isSelfManaged := d.Get(neDeviceSchemaNames["IsSelfManaged"]).(bool)
	if isSelfManaged {
		fileRequest["DeviceManagementType"] = ne.DeviceManagementTypeSelf
	} else {
		fileRequest["DeviceManagementType"] = ne.DeviceManagementTypeEquinix
	}
	isBYOL := d.Get(neDeviceSchemaNames["IsBYOL"]).(bool)
	if isBYOL {
		fileRequest["LicenseMode"] = ne.DeviceLicenseModeBYOL
	} else {
		fileRequest["LicenseMode"] = ne.DeviceLicenseModeSubscription
	}
	return fileRequest
}

func updateFileResource(file *ne.File, d *schema.ResourceData) error {
	if err := d.Set(networkFileSchemaNames["UUID"], file.UUID); err != nil {
		return fmt.Errorf("error reading %s: %s", networkFileSchemaNames["UUID"], err)
	}
	if err := d.Set(networkFileSchemaNames["FileName"], file.FileName); err != nil {
		return fmt.Errorf("error reading %s: %s", networkFileSchemaNames["FileName"], err)
	}
	if err := d.Set(networkFileSchemaNames["MetroCode"], file.MetroCode); err != nil {
		return fmt.Errorf("error reading %s: %s", networkFileSchemaNames["MetroCode"], err)
	}
	if err := d.Set(networkFileSchemaNames["DeviceTypeCode"], file.DeviceTypeCode); err != nil {
		return fmt.Errorf("error reading %s: %s", networkFileSchemaNames["DeviceTypeCode"], err)
	}
	if err := d.Set(networkFileSchemaNames["ProcessType"], file.ProcessType); err != nil {
		return fmt.Errorf("error reading %s: %s", networkFileSchemaNames["ProcessType"], err)
	}
	if err := d.Set(networkFileSchemaNames["Status"], file.Status); err != nil {
		return fmt.Errorf("error reading %s: %s", networkFileSchemaNames["Status"], err)
	}
	return nil
}
