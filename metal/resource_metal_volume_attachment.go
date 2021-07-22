package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		Create:             removedResourceOp(resourceVolumeAttachmentRemovedMsg),
		Read:               removedResourceOp(resourceVolumeAttachmentRemovedMsg),
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

func resourceMetalVolumeAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	resp, err := client.VolumeAttachments.Delete(d.Id())
	return ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err)
}
