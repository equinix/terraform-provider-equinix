package device_software

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var networkDeviceSoftwareSchemaNames = map[string]string{
	"DeviceTypeCode":   "device_type",
	"Version":          "version",
	"VersionRegex":     "version_regex",
	"ImageName":        "image_name",
	"Date":             "date",
	"Status":           "status",
	"IsStable":         "stable",
	"ReleaseNotesLink": "release_notes_link",
	"PackageCodes":     "packages",
	"MostRecent":       "most_recent",
}

var networkDeviceSoftwareDescriptions = map[string]string{
	"DeviceTypeCode":   "Code of a device type",
	"Version":          "Software version",
	"VersionRegex":     "A regex string to apply on returned versions and filter search results",
	"ImageName":        "Software image name",
	"Date":             "Version release date",
	"Status":           "Version status",
	"IsStable":         "Boolean value to limit query results to stable versions only",
	"ReleaseNotesLink": "Link to version release notes",
	"PackageCodes":     "Limits returned versions to those that are supported by given software package codes",
	"MostRecent":       "Boolean value to indicate that most recent version should be used, in case when more than one result is returned",
}

const networkDeviceSoftwareDateLayout = "2006-01-02"

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkDeviceSoftwareRead,
		Description: "Use this data source to get Equinix Network Edge device software details for a given device type.",
		Schema: map[string]*schema.Schema{
			networkDeviceSoftwareSchemaNames["DeviceTypeCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				Description:  networkDeviceSoftwareDescriptions["DeviceTypeCode"],
			},
			networkDeviceSoftwareSchemaNames["Version"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceSoftwareDescriptions["Version"],
			},
			networkDeviceSoftwareSchemaNames["VersionRegex"]: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  networkDeviceSoftwareDescriptions["VersionRegex"],
			},
			networkDeviceSoftwareSchemaNames["ImageName"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceSoftwareDescriptions["ImageName"],
			},
			networkDeviceSoftwareSchemaNames["Date"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceSoftwareDescriptions["Date"],
			},
			networkDeviceSoftwareSchemaNames["Status"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceSoftwareDescriptions["Status"],
			},
			networkDeviceSoftwareSchemaNames["IsStable"]: {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: networkDeviceSoftwareDescriptions["IsStable"],
			},
			networkDeviceSoftwareSchemaNames["ReleaseNotesLink"]: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: networkDeviceSoftwareDescriptions["ReleaseNotesLink"],
			},
			networkDeviceSoftwareSchemaNames["PackageCodes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				Description: networkDeviceSoftwareDescriptions["PackageCodes"],
			},
			networkDeviceSoftwareSchemaNames["MostRecent"]: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: networkDeviceSoftwareDescriptions["MostRecent"],
			},
		},
	}
}

func dataSourceNetworkDeviceSoftwareRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	conf := m.(*config.Config)
	var diags diag.Diagnostics
	typeCode := d.Get(networkDeviceSoftwareSchemaNames["DeviceTypeCode"]).(string)
	pkgCodes := converters.SetToStringList(d.Get(networkDeviceSoftwareSchemaNames["PackageCodes"]).(*schema.Set))
	versions, err := conf.Ne.GetDeviceSoftwareVersions(typeCode)
	if err != nil {
		return diag.FromErr(err)
	}
	var filtered []ne.DeviceSoftwareVersion
	for _, version := range versions {
		if v, ok := d.GetOk(networkDeviceSoftwareSchemaNames["VersionRegex"]); ok {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(ne.StringValue(version.Version)) {
				continue
			}
		}
		if v, ok := d.GetOk(networkDeviceSoftwareSchemaNames["IsStable"]); ok && v.(bool) != ne.BoolValue(version.IsStable) {
			continue
		}
		if !stringsFound(pkgCodes, version.PackageCodes) {
			continue
		}
		filtered = append(filtered, version)
	}
	if len(filtered) < 1 {
		return diag.Errorf("network device software query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		if !d.Get(networkDeviceSoftwareSchemaNames["MostRecent"]).(bool) {
			return diag.Errorf("network device software query returned more than one result, please try more specific search criteria")
		}
		sort.Slice(filtered, func(i, j int) bool {
			iTime, _ := time.Parse(networkDeviceSoftwareDateLayout, ne.StringValue(filtered[i].Date))
			jTime, _ := time.Parse(networkDeviceSoftwareDateLayout, ne.StringValue(filtered[j].Date))
			if iTime.Unix() == jTime.Unix() {
				return ne.StringValue(filtered[i].Version) > ne.StringValue(filtered[j].Version)
			}
			return iTime.Unix() > jTime.Unix()
		})
	}
	if err := updateNetworkDeviceSoftwareResource(filtered[0], typeCode, d); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func updateNetworkDeviceSoftwareResource(version ne.DeviceSoftwareVersion, typeCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", typeCode, ne.StringValue(version.Version)))
	if err := d.Set(networkDeviceSoftwareSchemaNames["Version"], version.Version); err != nil {
		return fmt.Errorf("error reading Version: %s", err)
	}
	if err := d.Set(networkDeviceSoftwareSchemaNames["ImageName"], version.ImageName); err != nil {
		return fmt.Errorf("error reading ImageName: %s", err)
	}
	if err := d.Set(networkDeviceSoftwareSchemaNames["Date"], version.Date); err != nil {
		return fmt.Errorf("error reading Date: %s", err)
	}
	if err := d.Set(networkDeviceSoftwareSchemaNames["Status"], version.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(networkDeviceSoftwareSchemaNames["IsStable"], version.IsStable); err != nil {
		return fmt.Errorf("error reading IsStable: %s", err)
	}
	if err := d.Set(networkDeviceSoftwareSchemaNames["ReleaseNotesLink"], version.ReleaseNotesLink); err != nil {
		return fmt.Errorf("error reading ReleaseNotesLink: %s", err)
	}
	if err := d.Set(networkDeviceSoftwareSchemaNames["PackageCodes"], version.PackageCodes); err != nil {
		return fmt.Errorf("error reading PackageCodes: %s", err)
	}
	return nil
}
