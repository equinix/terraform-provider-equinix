package metal_organization

import (
	"context"
	"fmt"
	"regexp"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/packethost/packngo"
	"github.com/equinix/terraform-provider-equinix/internal/helper"
)

var countryRE = regexp.MustCompile("(?i)^[a-z]{2}$")

type OrganizationResourceModel struct {
    ID          types.String   `tfsdk:"id"`
    Name        types.String   `tfsdk:"name"`
    Description types.String   `tfsdk:"description"`
    Website     types.String   `tfsdk:"website"`
    Twitter     types.String   `tfsdk:"twitter"`
    Logo        types.String   `tfsdk:"logo"`
    Created     types.String   `tfsdk:"created"`
    Updated     types.String   `tfsdk:"updated"`
    Address     OrganizationAddress `tfsdk:"address"`
}

type OrganizationAddress struct {
    Address  types.String `tfsdk:"address"`
    City     types.String `tfsdk:"city"`
    ZipCode  types.String `tfsdk:"zip_code"`
    Country  types.String `tfsdk:"country"`
    State    types.String `tfsdk:"state"`
}

func (rm *OrganizationResourceModel) parse(org *packngo.Organization) diag.Diagnostics {
    var diags diag.Diagnostics

    rm.ID = types.StringValue(org.ID)
    rm.Name = types.StringValue(org.Name)
    rm.Description = types.StringValue(org.Description)
    rm.Website = types.StringValue(org.Website)
    rm.Twitter = types.StringValue(org.Twitter)
    rm.Logo = types.StringValue(org.Logo)
    rm.Created = types.StringValue(org.Created)
    rm.Updated = types.StringValue(org.Updated)

	address := OrganizationAddress{}
	diags.Append(address.parse(&org.Address)...)
	rm.Address = address

    return diags
}

func (addr *OrganizationAddress) parse(a *packngo.Address) diag.Diagnostics {
    var diags diag.Diagnostics

    addr.Address = types.StringValue(a.Address)
    if a.City != nil {
        addr.City = types.StringValue(*a.City)
    }
    addr.ZipCode = types.StringValue(a.ZipCode)
    addr.Country = types.StringValue(a.Country)
    if a.State != nil {
        addr.State = types.StringValue(*a.State)
    }

    return diags
}


func NewResource() resource.Resource {
	return &Resource{
		BaseResource: helper.NewBaseResource(
			helper.BaseResourceConfig{
				Name:   "equinix_organization",
				Schema: &organizationResourceSchema,
			},
		),
	}
}

type Resource struct {
	helper.BaseResource
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Create an instance of your resource model to hold the planned state
    var plan OrganizationResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Required data for the API request
    createRequest := packngo.OrganizationCreateRequest{
        Name: plan.Name.ValueString(),
        // Expand the address
        Address: packngo.Address{
            Address:  plan.Address.Address.ValueString(),
            City:     plan.Address.City.ValueStringPointer(),
            ZipCode:  plan.Address.ZipCode.ValueString(),
            Country:  plan.Address.Country.ValueString(),
            State:    plan.Address.State.ValueStringPointer(),
        },
    }

    // Optional fields
    if !plan.Description.IsNull() {
        createRequest.Description = plan.Description.ValueString()
    }
    if !plan.Website.IsNull() {
        createRequest.Website = plan.Website.ValueString()
    }
    if !plan.Twitter.IsNull() {
        createRequest.Twitter = plan.Twitter.ValueString()
    }
    if !plan.Logo.IsNull() {
        createRequest.Logo = plan.Logo.ValueString()
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // API call to create the organization
    org, _, err := client.Organizations.Create(&createRequest)
    if err != nil {
        resp.Diagnostics.AddError(
            "Error creating Organization",
            "Could not create Organization: " + err.Error(),
        )
        return
    }

    // Parse API response into the Terraform state
    stateDiags := plan.parse(org)
    resp.Diagnostics.Append(stateDiags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Set the state
    diags = resp.State.Set(ctx, &plan)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Retrieve the current state
    var state OrganizationResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the resource from the state
    id := state.ID.ValueString()

    // API call to get the current state of the organization
    org, _, err := client.Organizations.Get(id, &packngo.GetOptions{Includes: []string{"address"}})
    if err != nil {
        err = helper.FriendlyError(err)

        // Check if the organization no longer exists
        if helper.IsNotFound(err) {
			resp.Diagnostics.AddWarning(
				"Organization",
				fmt.Sprintf("[WARN] Organization (%s) not found, removing from state", id),
			)
            resp.State.RemoveResource(ctx)
            return
        }

        resp.Diagnostics.AddError(
            "Error reading Organization",
            "Could not read Organization with ID " + id + ": " + err.Error(),
        )
        return
    }

    // Parse the API response into the Terraform state
    stateDiags := state.parse(org)
    resp.Diagnostics.Append(stateDiags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Update the Terraform state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Retrieve the current state and plan
    var state, plan OrganizationResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    diags = req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client from the provider metadata
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the organization from the state
    id := state.ID.ValueString()

    // Prepare update request based on the changes
    updateRequest := &packngo.OrganizationUpdateRequest{}
    if state.Name != plan.Name {
        updateRequest.Name = plan.Name.ValueStringPointer()
    }
    if state.Description != plan.Description {
        updateRequest.Description = plan.Description.ValueStringPointer()
    }
    if state.Website != plan.Website {
        updateRequest.Website = plan.Website.ValueStringPointer()
    }
    if state.Twitter != plan.Twitter {
        updateRequest.Twitter = plan.Twitter.ValueStringPointer()
    }
    if state.Logo != plan.Logo {
        updateRequest.Logo = plan.Logo.ValueStringPointer()
    }

    // Handle address updates
    if !reflect.DeepEqual(state.Address, plan.Address) {
        updateRequest.Address = &packngo.Address{
            Address:  plan.Address.Address.ValueString(),
            City:     plan.Address.City.ValueStringPointer(),
            ZipCode:  plan.Address.ZipCode.ValueString(),
            Country:  plan.Address.Country.ValueString(),
            State:    plan.Address.State.ValueStringPointer(),
        }
    }

    // API call to update the organization
    updatedOrg, _, err := client.Organizations.Update(id, updateRequest)
    if err != nil {
		err = helper.FriendlyError(err)
        resp.Diagnostics.AddError(
            "Error updating Organization",
            "Could not update Organization with ID " + id + ": " + err.Error(),
        )
        return
    }

   // Parse the updated API response into the Terraform state
   stateDiags := state.parse(updatedOrg)
   resp.Diagnostics.Append(stateDiags...)
   if stateDiags.HasError() {
	   return
   }

   // Update the Terraform state
   diags = resp.State.Set(ctx, &state)
   resp.Diagnostics.Append(diags...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve the current state
    var state OrganizationResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Retrieve the API client
    r.Meta.AddFwModuleToMetalUserAgent(ctx, req.ProviderMeta)
    client := r.Meta.Metal

    // Extract the ID of the organization
    id := state.ID.ValueString()

    // API call to delete the organization
    deleteResp, err := client.Organizations.Delete(id)
    if helper.IgnoreResponseErrors(helper.HttpForbidden, helper.HttpNotFound)(deleteResp, err) != nil {
		err = helper.FriendlyError(err)
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failed to delete Organization %s", id),
			err.Error(),
		)
	}
}

var organizationResourceSchema = schema.Schema{
    Attributes: map[string]schema.Attribute{
        "id": schema.StringAttribute{
            Description: "The unique identifier for the Organization",
            Computed:    true,
        },
        "name": schema.StringAttribute{
            Description: "The name of the Organization",
            Required:    true,
        },
        "description": schema.StringAttribute{
            Description: "Description string",
            Optional:    true,
        },
        "website": schema.StringAttribute{
            Description: "Website link",
            Optional:    true,
        },
        "twitter": schema.StringAttribute{
            Description: "Twitter handle",
            Optional:    true,
        },
        "logo": schema.StringAttribute{
            Description: "Logo URL",
            Optional:    true,
        },
        "created": schema.StringAttribute{
            Description: "The timestamp for when the Organization was created",
            Computed:    true,
        },
        "updated": schema.StringAttribute{
            Description: "The timestamp for the last time the Organization was updated",
            Computed:    true,
        },
        "address": schema.SingleNestedAttribute{
            Description: "Address information block",
            Required:    true,
            Attributes: addressSchema,
        },
    },
}

var addressSchema = map[string]schema.Attribute{
	"address": schema.StringAttribute{
		Description:  "Postal address",
		Required:     true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	},
	"city": schema.StringAttribute{
		Description:  "City name",
		Required:     true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	},
	"zip_code": schema.StringAttribute{
		Description:  "Zip Code",
		Required:     true,
		Validators: []validator.String{
			stringvalidator.LengthAtLeast(1),
		},
	},
	"country": schema.StringAttribute{
		Description:  "Two letter country code (ISO 3166-1 alpha-2), e.g. US",
		Required:     true,
		Validators: []validator.String{
			stringvalidator.RegexMatches(countryRE, "Address country must be a two letter code (ISO 3166-1 alpha-2)"),
		},
	},
	"state": schema.StringAttribute{
		Description: "State name",
		Optional:    true,
	},
}
