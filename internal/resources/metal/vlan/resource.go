package vlan

import (
	"context"
	"path"

	"github.com/equinix/terraform-provider-equinix/internal/converters"

	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

	"github.com/equinix/terraform-provider-equinix/internal/config"

	"github.com/equinix/equinix-sdk-go/services/metalv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceMetalVlanCreate,
		ReadWithoutTimeout:   resourceMetalVlanRead,
		DeleteWithoutTimeout: resourceMetalVlanDelete,
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

func resourceMetalVlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalGoUserAgent(d)
	client := meta.(*config.Config).Metalgo

	facRaw, facOk := d.GetOk("facility")
	metroRaw, metroOk := d.GetOk("metro")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")

	if !facOk && !metroOk {
		return diag.Errorf("one of facility or metro must be configured")
	}
	if facOk && vxlanOk {
		return diag.Errorf("you can set vxlan only for metro vlans")
	}

	createRequest := metalv1.VirtualNetworkCreateInput{
		Description: metalv1.PtrString(d.Get("description").(string)),
	}
	if metroOk {
		createRequest.Metro = metalv1.PtrString(metroRaw.(string))
		createRequest.Vxlan = metalv1.PtrInt32(int32(vxlanRaw.(int)))
	}
	if facOk {
		createRequest.Facility = facRaw.(string)
	}
	vlan, _, err := client.VLANsApi.CreateVirtualNetwork(ctx, d.Get("project_id").(string)).VirtualNetworkCreateInput(createRequest).Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}
	d.SetId(vlan.GetId())
	return resourceMetalVlanRead(ctx, d, meta)
}

func resourceMetalVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalGoUserAgent(d)
	client := meta.(*config.Config).Metalgo

	vlan, _, err := client.VLANsApi.GetVirtualNetwork(ctx, d.Id()).
		Include([]string{"assigned_to"}).Execute()
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		if equinix_errors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)

	}
	d.Set("description", vlan.Description)
	d.Set("project_id", vlan.AssignedTo.GetId()) // assigned_to is a project but specced as an href
	d.Set("vxlan", vlan.GetVxlan())
	d.Set("facility", vlan.FacilityCode) // vlan spec does not include facility_code; should we remove it?
	d.Set("metro", vlan.MetroCode)
	return nil
}

func resourceMetalVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalGoUserAgent(d)
	client := meta.(*config.Config).Metalgo

	id := d.Id()
	vlan, resp, err := client.VLANsApi.GetVirtualNetwork(ctx, id).Include([]string{"instances", "instances.network_ports.virtual_networks", "internet_gateway"}).Execute()
	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	} else if err != nil {
		// missing vlans are deleted
		return nil
	}

	// all device ports must be unassigned before delete
	for _, i := range vlan.Instances {
		for _, p := range i.NetworkPorts { // instances is specced as a list of href; should be devices?
			for _, a := range p.AttachedVirtualNetworks {
				// a.ID is not set despite including instaces.network_ports.virtual_networks
				// TODO(displague) packngo should offer GetID() that uses ID or Href
				aID := path.Base(a.Href)

				if aID == id {
					portInput := metalv1.PortAssignInput{
						Vnid: &id,
					}
					_, resp, err := client.PortsApi.UnassignPort(ctx, p.GetId()).PortAssignInput(portInput).Execute()

					if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
						return diag.FromErr(equinix_errors.FriendlyError(err))
					}
				}
			}
		}
	}

	// TODO(displague) do we need to unassign gateway connections before delete?
	_, resp, err = client.VLANsApi.DeleteVirtualNetwork(ctx, id).Execute()

	if equinix_errors.IgnoreHttpResponseErrors(equinix_errors.HttpForbidden, equinix_errors.HttpNotFound)(resp, err) != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}

	return nil
}
