package metal

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourceMetalHardwareReservation() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalHardwareReservationRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the hardware reservation to look up",
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
			"device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of device occupying the reservation",
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
				Description: "Flag indicating whether the reservation can be currently used to create a device",
			},
			"spare": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Flag indicating whether the reservation is spare (@displague help),",
			},
			"switch_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of switch (@displague help)",
			},
			"intervals": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "(@displague help)",
			},
			"current_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "(@displague help)",
			},
		},
	}
}

func dataSourceMetalHardwareReservationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	hrId := d.Get("id").(string)

	hr, _, err := client.HardwareReservations.Get(
		hrId,
		&packngo.GetOptions{Includes: []string{"project", "facility", "device"}})
	if err != nil {
		return err
	}
	deviceId := ""
	if hr.Device != nil {
		deviceId = hr.Device.ID
	}

	m := map[string]interface{}{
		"short_id":       hr.ShortID,
		"project_id":     hr.Project.ID,
		"device_id":      deviceId,
		"plan":           hr.Plan.Slug,
		"facility":       hr.Facility.Code,
		"provisionable":  hr.Provisionable,
		"spare":          hr.Spare,
		"switch_uuid":    hr.SwitchUUID,
		"intervals":      hr.Intervals,
		"current_period": hr.CurrentPeriod,
	}

	d.SetId(hr.ID)
	return setMap(d, m)
}
