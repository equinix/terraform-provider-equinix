package metal

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceOperatingSystem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalOperatingSystemRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "Name or part of the name of the distribution. Case insensitive",
				Optional:    true,
			},
			"distro": {
				Type:        schema.TypeString,
				Description: "Name of the OS distribution",
				Optional:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Version of the distribution",
				Optional:    true,
			},
			"provisionable_on": {
				Type:        schema.TypeString,
				Description: "Plan name",
				Optional:    true,
			},
			"slug": {
				Type:        schema.TypeString,
				Description: "Operating system slug (same as id)",
				Computed:    true,
			},
		},
	}
}

func dataSourceMetalOperatingSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	name, nameOK := d.GetOk("name")
	distro, distroOK := d.GetOk("distro")
	version, versionOK := d.GetOk("version")
	provisionableOn, provisionableOnOK := d.GetOk("provisionable_on")

	if !nameOK && !distroOK && !versionOK && !provisionableOnOK {
		return fmt.Errorf("One of name, distro, version, or provisionable_on must be assigned")
	}

	oss, _, err := client.OperatingSystems.List()
	if err != nil {
		return err
	}

	if nameOK {
		temp := []packngo.OS{}
		for _, os := range oss {
			if strings.Contains(strings.ToLower(os.Name), strings.ToLower(name.(string))) {
				temp = append(temp, os)
			}
		}
		oss = temp
	}

	if distroOK && (len(oss) != 0) {
		temp := []packngo.OS{}
		for _, v := range oss {
			if v.Distro == distro.(string) {
				temp = append(temp, v)
			}
		}
		oss = temp
	}

	if versionOK && (len(oss) != 0) {
		temp := []packngo.OS{}
		for _, v := range oss {
			if v.Version == version.(string) {
				temp = append(temp, v)
			}
		}
		oss = temp
	}

	if provisionableOnOK && (len(oss) != 0) {
		temp := []packngo.OS{}
		for _, v := range oss {
			for _, po := range v.ProvisionableOn {
				if po == provisionableOn.(string) {
					temp = append(temp, v)
				}
			}
		}
		oss = temp
	}

	if len(oss) == 0 {
		return fmt.Errorf("There are no operating systems that match the search criteria")
	}

	if len(oss) > 1 {
		return fmt.Errorf("There is more than one operating system that matches the search criteria")
	}
	d.Set("name", oss[0].Name)
	d.Set("distro", oss[0].Distro)
	d.Set("version", oss[0].Version)
	d.Set("slug", oss[0].Slug)
	d.SetId(oss[0].Slug)
	return nil
}
