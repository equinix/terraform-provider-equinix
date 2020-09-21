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

var neDeviceSoftwareSchemaNames = map[string]string{
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

const dateLayout = "2006-01-02"

func dataSourceNeDeviceSoftware() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNeDeviceSoftwareRead,
		Schema: map[string]*schema.Schema{
			neDeviceSoftwareSchemaNames["DeviceTypeCode"]: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			neDeviceSoftwareSchemaNames["Version"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceSoftwareSchemaNames["VersionRegex"]: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			neDeviceSoftwareSchemaNames["ImageName"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceSoftwareSchemaNames["Date"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceSoftwareSchemaNames["Status"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceSoftwareSchemaNames["IsStable"]: {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			neDeviceSoftwareSchemaNames["ReleaseNotesLink"]: {
				Type:     schema.TypeString,
				Computed: true,
			},
			neDeviceSoftwareSchemaNames["PackageCodes"]: {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			neDeviceSoftwareSchemaNames["MostRecent"]: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceNeDeviceSoftwareRead(d *schema.ResourceData, m interface{}) error {
	conf := m.(*Config)
	typeCode := d.Get(neDeviceSoftwareSchemaNames["DeviceTypeCode"]).(string)
	pkgCodes := expandSetToStringList(d.Get(neDeviceSoftwareSchemaNames["PackageCodes"]).(*schema.Set))
	versions, err := conf.ne.GetDeviceSoftwareVersions(typeCode)
	if err != nil {
		return err
	}
	var filtered []ne.DeviceSoftwareVersion
	for _, version := range versions {
		if v, ok := d.GetOk(neDeviceSoftwareSchemaNames["VersionRegex"]); ok {
			r := regexp.MustCompile(v.(string))
			if !r.MatchString(version.Version) {
				continue
			}
		}
		if v, ok := d.GetOk(neDeviceSoftwareSchemaNames["IsStable"]); ok && v.(bool) == version.IsStable {
			continue
		}
		if !stringsFound(pkgCodes, version.PackageCodes) {
			continue
		}
		filtered = append(filtered, version)
	}
	if len(filtered) < 1 {
		return fmt.Errorf("device software query returned no results, please change your search criteria")
	}
	if len(filtered) > 1 {
		if !d.Get(neDeviceSoftwareSchemaNames["MostRecent"]).(bool) {
			return fmt.Errorf("device software query returned more than one result, please try more specific search criteria")
		}
		sort.Slice(filtered, func(i, j int) bool {
			iTime, _ := time.Parse(dateLayout, filtered[i].Date)
			jTime, _ := time.Parse(dateLayout, filtered[j].Date)
			return iTime.Unix() > jTime.Unix()
		})
	}
	return updateNeDeviceSoftwareResource(filtered[0], typeCode, d)
}

func updateNeDeviceSoftwareResource(version ne.DeviceSoftwareVersion, typeCode string, d *schema.ResourceData) error {
	d.SetId(fmt.Sprintf("%s-%s", typeCode, version.Version))
	if err := d.Set(neDeviceSoftwareSchemaNames["Version"], version.Version); err != nil {
		return fmt.Errorf("error reading Version: %s", err)
	}
	if err := d.Set(neDeviceSoftwareSchemaNames["ImageName"], version.ImageName); err != nil {
		return fmt.Errorf("error reading ImageName: %s", err)
	}
	if err := d.Set(neDeviceSoftwareSchemaNames["Date"], version.Date); err != nil {
		return fmt.Errorf("error reading Date: %s", err)
	}
	if err := d.Set(neDeviceSoftwareSchemaNames["Status"], version.Status); err != nil {
		return fmt.Errorf("error reading Status: %s", err)
	}
	if err := d.Set(neDeviceSoftwareSchemaNames["IsStable"], version.IsStable); err != nil {
		return fmt.Errorf("error reading IsStable: %s", err)
	}
	if err := d.Set(neDeviceSoftwareSchemaNames["ReleaseNotesLink"], version.ReleaseNotesLink); err != nil {
		return fmt.Errorf("error reading ReleaseNotesLink: %s", err)
	}
	if err := d.Set(neDeviceSoftwareSchemaNames["PackageCodes"], version.PackageCodes); err != nil {
		return fmt.Errorf("error reading PackageCodes: %s", err)
	}
	return nil
}
