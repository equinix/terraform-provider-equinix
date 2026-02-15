package device_type

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkDeviceTypeSchemaNames = map[string]string{
	"Name":        "name",
	"Code":        "code",
	"Description": "description",
	"Vendor":      "vendor",
	"Category":    "category",
	"MetroCodes":  "metro_codes",
}

var networkDeviceTypeDescriptions = map[string]string{
	"Name":        "Device type name",
	"Code":        "Device type short code, unique identifier of a network device type",
	"Description": "Device type textual description",
	"Vendor":      "Device type vendor i.e. Cisco, Juniper Networks, VERSA Networks",
	"Category":    "Device type category, one of: Router, Firewall, SDWAN",
	"MetroCodes":  "List of metro codes where device type has to be available",
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkDeviceTypeRead,
		Description: "Use this data source to get Equinix Network Edge device type details",
		Schema: map[string]*schema.Schema{
			networkDeviceTypeSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  networkDeviceTypeDescriptions["Name"],
			},
			networkDeviceTypeSchemaNames["Code"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceTypeDescriptions["Code"],
			},
			networkDeviceTypeSchemaNames["Description"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceTypeDescriptions["Description"],
			},
			networkDeviceTypeSchemaNames["Vendor"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  networkDeviceTypeDescriptions["Vendor"],
			},
			networkDeviceTypeSchemaNames["Category"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Router", "Firewall", "SDWAN"}, true),
				Description:  networkDeviceTypeDescriptions["Category"],
			},
			networkDeviceTypeSchemaNames["MetroCodes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: equinix_validation.StringIsMetroCode,
				},
				Description: networkDeviceTypeDescriptions["MetroCodes"],
			},
		},
	}
}

func dataSourceNetworkDeviceTypeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*config.Config)
	var diags diag.Diagnostics
	types, err := conf.Ne.GetDeviceTypes()
	name := d.Get(networkDeviceTypeSchemaNames["Name"]).(string)
	vendor := d.Get(networkDeviceTypeSchemaNames["Vendor"]).(string)
	category := d.Get(networkDeviceTypeSchemaNames["Category"]).(string)
	metroCodes := converters.SetToStringList(d.Get(networkDeviceTypeSchemaNames["MetroCodes"]).(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}
	filtered := make([]ne.DeviceType, 0, len(types))
	for _, deviceType := range types {
		if name != "" && ne.StringValue(deviceType.Name) != name {
			continue
		}
		if vendor != "" && ne.StringValue(deviceType.Vendor) != vendor {
			continue
		}
		if category != "" && !strings.EqualFold(ne.StringValue(deviceType.Category), category) {
			continue
		}
		if !comparisons.Subsets(metroCodes, deviceType.MetroCodes) {
			continue
		}
		filtered = append(filtered, deviceType)
	}
	if len(filtered) < 1 {
		return diag.Errorf("network device type query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return diag.Errorf("network device type query returned more than one result, please try more specific search criteria")
	}
	if err := updateNetworkDeviceTypeResource(filtered[0], d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateNetworkDeviceTypeResource(deviceType ne.DeviceType, d *schema.ResourceData) error {
	d.SetId(ne.StringValue(deviceType.Code))
	if err := d.Set(networkDeviceTypeSchemaNames["Name"], deviceType.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Code"], deviceType.Code); err != nil {
		return fmt.Errorf("error reading Code: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Description"], deviceType.Description); err != nil {
		return fmt.Errorf("error reading Description: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Vendor"], deviceType.Vendor); err != nil {
		return fmt.Errorf("error reading Vendor: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Category"], deviceType.Category); err != nil {
		return fmt.Errorf("error reading Category: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["MetroCodes"], deviceType.MetroCodes); err != nil {
		return fmt.Errorf("error reading MetroCodes: %s", err)
	}
	return nil
}
