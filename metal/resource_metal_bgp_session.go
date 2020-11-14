package metal

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

func resourceMetalBGPSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalBGPSessionCreate,
		Read:   resourceMetalBGPSessionRead,
		Delete: resourceMetalBGPSessionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"address_family": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ipv4", "ipv6"}, false),
			},
			"default_route": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMetalBGPSessionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	dID := d.Get("device_id").(string)
	addressFamily := d.Get("address_family").(string)
	defaultRoute := d.Get("default_route").(bool)
	log.Printf("[DEBUG] creating %s BGP session to device (%s)\n", addressFamily, dID)
	bgpSession, _, err := client.BGPSessions.Create(
		dID, packngo.CreateBGPSessionRequest{
			AddressFamily: addressFamily,
			DefaultRoute:  &defaultRoute})
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(bgpSession.ID)
	return resourceMetalBGPSessionRead(d, meta)
}

func resourceMetalBGPSessionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	bgpSession, _, err := client.BGPSessions.Get(d.Id(),
		&packngo.GetOptions{Includes: []string{"device"}})
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
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
	client := meta.(*packngo.Client)
	_, err := client.BGPSessions.Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}
