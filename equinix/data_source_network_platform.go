package equinix

import (
	"context"
	"fmt"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkDevicePlatformSchemaNames = map[string]string{
	"DeviceTypeCode":  "device_type",
	"Flavor":          "flavor",
	"CoreCount":       "core_count",
	"Memory":          "memory",
	"MemoryUnit":      "memory_unit",
	"PackageCodes":    "packages",
	"ManagementTypes": "management_types",
	"LicenseOptions":  "license_options",
}

func dataSourceNetworkDevicePlatform() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkDevicePlatformRead,
		Schema: map[string]*schema.Schema{
			networkDevicePlatformSchemaNames["DeviceTypeCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			networkDevicePlatformSchemaNames["Flavor"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"small", "medium", "large", "xlarge"}, false),
			},
			networkDevicePlatformSchemaNames["CoreCount"]: {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			networkDevicePlatformSchemaNames["Memory"]: {
				Type:     schema.TypeInt,
				Computed: true,
			},
			networkDevicePlatformSchemaNames["MemoryUnit"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDevicePlatformSchemaNames["PackageCodes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			networkDevicePlatformSchemaNames["ManagementTypes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"EQUINIX-CONFIGURED", "SELF-CONFIGURED"}, false),
				},
			},
			networkDevicePlatformSchemaNames["LicenseOptions"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"BYOL", "Sub"}, false),
				},
			},
		},
	}
}

func dataSourceNetworkDevicePlatformRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*Config)
	var diags diag.Diagnostics
	typeCode := d.Get(networkDevicePlatformSchemaNames["DeviceTypeCode"]).(string)
	platforms, err := conf.ne.GetDevicePlatforms(typeCode)
	if err != nil {
		return diag.FromErr(err)
	}
	var filtered []ne.DevicePlatform
	for _, platform := range platforms {
		if v, ok := d.GetOk(networkDevicePlatformSchemaNames["Flavor"]); ok && ne.StringValue(platform.Flavor) != v.(string) {
			continue
		}
		if v, ok := d.GetOk(networkDevicePlatformSchemaNames["CoreCount"]); ok && ne.IntValue(platform.CoreCount) != v.(int) {
			continue
		}
		if v, ok := d.GetOk(networkDevicePlatformSchemaNames["PackageCodes"]); ok {
			pkgCodes := expandSetToStringList(v.(*schema.Set))
			if !stringsFound(pkgCodes, platform.PackageCodes) {
				continue
			}
		}
		if v, ok := d.GetOk(networkDevicePlatformSchemaNames["ManagementTypes"]); ok {
			mgmtTypes := expandSetToStringList(v.(*schema.Set))
			if !stringsFound(mgmtTypes, platform.ManagementTypes) {
				continue
			}
		}
		if v, ok := d.GetOk(networkDevicePlatformSchemaNames["LicenseOptions"]); ok {
			licOptions := expandSetToStringList(v.(*schema.Set))
			if !stringsFound(licOptions, platform.LicenseOptions) {
				continue
			}
		}
		filtered = append(filtered, platform)
	}
	if len(filtered) < 1 {
		return diag.Errorf("network device platform query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return diag.Errorf("network device platform query returned more than one result, please try more specific search criteria")
	}
	if err := updateNetworkDevicePlatformResource(filtered[0], typeCode, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateNetworkDevicePlatformResource(platform ne.DevicePlatform, typeCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", typeCode, ne.StringValue(platform.Flavor)))
	if err := d.Set(networkDevicePlatformSchemaNames["Flavor"], platform.Flavor); err != nil {
		return fmt.Errorf("error reading Flavor: %s", err)
	}
	if err := d.Set(networkDevicePlatformSchemaNames["CoreCount"], platform.CoreCount); err != nil {
		return fmt.Errorf("error reading CoreCount: %s", err)
	}
	if err := d.Set(networkDevicePlatformSchemaNames["Memory"], platform.Memory); err != nil {
		return fmt.Errorf("error reading Memory: %s", err)
	}
	if err := d.Set(networkDevicePlatformSchemaNames["MemoryUnit"], platform.MemoryUnit); err != nil {
		return fmt.Errorf("error reading MemoryUnit: %s", err)
	}
	if err := d.Set(networkDevicePlatformSchemaNames["PackageCodes"], platform.PackageCodes); err != nil {
		return fmt.Errorf("error reading PackageCodes: %s", err)
	}
	if err := d.Set(networkDevicePlatformSchemaNames["ManagementTypes"], platform.ManagementTypes); err != nil {
		return fmt.Errorf("error reading ManagementTypes: %s", err)
	}
	if err := d.Set(networkDevicePlatformSchemaNames["LicenseOptions"], platform.LicenseOptions); err != nil {
		return fmt.Errorf("error reading LicenseOptions: %s", err)
	}
	return nil
}
