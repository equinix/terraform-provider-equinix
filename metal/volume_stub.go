package metal

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const volumeRemovedMsgSuffix = "was removed in version 3.0.0, see https://metal.equinix.com/developers/docs/storage/elastic-block-storage/#elastic-block-storage"

var (
	resourceVolumeRemovedMsg           = fmt.Sprintf("resource metal_volume %s", volumeRemovedMsgSuffix)
	dataSourceVolumeRemovedMsg         = fmt.Sprintf("datasource metal_volume %s", volumeRemovedMsgSuffix)
	resourceVolumeAttachmentRemovedMsg = fmt.Sprintf("resource metal_volume_attachment %s", volumeRemovedMsgSuffix)
)

func removedResourceOp(message string) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		return errors.New(message)
	}
}
