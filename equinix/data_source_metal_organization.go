package equinix

import (
	"fmt"
	"path"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalOrganization() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalOrganizationRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Description:   "The organization name",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"organization_id"},
			},
			"organization_id": {
				Type:          schema.TypeString,
				Description:   "The UUID of the organization resource",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description string",
				Computed:    true,
			},
			"website": {
				Type:        schema.TypeString,
				Description: "Website link",
				Computed:    true,
			},
			"twitter": {
				Type:        schema.TypeString,
				Description: "Twitter handle",
				Computed:    true,
			},
			"logo": {
				Type:        schema.TypeString,
				Description: "Logo URL",
				Computed:    true,
			},
			"project_ids": {
				Type:        schema.TypeList,
				Description: "UUIDs of project resources which belong to this organization",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"address": {
				Type:        schema.TypeList,
				Description: "Business' address",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: createOrganizationAddressDataSourceSchema(),
				},
			},
		},
	}
}

func createOrganizationAddressDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"city": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"zip_code": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"country": {
			Type:        schema.TypeString,
			Description: "Two letter country code (ISO 3166-1 alpha-2), e.g. US",
			Computed:    true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func findOrgByName(os []packngo.Organization, name string) (*packngo.Organization, error) {
	results := make([]packngo.Organization, 0)
	for _, o := range os {
		if o.Name == name {
			results = append(results, o)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no organization found with name %s", name)
	}
	return nil, fmt.Errorf("too many organizations found with name %s (found %d, expected 1)", name, len(results))
}

func dataSourceMetalOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*config.Config).Metal
	nameRaw, nameOK := d.GetOk("name")
	orgIdRaw, orgIdOK := d.GetOk("organization_id")

	if !orgIdOK && !nameOK {
		return fmt.Errorf("you must supply organization_id or name")
	}
	var org *packngo.Organization

	if nameOK {
		name := nameRaw.(string)

		os, _, err := client.Organizations.List(&packngo.GetOptions{Includes: []string{"address"}})
		if err != nil {
			return err
		}

		org, err = findOrgByName(os, name)
		if err != nil {
			return err
		}
	} else {
		orgId := orgIdRaw.(string)
		var err error

		org, _, err = client.Organizations.Get(orgId, &packngo.GetOptions{Includes: []string{"address"}})
		if err != nil {
			return err
		}
	}
	projectIds := []string{}

	for _, p := range org.Projects {
		projectIds = append(projectIds, path.Base(p.URL))
	}

	d.SetId(org.ID)
	return setMap(d, map[string]interface{}{
		"organization_id": org.ID,
		"name":            org.Name,
		"description":     org.Description,
		"website":         org.Website,
		"twitter":         org.Twitter,
		"logo":            org.Logo,
		"project_ids":     projectIds,
		"address":         flattenMetalOrganizationAddress(org.Address),
	})
}
