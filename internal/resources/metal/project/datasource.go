package project

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceMetalProjectRead,
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

func dataSourceMetalProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	nameRaw, nameOK := d.GetOk("name")
	projectIdRaw, projectIdOK := d.GetOk("project_id")

	if !projectIdOK && !nameOK {
		return diag.Errorf("you must supply project_id or name")
	}
	var project *metalv1.Project

	if nameOK {
		name := nameRaw.(string)

		projects, err := client.ProjectsApi.FindProjects(ctx).Name(name).ExecuteWithPagination()
		if err != nil {
			return diag.FromErr(err)
		}

		project, err = findProjectByName(projects, name)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		projectId := projectIdRaw.(string)
		var err error
		project, _, err = client.ProjectsApi.FindProjectById(ctx, projectId).Execute()
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(project.GetId())
	d.Set("payment_method_id", path.Base(project.PaymentMethod.GetHref()))
	d.Set("name", project.GetName())
	d.Set("project_id", project.GetId())
	d.Set("organization_id", path.Base(project.Organization.AdditionalProperties["href"].(string))) // spec: organization has no href
	d.Set("created", project.GetCreatedAt().Format(time.RFC3339))
	d.Set("updated", project.GetUpdatedAt().Format(time.RFC3339))
	d.Set("backend_transfer", project.AdditionalProperties["backend_transfer_enabled"].(bool)) // No backend_transfer_enabled property in API spec

	bgpConf, _, err := client.BGPApi.FindBgpConfigByProject(ctx, project.GetId()).Execute()
	userIds := []string{}
	for _, u := range project.GetMembers() {
		userIds = append(userIds, path.Base(u.GetHref()))
	}
	d.Set("user_ids", userIds)

	if (err == nil) && (bgpConf != nil) {
		// guard against an empty struct
		if bgpConf.GetId() != "" {
			err := d.Set("bgp_config", flattenBGPConfig(bgpConf))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

func findProjectByName(ps *metalv1.ProjectList, name string) (*metalv1.Project, error) {
	results := make([]metalv1.Project, 0)
	for _, p := range ps.Projects {
		if p.GetName() == name {
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
