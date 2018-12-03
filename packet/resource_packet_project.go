package packet

import (
	"path"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/packethost/packngo"
)

var uuidRE = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func resourcePacketProject() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketProjectCreate,
		Read:   resourcePacketProjectRead,
		Update: resourcePacketProjectUpdate,
		Delete: resourcePacketProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"payment_method_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(strings.Trim(old, `"`)) == strings.ToLower(strings.Trim(new, `"`))
				},
				ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
			},

			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(strings.Trim(old, `"`)) == strings.ToLower(strings.Trim(new, `"`))
				},
				ValidateFunc: validation.StringMatch(uuidRE, "must be a valid UUID"),
			},
		},
	}
}

func resourcePacketProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.ProjectCreateRequest{
		Name:           d.Get("name").(string),
		OrganizationID: d.Get("organization_id").(string),
	}

	project, _, err := client.Projects.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(project.ID)

	return resourcePacketProjectRead(d, meta)
}

func resourcePacketProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	proj, _, err := client.Projects.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.Set("id", proj.ID)
	d.Set("payment_method_id", path.Base(proj.PaymentMethod.URL))
	d.Set("name", proj.Name)
	d.Set("organization_id", path.Base(proj.Organization.URL))
	d.Set("created", proj.Created)
	d.Set("updated", proj.Updated)

	return nil
}

func resourcePacketProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	updateRequest := &packngo.ProjectUpdateRequest{}
	if d.HasChange("name") {
		pName := d.Get("name").(string)
		updateRequest.Name = &pName
	}
	if d.HasChange("payment_method_id") {
		pPayment := d.Get("payment_method_id").(string)
		updateRequest.PaymentMethodID = &pPayment
	}
	_, _, err := client.Projects.Update(d.Id(), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourcePacketProjectRead(d, meta)
}

func resourcePacketProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.Projects.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
