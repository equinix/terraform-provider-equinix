// Package devicetype provides the network_device_type data source
package devicetype

import (
	"context"
	"fmt"
	"strings"

	"github.com/equinix/terraform-provider-equinix/internal/comparisons"
	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_validation "github.com/equinix/terraform-provider-equinix/internal/validation"

	"github.com/equinix/equinix-sdk-go/services/networkedgev1"
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

// DataSource creates a new Terraform data source for retrieving device type data
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
	var diags diag.Diagnostics
	client := m.(*config.Config).NewNetworkEdgeClientForSDK(ctx, d)
	types, _, err := client.SetupApi.GetVirtualDevicesUsingGET(ctx).Execute()
	name := d.Get(networkDeviceTypeSchemaNames["Name"]).(string)
	vendor := d.Get(networkDeviceTypeSchemaNames["Vendor"]).(string)
	category := d.Get(networkDeviceTypeSchemaNames["Category"]).(string)
	metroCodes := converters.SetToStringList(d.Get(networkDeviceTypeSchemaNames["MetroCodes"]).(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}
	filtered := make([]networkedgev1.VirtualDeviceType, 0, len(types.Data))
	for _, deviceType := range types.Data {
		if name != "" && deviceType.GetName() != name {
			continue
		}
		if vendor != "" && deviceType.GetVendor() != vendor {
			continue
		}
		if category != "" && !strings.EqualFold(deviceType.GetCategory(), category) {
			continue
		}

		availableMetros := make([]string, len(deviceType.AvailableMetros))
		for _, metro := range deviceType.AvailableMetros {
			availableMetros = append(availableMetros, metro.GetMetroCode())
		}
		if !comparisons.Subsets(metroCodes, availableMetros) {
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

func updateNetworkDeviceTypeResource(deviceType networkedgev1.VirtualDeviceType, d *schema.ResourceData) error {
	d.SetId(deviceType.GetDeviceTypeCode())
	if err := d.Set(networkDeviceTypeSchemaNames["Name"], deviceType.GetName()); err != nil {
		return fmt.Errorf("error reading Name: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Code"], deviceType.GetDeviceTypeCode()); err != nil {
		return fmt.Errorf("error reading Code: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Description"], deviceType.GetDescription()); err != nil {
		return fmt.Errorf("error reading Description: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Vendor"], deviceType.GetVendor()); err != nil {
		return fmt.Errorf("error reading Vendor: %s", err)
	}
	if err := d.Set(networkDeviceTypeSchemaNames["Category"], deviceType.GetCategory()); err != nil {
		return fmt.Errorf("error reading Category: %s", err)
	}

	availableMetros := make([]string, len(deviceType.AvailableMetros))
	for _, metro := range deviceType.AvailableMetros {
		availableMetros = append(availableMetros, metro.GetMetroCode())
	}
	if err := d.Set(networkDeviceTypeSchemaNames["MetroCodes"], availableMetros); err != nil {
		return fmt.Errorf("error reading MetroCodes: %s", err)
	}
	return nil
}
