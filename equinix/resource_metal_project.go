package equinix

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

var uuidRE = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func resourceMetalProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalProjectCreate,
		Read:   resourceMetalProjectRead,
		Update: resourceMetalProjectUpdate,
		Delete: resourceMetalProjectDelete,
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
					validation.StringMatch(uuidRE, "must be a valid UUID"),
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
				ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
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

func expandBGPConfig(d *schema.ResourceData) packngo.CreateBGPConfigRequest {
	bgpCreateRequest := packngo.CreateBGPConfigRequest{
		DeploymentType: d.Get("bgp_config.0.deployment_type").(string),
		Asn:            d.Get("bgp_config.0.asn").(int),
	}
	md5, ok := d.GetOk("bgp_config.0.md5")
	if ok {
		bgpCreateRequest.Md5 = md5.(string)
	}

	return bgpCreateRequest
}

func resourceMetalProjectCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

	createRequest := &packngo.ProjectCreateRequest{
		Name:           d.Get("name").(string),
		OrganizationID: d.Get("organization_id").(string),
	}

	project, _, err := client.Projects.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(project.ID)

	_, hasBGPConfig := d.GetOk("bgp_config")
	if hasBGPConfig {
		bgpCR := expandBGPConfig(d)
		_, err := client.BGPConfig.Create(project.ID, bgpCR)
		if err != nil {
			return friendlyError(err)
		}
	}

	backendTransfer := d.Get("backend_transfer").(bool)
	if backendTransfer {
		pur := packngo.ProjectUpdateRequest{BackendTransfer: &backendTransfer}
		_, _, err := client.Projects.Update(project.ID, &pur)
		if err != nil {
			return friendlyError(err)
		}
	}
	return resourceMetalProjectRead(d, meta)
}

func resourceMetalProjectRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

	proj, _, err := client.Projects.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as successfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(proj.ID)
	if len(proj.PaymentMethod.URL) != 0 {
		d.Set("payment_method_id", path.Base(proj.PaymentMethod.URL))
	}
	d.Set("name", proj.Name)
	d.Set("organization_id", path.Base(proj.Organization.URL))
	d.Set("created", proj.Created)
	d.Set("updated", proj.Updated)
	d.Set("backend_transfer", proj.BackendTransfer)

	bgpConf, _, err := client.BGPConfig.Get(proj.ID, nil)

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

func flattenBGPConfig(l *packngo.BGPConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)

	if l == nil {
		return nil
	}

	r := make(map[string]interface{})

	if l.Status != "" {
		r["status"] = l.Status
	}
	if l.DeploymentType != "" {
		r["deployment_type"] = l.DeploymentType
	}
	if l.Md5 != "" {
		r["md5"] = l.Md5
	}
	if l.Asn != 0 {
		r["asn"] = l.Asn
	}
	if l.MaxPrefix != 0 {
		r["max_prefix"] = l.MaxPrefix
	}

	result = append(result, r)

	return result
}

func resourceMetalProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal
	updateRequest := &packngo.ProjectUpdateRequest{}
	if d.HasChange("name") {
		pName := d.Get("name").(string)
		updateRequest.Name = &pName
	}
	if d.HasChange("payment_method_id") {
		pPayment := d.Get("payment_method_id").(string)
		updateRequest.PaymentMethodID = &pPayment
	}
	if d.HasChange("backend_transfer") {
		pBT := d.Get("backend_transfer").(bool)
		updateRequest.BackendTransfer = &pBT
	}
	if d.HasChange("bgp_config") {
		o, n := d.GetChange("bgp_config")
		oldarr := o.([]interface{})
		newarr := n.([]interface{})
		if len(newarr) == 1 {
			bgpCreateRequest := expandBGPConfig(d)
			_, err := client.BGPConfig.Create(d.Id(), bgpCreateRequest)
			if err != nil {
				return friendlyError(err)
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

				errStr := fmt.Errorf("BGP Config can not be removed from a project, please add back\n%s", bgpConfStr)
				return friendlyError(errStr)
			}
		}
	} else {
		_, _, err := client.Projects.Update(d.Id(), updateRequest)
		if err != nil {
			return friendlyError(err)
		}
	}

	return resourceMetalProjectRead(d, meta)
}

func resourceMetalProjectDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*Config).addModuleToMetalUserAgent(d)
	client := meta.(*Config).metal

	resp, err := client.Projects.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
