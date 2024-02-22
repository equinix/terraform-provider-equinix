package equinix

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/equinix/terraform-provider-equinix/internal/converters"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/packethost/packngo"
)

func resourceMetalVlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalVlanCreate,
		Read:   resourceMetalVlanRead,
		Delete: resourceMetalVlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Description: "ID of parent project",
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description string",
				Optional:    true,
				ForceNew:    true,
			},
			"facility": {
				Type:          schema.TypeString,
				Description:   "Facility where to create the VLAN",
				Deprecated:    "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"metro"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// suppress diff when unsetting facility
					if len(old) > 0 && new == "" {
						return true
					}
					return old == new
				},
			},
			"metro": {
				Type:          schema.TypeString,
				Description:   "Metro in which to create the VLAN",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"facility"},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, facOk := d.GetOk("facility")
					// new - new val from template
					// old - old val from state
					//
					// suppress diff if metro is manually set for first time, and
					// facility is already set
					if len(new) > 0 && old == "" && facOk {
						return facOk
					}
					return old == new
				},
				StateFunc: converters.ToLowerIf,
			},
			"vxlan": {
				Type:        schema.TypeInt,
				Description: "VLAN ID, must be unique in metro",
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceMetalVlanCreate(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	facRaw, facOk := d.GetOk("facility")
	metroRaw, metroOk := d.GetOk("metro")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")

	if !facOk && !metroOk {
		return equinix_errors.FriendlyError(errors.New("one of facility or metro must be configured"))
	}
	if facOk && vxlanOk {
		return equinix_errors.FriendlyError(errors.New("you can set vxlan only for metro vlans"))
	}

	createRequest := &packngo.VirtualNetworkCreateRequest{
		ProjectID:   d.Get("project_id").(string),
		Description: d.Get("description").(string),
	}
	if metroOk {
		createRequest.Metro = metroRaw.(string)
		createRequest.VXLAN = vxlanRaw.(int)
	}
	if facOk {
		createRequest.Facility = facRaw.(string)
	}
	vlan, _, err := client.ProjectVirtualNetworks.Create(createRequest)
	if err != nil {
		return equinix_errors.FriendlyError(err)
	}
	d.SetId(vlan.ID)
	return resourceMetalVlanRead(d, meta)
}

func resourceMetalVlanRead(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	vlan, _, err := client.ProjectVirtualNetworks.Get(d.Id(),
		&packngo.GetOptions{Includes: []string{"assigned_to"}})
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		if equinix_errors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return err

	}
	d.Set("description", vlan.Description)
	d.Set("project_id", vlan.Project.ID)
	d.Set("vxlan", vlan.VXLAN)
	d.Set("facility", vlan.FacilityCode)
	d.Set("metro", vlan.MetroCode)
	return nil
}

func resourceMetalVlanDelete(d *schema.ResourceData, meta interface{}) error {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metal

	id := d.Id()
	vlan, resp, err := client.ProjectVirtualNetworks.Get(id, &packngo.GetOptions{Includes: []string{"instances", "instances.network_ports.virtual_networks", "internet_gateway"}})
	if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return equinix_errors.FriendlyError(err)
	} else if err != nil {
		// missing vlans are deleted
		return nil
	}

	// all device ports must be unassigned before delete
	for _, i := range vlan.Instances {
		for _, p := range i.NetworkPorts {
			for _, a := range p.AttachedVirtualNetworks {
				// a.ID is not set despite including instaces.network_ports.virtual_networks
				// TODO(displague) packngo should offer GetID() that uses ID or Href
				aID := path.Base(a.Href)

				if aID == id {
					_, resp, err := client.Ports.Unassign(p.ID, id)

					if equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
						return equinix_errors.FriendlyError(err)
					}
				}
			}
		}
	}

	// TODO(displague) do we need to unassign gateway connections before delete?

	return equinix_errors.FriendlyError(equinix_errors.IgnoreResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(client.ProjectVirtualNetworks.Delete(id)))
}

func addMetalVlanSweeper() {
	resource.AddTestSweepers("equinix_metal_vlan", &resource.Sweeper{
		Name:         "equinix_metal_vlan",
		Dependencies: []string{"equinix_metal_virtual_circuit", "equinix_metal_vrf", "equinix_metal_device"},
		F:            testSweepVlans,
	})
}

func testSweepVlans(region string) error {
	log.Printf("[DEBUG] Sweeping vlans")
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting configuration for sweeping vlans: %s", err)
	}
	metal := config.NewMetalClient()
	ps, _, err := metal.Projects.List(nil)
	if err != nil {
		return fmt.Errorf("[INFO][SWEEPER_LOG] Error getting project list for sweeping vlans: %s", err)
	}
	pids := []string{}
	for _, p := range ps {
		if isSweepableTestResource(p.Name) {
			pids = append(pids, p.ID)
		}
	}
	dids := []string{}
	for _, pid := range pids {
		ds, _, err := metal.ProjectVirtualNetworks.List(pid, nil)
		if err != nil {
			log.Printf("Error listing vlans to sweep: %s", err)
			continue
		}
		for _, d := range ds.VirtualNetworks {
			if isSweepableTestResource(d.Description) {
				dids = append(dids, d.ID)
			}
		}
	}

	for _, did := range dids {
		log.Printf("Removing vlan %s", did)
		_, err := metal.ProjectVirtualNetworks.Delete(did)
		if err != nil {
			return fmt.Errorf("Error deleting vlan %s", err)
		}
	}
	return nil
}
