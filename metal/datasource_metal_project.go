package metal

import (
	"fmt"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func dataSourceMetalProject() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalProjectRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Description:   "The name which is used to look up the project",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"project_id"},
			},
			"project_id": {
				Type:          schema.TypeString,
				Description:   "The UUID by which to look up the project",
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
			},

			"created": {
				Type:        schema.TypeString,
				Description: "The timestamp for when the project was created",
				Computed:    true,
			},

			"updated": {
				Type:        schema.TypeString,
				Description: "The timestamp for the last time the project was updated",
				Computed:    true,
			},

			"backend_transfer": {
				Type:        schema.TypeBool,
				Description: "Whether Backend Transfer is enabled for this project",
				Computed:    true,
			},

			"payment_method_id": {
				Type:        schema.TypeString,
				Description: "The UUID of payment method for this project",
				Computed:    true,
			},

			"organization_id": {
				Type:        schema.TypeString,
				Description: "The UUID of this project's parent organization",
				Computed:    true,
			},
			"user_ids": {
				Type:        schema.TypeList,
				Description: "List of UUIDs of user accounts which belong to this project",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"bgp_config": {
				Type:        schema.TypeList,
				Description: "Optional BGP settings. Refer to [Equinix Metal guide for BGP](https://metal.equinix.com/developers/docs/networking/local-global-bgp/)",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_type": {
							Type:         schema.TypeString,
							Description:  "Private or public, the private is likely to be usable immediately, the public will need to be review by Equinix Metal engineers",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"local", "global"}, false),
						},
						"asn": {
							Type:        schema.TypeInt,
							Description: "Autonomous System Number for local BGP deployment",
							Required:    true,
						},
						"md5": {
							Type:        schema.TypeString,
							Description: "Password for BGP session in plaintext (not a checksum)",
							Optional:    true,
							Sensitive:   true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Status of BGP configuration in the project",
							Computed:    true,
						},
						"max_prefix": {
							Type:        schema.TypeInt,
							Description: "The maximum number of route filters allowed per server",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceMetalProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	nameRaw, nameOK := d.GetOk("name")
	projectIdRaw, projectIdOK := d.GetOk("project_id")

	if !projectIdOK && !nameOK {
		return fmt.Errorf("You must supply project_id or name")
	}
	var project *packngo.Project

	if nameOK {
		name := nameRaw.(string)

		os, _, err := client.Projects.List(nil)
		if err != nil {
			return err
		}

		project, err = findProjectByName(os, name)
		if err != nil {
			return err
		}
	} else {
		projectId := projectIdRaw.(string)
		var err error
		project, _, err = client.Projects.Get(projectId, nil)
		if err != nil {
			return err
		}
	}

	d.SetId(project.ID)
	d.Set("payment_method_id", path.Base(project.PaymentMethod.URL))
	d.Set("name", project.Name)
	d.Set("project_id", project.ID)
	d.Set("organization_id", path.Base(project.Organization.URL))
	d.Set("created", project.Created)
	d.Set("updated", project.Updated)
	d.Set("backend_transfer", project.BackendTransfer)

	bgpConf, _, err := client.BGPConfig.Get(project.ID, nil)
	userIds := []string{}
	for _, u := range project.Users {
		userIds = append(userIds, path.Base(u.URL))
	}
	d.Set("user_ids", userIds)

	if (err == nil) && (bgpConf != nil) {
		// guard against an empty struct
		if bgpConf.ID != "" {
			err := d.Set("bgp_config", flattenBGPConfig(bgpConf))
			if err != nil {
				err = friendlyError(err)
				return err
			}
		}
	}
	return nil
}

func findProjectByName(ps []packngo.Project, name string) (*packngo.Project, error) {
	results := make([]packngo.Project, 0)
	for _, p := range ps {
		if p.Name == name {
			results = append(results, p)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no project found with name %s", name)
	}
	return nil, fmt.Errorf("too many projects found with name %s (found %d, expected 1)", name, len(results))
}
