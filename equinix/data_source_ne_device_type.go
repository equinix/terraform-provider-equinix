package equinix

import (
	"fmt"
	"strings"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var neDeviceTypeSchemaNames = map[string]string{
	"Name":        "name",
	"Code":        "code",
	"Description": "description",
	"Vendor":      "vendor",
	"Category":    "category",
	"MetroCodes":  "metro_codes",
}

func dataSourceNeDeviceType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDeviceTypeRead,
		Schema: map[string]*schema.Schema{
			neDeviceTypeSchemaNames["Name"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			neDeviceTypeSchemaNames["Code"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceTypeSchemaNames["Description"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceTypeSchemaNames["Vendor"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			neDeviceTypeSchemaNames["Category"]: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Router", "Firewall", "SDWAN"}, true),
			},
			neDeviceTypeSchemaNames["MetroCodes"]: {
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

func dataSourceDeviceTypeRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	types, err := conf.ne.GetDeviceTypes()
	name := d.Get(neDeviceTypeSchemaNames["Name"]).(string)
	vendor := d.Get(neDeviceTypeSchemaNames["Vendor"]).(string)
	category := d.Get(neDeviceTypeSchemaNames["Category"]).(string)
	metroCodes := expandSetToStringList(d.Get(neDeviceTypeSchemaNames["MetroCodes"]).(*schema.Set))
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
		if !metroCodesFound(metroCodes, deviceType.MetroCodes) {
			continue
		}
		filtered = append(filtered, deviceType)
	}
	if len(filtered) < 1 {
		return fmt.Errorf("device type query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		return fmt.Errorf("device type query returned more than one result, please try more specific search criteria")
	}
	return updateNeDeviceTypeResource(filtered[0], d)
}

func updateNeDeviceTypeResource(deviceType ne.DeviceType, d *schema.ResourceData) error {
	d.SetId(deviceType.Code)
	if err := d.Set(neDeviceTypeSchemaNames["Name"], deviceType.Name); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(neDeviceTypeSchemaNames["Code"], deviceType.Code); err != nil {
		return fmt.Errorf("error reading Code: %s", err)
	}
	if err := d.Set(neDeviceTypeSchemaNames["Description"], deviceType.Description); err != nil {
		return fmt.Errorf("error reading Description: %s", err)
	}
	if err := d.Set(neDeviceTypeSchemaNames["Vendor"], deviceType.Vendor); err != nil {
		return fmt.Errorf("error reading Vendor: %s", err)
	}
	if err := d.Set(neDeviceTypeSchemaNames["Category"], deviceType.Category); err != nil {
		return fmt.Errorf("error reading Category: %s", err)
	}
	if err := d.Set(neDeviceTypeSchemaNames["MetroCodes"], deviceType.MetroCodes); err != nil {
		return fmt.Errorf("error reading MetroCodes: %s", err)
	}
	return nil
}

func metroCodesFound(source []string, target []string) bool {
	for i := range source {
		if !isStringInSlice(source[i], target) {
			return false
		}
	}
	return true
}

func isStringInSlice(needle string, hay []string) bool {
	for i := range hay {
		if needle == hay[i] {
			return true
		}
	}
	return false
}
