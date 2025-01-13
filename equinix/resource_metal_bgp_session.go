package equinix

import (
	"log"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/packethost/packngo"
)

func resourceMetalBGPSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalBGPSessionCreate,
		Read:   resourceMetalBGPSessionRead,
		Delete: resourceMetalBGPSessionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:        schema.TypeString,
				Description: "ID of device",
				Required:    true,
				ForceNew:    true,
			},
			"address_family": {
				Type:         schema.TypeString,
				Description:  "ipv4 or ipv6",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"default_route": {
				Type:        schema.TypeBool,
				Description: "Boolean flag to set the default route policy. False by default",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},

			"status": {
				Type:        schema.TypeString,
				Description: "Status of the session - up or down",
				Computed:    true,
			},
		},
	}
}

func resourceMetalBGPSessionCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	dID := d.Get("device_id").(string)
	addressFamily := d.Get("address_family").(string)
	defaultRoute := d.Get("default_route").(bool)
	log.Printf("[DEBUG] creating %s BGP session to device (%s)\n", addressFamily, dID)
	bgpSession, _, err := client.BGPSessions.Create(
		dID, packngo.CreateBGPSessionRequest{
			AddressFamily: addressFamily,
			DefaultRoute:  &defaultRoute,
		})
	if err != nil {
		return equinix_errors.FriendlyError(err)
	}

	d.SetId(bgpSession.ID)
	return resourceMetalBGPSessionRead(d, meta)
}

func resourceMetalBGPSessionRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	bgpSession, _, err := client.BGPSessions.Get(d.Id(),
		&packngo.GetOptions{Includes: []string{"device"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		if equinix_errors.IsNotFound(err) {
			log.Printf("[WARN] BGP Session (%s) not found, removing from state", d.Id())

			d.SetId("")
			return nil
		}
		return err
	}
	defaultRoute := false
	if bgpSession.DefaultRoute != nil {
		if *(bgpSession.DefaultRoute) {
			defaultRoute = true
		}
	}
	d.Set("device_id", bgpSession.Device.ID)
	d.Set("address_family", bgpSession.AddressFamily)
	d.Set("status", bgpSession.Status)
	d.Set("default_route", defaultRoute)
	d.SetId(bgpSession.ID)
	return nil
}

func resourceMetalBGPSessionDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal
	resp, err := client.BGPSessions.Delete(d.Id())
	return equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err)
}
