package equinix

import (
	"regexp"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func resourceMetalOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalOrganizationCreate,
		Read:   resourceMetalOrganizationRead,
		Update: resourceMetalOrganizationUpdate,
		Delete: resourceMetalOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the Organization",
				Required:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description string",
				Optional:    true,
			},
			"website": {
				Type:        schema.TypeString,
				Description: "Website link",
				Optional:    true,
			},
			"twitter": {
				Type:        schema.TypeString,
				Description: "Twitter handle",
				Optional:    true,
			},
			"logo": {
				Type:        schema.TypeString,
				Description: "Logo URL",
				Optional:    true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:        schema.TypeList,
				Description: "Address information block",
				Required:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: createMetalOrganizationAddressResourceSchema(),
				},
			},
		},
	}
}

func createMetalOrganizationAddressResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:         schema.TypeString,
			Description:  "Postal address",
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"city": {
			Type:         schema.TypeString,
			Description:  "City name",
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"zip_code": {
			Type:         schema.TypeString,
			Description:  "Zip Code",
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"country": {
			Type:         schema.TypeString,
			Description:  "Two letter country code (ISO 3166-1 alpha-2), e.g. US",
			Required:     true,
			ValidateFunc: validation.StringMatch(regexp.MustCompile("(?i)^[a-z]{2}$"), "Address country must be a two letter code (ISO 3166-1 alpha-2)"),
		},
		"state": {
			Type:        schema.TypeString,
			Description: "State name",
			Optional:    true,
		},
	}
}

func resourceMetalOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	createRequest := &packngo.OrganizationCreateRequest{
		Name:    d.Get("name").(string),
		Address: expandMetalOrganizationAddress(d.Get("address").([]interface{})),
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

	return resourceMetalOrganizationRead(d, meta)
}

func resourceMetalOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	key, _, err := client.Organizations.Get(d.Id(), &packngo.GetOptions{Includes: []string{"address"}})
	if err != nil {
		err = friendlyError(err)

		// If the project somehow already destroyed, mark as succesfully gone.
		if isNotFound(err) {
			d.SetId("")

			return nil
		}

		return err
	}

	d.SetId(key.ID)
	return setMap(d, map[string]interface{}{
		"name":        key.Name,
		"description": key.Description,
		"website":     key.Website,
		"twitter":     key.Twitter,
		"logo":        key.Logo,
		"created":     key.Created,
		"updated":     key.Updated,
		"address":     flattenMetalOrganizationAddress(key.Address),
	})
}

func resourceMetalOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	changes := getResourceDataChangedKeys([]string{"name", "description", "website", "twitter", "logo", "address"}, d)
	updateRequest := &packngo.OrganizationUpdateRequest{}
	for change, changeValue := range changes {
		switch change {
		case "name":
			cv := changeValue.(string)
			updateRequest.Name = &cv
		case "description":
			cv := changeValue.(string)
			updateRequest.Description = &cv
		case "website":
			cv := changeValue.(string)
			updateRequest.Website = &cv
		case "twitter":
			cv := changeValue.(string)
			updateRequest.Twitter = &cv
		case "logo":
			cv := changeValue.(string)
			updateRequest.Logo = &cv
		case "address":
			cv := expandMetalOrganizationAddress(changeValue.([]interface{}))
			updateRequest.Address = &cv
		}
	}

	_, _, err := client.Organizations.Update(d.Id(), updateRequest)
	if err != nil {
		return friendlyError(err)
	}

	return resourceMetalOrganizationRead(d, meta)
}

func resourceMetalOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	resp, err := client.Organizations.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	d.SetId("")
	return nil
}

func flattenMetalOrganizationAddress(addr packngo.Address) interface{} {
	result := make(map[string]interface{})
	if addr.Address != "" {
		result["address"] = addr.Address
	}
	if addr.City != nil && *addr.City != "" {
		result["city"] = addr.City
	}
	if addr.Country != "" {
		result["country"] = addr.Country
	}
	if addr.State != nil && *addr.State != "" {
		result["state"] = addr.State
	}
	if addr.ZipCode != "" {
		result["zip_code"] = addr.ZipCode
	}

	return []interface{}{result}
}

func expandMetalOrganizationAddress(address []interface{}) packngo.Address {
	transformed := packngo.Address{}
	addr := address[0].(map[string]interface{})

	if v, ok := addr["address"]; ok {
		transformed.Address = v.(string)
	}
	if v, ok := addr["city"]; ok {
		city := v.(string)
		transformed.City = &city
	}
	if v, ok := addr["zip_code"]; ok {
		transformed.ZipCode = v.(string)
	}
	if v, ok := addr["country"]; ok {
		transformed.Country = v.(string)
	}
	if v, ok := addr["state"]; ok {
		state := v.(string)
		transformed.State = &state
	}

	return transformed
}
