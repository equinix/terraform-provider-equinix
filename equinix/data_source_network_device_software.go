package equinix

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/equinix/ne-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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

const networkDeviceSoftwareDateLayout = "2006-01-02"

func dataSourceNetworkDeviceSoftware() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkDeviceSoftwareRead,
		Schema: map[string]*schema.Schema{
			networkDeviceSoftwareSchemaNames["DeviceTypeCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			networkDeviceSoftwareSchemaNames["Version"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDeviceSoftwareSchemaNames["VersionRegex"]: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			networkDeviceSoftwareSchemaNames["ImageName"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDeviceSoftwareSchemaNames["Date"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDeviceSoftwareSchemaNames["Status"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			networkDeviceSoftwareSchemaNames["IsStable"]: {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			networkDeviceSoftwareSchemaNames["ReleaseNotesLink"]: {
				Type:     schema.TypeString,
				Computed: true,
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
			},
			networkDeviceSoftwareSchemaNames["MostRecent"]: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceNetworkDeviceSoftwareRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	typeCode := d.Get(networkDeviceSoftwareSchemaNames["DeviceTypeCode"]).(string)
	pkgCodes := expandSetToStringList(d.Get(networkDeviceSoftwareSchemaNames["PackageCodes"]).(*schema.Set))
	versions, err := conf.ne.GetDeviceSoftwareVersions(typeCode)
	if err != nil {
		return err
	}
	var filtered []ne.DeviceSoftwareVersion
	for _, version := range versions {
		if v, ok := d.GetOk(networkDeviceSoftwareSchemaNames["VersionRegex"]); ok {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(version.Version) {
				continue
			}
		}
		if v, ok := d.GetOk(networkDeviceSoftwareSchemaNames["IsStable"]); ok && v.(bool) == version.IsStable {
			continue
		}
		if !stringsFound(pkgCodes, version.PackageCodes) {
			continue
		}
		filtered = append(filtered, version)
	}
	if len(filtered) < 1 {
		return fmt.Errorf("network device software query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		if !d.Get(networkDeviceSoftwareSchemaNames["MostRecent"]).(bool) {
			return fmt.Errorf("network device software query returned more than one result, please try more specific search criteria")
		}
		sort.Slice(filtered, func(i, j int) bool {
			iTime, _ := time.Parse(networkDeviceSoftwareDateLayout, filtered[i].Date)
			jTime, _ := time.Parse(networkDeviceSoftwareDateLayout, filtered[j].Date)
			return iTime.Unix() > jTime.Unix()
		})
	}
	return updateNetworkDeviceSoftwareResource(filtered[0], typeCode, d)
}

func updateNetworkDeviceSoftwareResource(version ne.DeviceSoftwareVersion, typeCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", typeCode, version.Version))
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
