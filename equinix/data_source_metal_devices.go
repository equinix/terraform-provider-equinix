package equinix

import (
	"context"
	"fmt"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/equinix/terraform-provider-equinix/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMetalDevices() *schema.Resource {
	dsmd := dataSourceMetalDevice()
	sch := dsmd.Schema
	for _, v := range sch {
		if v.Optional {
			v.Optional = false
		}
		if v.ConflictsWith != nil {
			v.ConflictsWith = nil
		}
	}
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:               sch,
		ResultAttributeName:        "devices",
		ResultAttributeDescription: "List of devices that match specified filters",
		FlattenRecord:              flattenDevice,
		GetRecords:                 getDevices,
		ExtraQuerySchema: map[string]*schema.Schema{
			"project_id": {
				Type:          schema.TypeString,
				Description:   "The id of the project to query for devices",
				Optional:      true,
				ConflictsWith: []string{"organization_id"},
			},
			"organization_id": {
				Type:          schema.TypeString,
				Description:   "The id of the organization to query for devices",
				Optional:      true,
				ConflictsWith: []string{"project_id"},
			},
			"search": {
				Type:        schema.TypeString,
				Description: "Search string to filter devices by hostname, description, short_id, reservation short_id, tags, plan name, plan slug, facility code, facility name, operating system name, operating system slug, IP addresses.",
				Optional:    true,
			},
		},
	}
	return datalist.NewResource(dataListConfig)
}

func getDevices(ctx context.Context, d *schema.ResourceData, meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.Config).NewMetalClientForSDK(d)
	projectID := extra["project_id"].(string)
	orgID := extra["organization_id"].(string)

	if (len(projectID) == 0) && (len(orgID) == 0) {
		return nil, fmt.Errorf("one of project_id or organization_id must be specified")
	}

	search := extra["search"].(string)

	var devices *metalv1.DeviceList
	devicesIf := []interface{}{}
	var err error

	if len(projectID) > 0 {
		query := client.DevicesApi.FindProjectDevices(
			ctx, projectID).Include(deviceCommonIncludes)
		if len(search) > 0 {
			query = query.Search(search)
		}
		devices, err = query.ExecuteWithPagination()
	}

	if len(orgID) > 0 {
		query := client.DevicesApi.FindOrganizationDevices(
			ctx, orgID).Include(deviceCommonIncludes)
		if len(search) > 0 {
			query = query.Search(search)
		}
		devices, err = query.ExecuteWithPagination()
	}

	for _, d := range devices.Devices {
		devicesIf = append(devicesIf, d)
	}
	return devicesIf, err
}

func flattenDevice(rawDevice interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	device, ok := rawDevice.(metalv1.Device)
	if !ok {
		return nil, fmt.Errorf("expected device to be of type *metalv1.Device, got %T", rawDevice)
	}
	return getDeviceMap(device), nil
}
