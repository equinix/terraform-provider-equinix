package equinix

import (
	"fmt"
	"strings"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var networkDeviceTypeSchemaNames = map[string]string{
	"Name":        "name",
	"Code":        "code",
	"Description": "description",
	"Vendor":      "vendor",
	"Category":    "category",
	"MetroCodes":  "metro_codes",
}

func dataSourceNetworkDeviceType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkDeviceTypeRead,
		Schema: map[string]*schema.Schema{
			networkDeviceTypeSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			networkDeviceTypeSchemaNames["Code"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDeviceTypeSchemaNames["Description"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDeviceTypeSchemaNames["Vendor"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			networkDeviceTypeSchemaNames["Category"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Router", "Firewall", "SDWAN"}, true),
			},
			networkDeviceTypeSchemaNames["MetroCodes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: stringIsMetroCode(),
				},
			},
		},
	}
}

func dataSourceNetworkDeviceTypeRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	types, err := conf.ne.GetDeviceTypes()
	name := d.Get(networkDeviceTypeSchemaNames["Name"]).(string)
	vendor := d.Get(networkDeviceTypeSchemaNames["Vendor"]).(string)
	category := d.Get(networkDeviceTypeSchemaNames["Category"]).(string)
	metroCodes := expandSetToStringList(d.Get(networkDeviceTypeSchemaNames["MetroCodes"]).(*schema.Set))
	if err != nil {
		return err
	}
	filtered := make([]ne.DeviceType, 0, len(types))
	for _, deviceType := range types {
		if name != "" && deviceType.Name != name {
			continue
		}
		if vendor != "" && deviceType.Vendor != vendor {
			continue
		}
		if category != "" && !strings.EqualFold(deviceType.Category, category) {
			continue
		}
		if !stringsFound(metroCodes, deviceType.MetroCodes) {
			continue
		}
		filtered = append(filtered, deviceType)
	}
	if len(filtered) < 1 {
		return fmt.Errorf("network device type query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("network device type query returned more than one result, please try more specific search criteria")
	}
	return updateNetworkDeviceTypeResource(filtered[0], d)
}

func updateNetworkDeviceTypeResource(deviceType ne.DeviceType, d *schema.ResourceData) error {
	d.SetId(deviceType.Code)
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
