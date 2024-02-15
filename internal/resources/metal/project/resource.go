package project

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMetalProjectCreate,
		ReadWithoutTimeout:   resourceMetalProjectRead,
		UpdateWithoutTimeout: resourceMetalProjectUpdate,
		DeleteWithoutTimeout: resourceMetalProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the project.  The maximum length is 80 characters.",
				Required:    true,
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
				Description: "Enable or disable [Backend Transfer](https://metal.equinix.com/developers/docs/networking/backend-transfer/), default is false",
				Optional:    true,
				Default:     false,
			},
			"payment_method_id": {
				Type:        schema.TypeString,
				Description: "The UUID of payment method for this project. The payment method and the project need to belong to the same organization (passed with organization_id, or default)",
				Optional:    true,
				Computed:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(strings.Trim(old, `"`), strings.Trim(new, `"`))
				},
				ValidateFunc: validation.Any(
					validation.IsUUID,
					validation.StringIsEmpty,
				),
			},
			"organization_id": {
				Type:        schema.TypeString,
				Description: "The UUID of organization under which you want to create the project. If you leave it out, the project will be create under your the default organization of your account",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(strings.Trim(old, `"`), strings.Trim(new, `"`))
				},
				ValidateFunc: validation.IsUUID,
			},
			"bgp_config": {
				Type:        schema.TypeList,
				Description: "Optional BGP settings. Refer to [Equinix Metal guide for BGP](https://metal.equinix.com/developers/docs/networking/local-global-bgp/)",
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_type": {
							Type:         schema.TypeString,
							Description:  "\"local\" or \"global\", the local is likely to be usable immediately, the global will need to be review by Equinix Metal engineers",
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
							Sensitive:   true,
							Optional:    true,
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

func expandBGPConfig(d *schema.ResourceData) (*metalv1.BgpConfigRequestInput, error) {
	bgpDeploymentType, err := metalv1.NewBgpConfigRequestInputDeploymentTypeFromValue(d.Get("bgp_config.0.deployment_type").(string))
	if err != nil {
		return nil, err
	}

	bgpCreateRequest := metalv1.BgpConfigRequestInput{
		DeploymentType: *bgpDeploymentType,
		Asn:            int32(d.Get("bgp_config.0.asn").(int)),
	}
	md5, ok := d.GetOk("bgp_config.0.md5")
	if ok {
		bgpCreateRequest.Md5 = metalv1.PtrString(md5.(string))
	}

	return &bgpCreateRequest, nil
}

func resourceMetalProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	createRequest := metalv1.ProjectCreateFromRootInput{
		Name: d.Get("name").(string),
	}

	organization_id := d.Get("organization_id").(string)
	if organization_id != "" {
		createRequest.OrganizationId = &organization_id
	}

	project, resp, err := client.ProjectsApi.CreateProject(ctx).ProjectCreateFromRootInput(createRequest).Execute()

	if err != nil {
		return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
	}

	d.SetId(project.GetId())

	_, hasBGPConfig := d.GetOk("bgp_config")
	if hasBGPConfig {
		bgpCR, err := expandBGPConfig(d)
		if err == nil {
			resp, err = client.BGPApi.RequestBgpConfig(ctx, project.GetId()).BgpConfigRequestInput(*bgpCR).Execute()
		}
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
		}
	}

	backendTransfer := d.Get("backend_transfer").(bool)
	if backendTransfer {
		pur := metalv1.ProjectUpdateInput{
			BackendTransferEnabled: &backendTransfer,
		}
		_, _, err := client.ProjectsApi.UpdateProject(ctx, project.GetId()).ProjectUpdateInput(pur).Execute()
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
		}
	}
	return resourceMetalProjectRead(ctx, d, meta)
}

func resourceMetalProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	proj, resp, err := client.ProjectsApi.FindProjectById(ctx, d.Id()).Execute()
	if err != nil {
		err = equinix_errors.FriendlyErrorForMetalGo(err, resp)

		// If the project somehow already destroyed, mark as successfully gone.
		if equinix_errors.IsNotFound(err) {
			d.SetId("")

			return nil
		}

		return diag.FromErr(err)
	}

	d.SetId(proj.GetId())
	if len(proj.PaymentMethod.GetHref()) != 0 {
		d.Set("payment_method_id", path.Base(proj.PaymentMethod.GetHref()))
	}
	d.Set("name", proj.Name)
	d.Set("organization_id", path.Base(proj.Organization.AdditionalProperties["href"].(string))) // spec: organization has no href
	d.Set("created", proj.GetCreatedAt().Format(time.RFC3339))
	d.Set("updated", proj.GetUpdatedAt().Format(time.RFC3339))
	d.Set("backend_transfer", proj.AdditionalProperties["backend_transfer_enabled"].(bool)) // No backend_transfer_enabled property in API spec

	bgpConf, _, err := client.BGPApi.FindBgpConfigByProject(ctx, proj.GetId()).Execute()

	if (err == nil) && (bgpConf != nil) {
		// guard against an empty struct
		if bgpConf.GetId() != "" {
			err := d.Set("bgp_config", flattenBGPConfig(bgpConf))
			if err != nil {
				return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
			}
		}
	}
	return nil
}

func flattenBGPConfig(l *metalv1.BgpConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if l == nil {
		return nil
	}

	r := make(map[string]interface{})

	if l.GetStatus() != "" {
		r["status"] = l.GetStatus()
	}
	if l.GetDeploymentType() != "" {
		r["deployment_type"] = l.GetDeploymentType()
	}
	if l.GetMd5() != "" {
		r["md5"] = l.GetMd5()
	}
	if l.GetAsn() != 0 {
		r["asn"] = l.GetAsn()
	}
	if l.GetMaxPrefix() != 0 {
		r["max_prefix"] = l.GetMaxPrefix()
	}

	result = append(result, r)

	return result
}

func resourceMetalProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	updateRequest := metalv1.ProjectUpdateInput{}
	if d.HasChange("name") {
		pName := d.Get("name").(string)
		updateRequest.Name = &pName
	}
	if d.HasChange("payment_method_id") {
		pPayment := d.Get("payment_method_id").(string)
		updateRequest.PaymentMethodId = &pPayment
	}
	if d.HasChange("backend_transfer") {
		pBT := d.Get("backend_transfer").(bool)
		updateRequest.BackendTransferEnabled = &pBT
	}
	if d.HasChange("bgp_config") {
		o, n := d.GetChange("bgp_config")
		oldarr := o.([]interface{})
		newarr := n.([]interface{})
		if len(newarr) == 1 {
			bgpCreateRequest, err := expandBGPConfig(d)
			if err != nil {
				return diag.FromErr(err)
			}

			resp, err := client.BGPApi.RequestBgpConfig(ctx, d.Id()).BgpConfigRequestInput(*bgpCreateRequest).Execute()
			if err != nil {
				return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
			}
		} else {
			if len(oldarr) == 1 {
				m := oldarr[0].(map[string]interface{})

				bgpConfStr := fmt.Sprintf(
					"bgp_config {\n"+
						"  deployment_type = \"%s\"\n"+
						"  md5 = \"%s\"\n"+
						"  asn = %d\n"+
						"}", m["deployment_type"].(string), m["md5"].(string),
					m["asn"].(int))

				return diag.Errorf("BGP Config can not be removed from a project, please add back\n%s", bgpConfStr)
			}
		}
	} else {
		_, resp, err := client.ProjectsApi.UpdateProject(ctx, d.Id()).ProjectUpdateInput(updateRequest).Execute()
		if err != nil {
			return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
		}
	}

	return resourceMetalProjectRead(ctx, d, meta)
}

func resourceMetalProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.Config).NewMetalClientForSDK(d)

	resp, err := client.ProjectsApi.DeleteProject(ctx, d.Id()).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return diag.FromErr(equinix_errors.FriendlyErrorForMetalGo(err, resp))
	}

	d.SetId("")
	return nil
}
