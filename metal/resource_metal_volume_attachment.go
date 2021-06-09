package metal

import (
	"log"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create:             resourceMetalVolumeAttachmentCreate,
		Read:               resourceMetalVolumeAttachmentRead,
		DeprecationMessage: "Volumes are deprecated, see https://metal.equinix.com/developers/docs/resilience-recovery/elastic-block-storage/#elastic-block-storage",
		Delete:             resourceMetalVolumeAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceMetalVolumeAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	dID := d.Get("device_id").(string)
	vID := d.Get("volume_id").(string)
	log.Printf("[DEBUG] Attaching Volume (%s) to Instance (%s)\n", vID, dID)
	va, _, err := client.VolumeAttachments.Create(vID, dID)
	if err != nil {
		switch err.(type) {
		case *packngo.ErrorResponse:
			e := err.(*packngo.ErrorResponse)
			if len(e.Errors) == 1 {
				if e.Errors[0] == "Instance is already attached to this volume" {
					log.Printf("[DEBUG] Volume (%s) is already attached to Instance (%s)", vID, dID)
					break
				}
			}
		}
		return err
	}

	d.SetId(va.ID)
	return resourceMetalVolumeAttachmentRead(d, meta)
}

func resourceMetalVolumeAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	va, _, err := client.VolumeAttachments.Get(d.Id(), nil)
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			log.Printf("[WARN] Volume Attachment (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("device_id", path.Base(va.Device.Href))
	d.Set("volume_id", path.Base(va.Volume.Href))
	return nil
}

func resourceMetalVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	resp, err := client.VolumeAttachments.Delete(d.Id())
	return ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err)
}
