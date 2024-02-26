package vlans

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/equinix/terraform-provider-equinix/internal/config"
	"github.com/equinix/terraform-provider-equinix/internal/converters"
	equinix_errors "github.com/equinix/terraform-provider-equinix/internal/errors"

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
			StateContext: schema.ImportStatePassthroughContext,
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
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metalgo

	facRaw, facOk := d.GetOk("facility")
	metroRaw, metroOk := d.GetOk("metro")
	vxlanRaw, vxlanOk := d.GetOk("vxlan")

	if !facOk && !metroOk {
		return diag.FromErr(equinix_errors.FriendlyError(errors.New("one of facility or metro must be configured")))
	}
	if facOk && vxlanOk {
		return diag.FromErr(equinix_errors.FriendlyError(errors.New("you can set vxlan only for metro vlans")))
	}

	createRequest := metalv1.VirtualNetworkCreateInput{
		Description: metalv1.PtrString(d.Get("description").(string)),
	}
	if metroOk {
		createRequest.Metro = metalv1.PtrString(metroRaw.(string))
		createRequest.Vxlan = metalv1.PtrInt32(int32(vxlanRaw.(int)))
	}
	if facOk {
		createRequest.Facility = metalv1.PtrString(facRaw.(string))
	}
	projectId := d.Get("project_id").(string)
	vlan, _, err := client.VLANsApi.
		CreateVirtualNetwork(context.Background(), projectId).
		VirtualNetworkCreateInput(createRequest).
		Execute()
	if err != nil {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}
	d.SetId(vlan.GetId())
	return resourceMetalVlanRead(ctx, d, meta)
}

func resourceMetalVlanRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metalgo

	vlan, _, err := client.VLANsApi.
		GetVirtualNetwork(context.Background(), d.Id()).
		Include([]string{"assigned_to"}).
		Execute()
	if err != nil {
		err = equinix_errors.FriendlyError(err)
		if equinix_errors.IsNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)

	}
	d.Set("description", vlan.GetDescription())
	//d.Set("project_id", vlan)
	d.Set("vxlan", vlan.Vxlan)
	d.Set("facility", vlan.GetFacility())
	d.Set("metro", vlan.GetMetro())
	return nil
}

func resourceMetalVlanDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	meta.(*config.Config).AddModuleToMetalUserAgent(d)
	client := meta.(*config.Config).Metalgo

	vlan, resp, err := client.VLANsApi.
		GetVirtualNetwork(ctx, d.Id()).
		Include([]string{"instances", "meta_gateway"}).
		Execute()
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusNotFound {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	} else if err != nil {
		// missing vlans are deleted
		return nil
	}

	// all device ports must be unassigned before delete
	for _, i := range vlan.GetInstances() {
		devideId := path.Base(i.GetHref())
		device, resp, _ := client.DevicesApi.FindDeviceById(ctx, devideId).Execute()
		if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusNotFound {
			break
		}
		for _, p := range device.GetNetworkPorts() {
			for _, vlanHref := range p.GetVirtualNetworks() {
				vlanId := path.Base(vlanHref.GetHref())

				if vlanId == vlan.GetId() {
					_, resp, err := client.PortsApi.
						UnassignPort(ctx, p.GetId()).
						PortAssignInput(metalv1.PortAssignInput{Vnid: &vlanId}).
						Execute()
					if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusNotFound {
						return diag.FromErr(equinix_errors.FriendlyError(err))
					}
				}
			}
		}
	}

	_, resp, err = client.VLANsApi.DeleteVirtualNetwork(ctx, vlan.GetId()).Execute()
	if resp.StatusCode != http.StatusForbidden && resp.StatusCode != http.StatusNotFound {
		return diag.FromErr(equinix_errors.FriendlyError(err))
	}

	return nil
}
