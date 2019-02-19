package packet

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/packethost/packngo"
)

func resourcePacketBGPSession() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketBGPSessionCreate,
		Read:   resourcePacketBGPSessionRead,
		Delete: resourcePacketBGPSessionDelete,
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

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePacketBGPSessionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	dID := d.Get("device_id").(string)
	addressFamily := d.Get("address_family").(string)
	log.Printf("[DEBUG] creating %s BGP session to device (%s)\n", addressFamily, dID)
	bgpSession, _, err := client.BGPSessions.Create(
		dID, packngo.CreateBGPSessionRequest{AddressFamily: "ipv4"})
	if err != nil {
		return friendlyError(err)
	}

	d.SetId(bgpSession.ID)
	return resourcePacketBGPSessionRead(d, meta)
}

func resourcePacketBGPSessionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	bgpSession, _, err := client.BGPSessions.Get(d.Id(),
		&packngo.GetOptions{Includes: []string{"device"}})
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("device_id", bgpSession.Device.ID)
	d.Set("address_family", bgpSession.AddressFamily)
	d.Set("status", bgpSession.Status)
	d.SetId(bgpSession.ID)
	return nil
}

func resourcePacketBGPSessionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	_, err := client.BGPSessions.Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}
