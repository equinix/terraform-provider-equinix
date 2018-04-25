package packet

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketOrganizationCreate,
		Read:   resourcePacketOrganizationRead,
		Update: resourcePacketOrganizationUpdate,
		Delete: resourcePacketOrganizationDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"website": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"twitter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"logo": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Required: false,
			},

			"created": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePacketOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	createRequest := &packngo.OrganizationCreateRequest{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("website"); ok {
		createRequest.Website = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		createRequest.Description = attr.(string)
	}

	if attr, ok := d.GetOk("twitter"); ok {
		createRequest.Twitter = attr.(string)
	}

	if attr, ok := d.GetOk("logo"); ok {
		createRequest.Logo = attr.(string)
	}

	org, _, err := client.Organizations.Create(createRequest)
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(org.ID)

	return resourcePacketOrganizationRead(d, meta)
}

func resourcePacketOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	key, _, err := client.Organizations.Get(d.Id())
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.Set("id", key.ID)
	d.Set("name", key.Name)
	d.Set("description", key.Description)
	d.Set("website", key.Website)
	d.Set("twitter", key.Twitter)
	d.Set("logo", key.Logo)
	d.Set("created", key.Created)
	d.Set("updated", key.Updated)

	return nil
}

func resourcePacketOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	website := d.Get("website").(string)
	twitter := d.Get("twitter").(string)
	logo := d.Get("logo").(string)

	updateRequest := &packngo.OrganizationUpdateRequest{
		Name:        &name,
		Description: &description,
		Website:     &website,
		Twitter:     &twitter,
		Logo:        &logo,
	}

	_, _, err := client.Organizations.Update(d.Get("id").(string), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourcePacketOrganizationRead(d, meta)
}

func resourcePacketOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	_, err := client.Organizations.Delete(d.Id())
	if err != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}
