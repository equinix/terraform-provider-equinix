package metal

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalHardwareReservation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalHardwareReservationRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the hardware reservation to look up",
			},
			"device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "UUID of device occupying the reservation",
			},
			"short_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Reservation short ID",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of project this reservation is scoped to",
			},
			"plan": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan type for the reservation",
			},
			"facility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan type for the reservation",
			},
			"provisionable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the reserved server is provisionable or not. Spare devices can't be provisioned unless they are activated first",
			},
			"spare": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the Hardware Reservation is a spare. Spare Hardware Reservations are used when a Hardware Reservations requires service from Metal Equinix",
			},
			"switch_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Switch short ID, can be used to determine if two devices are connected to the same switch",
			},
		},
	}
}

func dataSourceMetalHardwareReservationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	hrIdRaw, hrIdOk := d.GetOk("id")
	dIdRaw, dIdOk := d.GetOk("device_id")

	if dIdOk == hrIdOk {
		return fmt.Errorf("You must set one of id and device_id")
	}

	var deviceId string
	var hr *packngo.HardwareReservation

	if dIdOk {
		deviceId = dIdRaw.(string)
		includes := []string{
			"hardware_reservation.project",
			"hardware_reservation.facility",
		}
		d, _, err := client.Devices.Get(deviceId, &packngo.GetOptions{Includes: includes})
		if err != nil {
			return err
		}
		if d.HardwareReservation == nil {
			return fmt.Errorf("Device %s is not in a hardware reservation", deviceId)
		}
		hr = d.HardwareReservation

	} else {
		var err error
		hr, _, err = client.HardwareReservations.Get(
			hrIdRaw.(string),
			&packngo.GetOptions{Includes: []string{"project", "facility", "device"}})
		if err != nil {
			return err
		}
		if hr.Device != nil {
			deviceId = hr.Device.ID
		}
	}

	m := map[string]interface{}{
		"short_id":      hr.ShortID,
		"project_id":    hr.Project.ID,
		"device_id":     deviceId,
		"plan":          hr.Plan.Slug,
		"facility":      hr.Facility.Code,
		"provisionable": hr.Provisionable,
		"spare":         hr.Spare,
		"switch_uuid":   hr.SwitchUUID,
	}

	d.SetId(hr.ID)
	return setMap(d, m)
}
