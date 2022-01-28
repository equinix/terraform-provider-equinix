package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func resourceMetalVolume() *schema.Resource {
	return &schema.Resource{
		Create:             removedResourceOp(resourceVolumeRemovedMsg),
		Read:               removedResourceOp(resourceVolumeRemovedMsg),
		DeprecationMessage: "Volumes are deprecated, see https://metal.equinix.com/developers/docs/resilience-recovery/elastic-block-storage/#elastic-block-storage",
		Update:             removedResourceOp(resourceVolumeRemovedMsg),
		Delete:             resourceMetalVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"facility": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"plan": {
				Type:     schema.TypeString,
				Required: true,
			},

			"billing_cycle": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"locked": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"snapshot_policies": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_frequency": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"snapshot_count": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"attachments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceMetalVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	resp, err := client.Volumes.Delete(d.Id())
	if ignoreResponseErrors(httpForbidden, httpNotFound)(resp, err) != nil {
		return friendlyError(err)
	}

	return nil
}
