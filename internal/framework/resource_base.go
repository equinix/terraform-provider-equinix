package framework

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/equinix/terraform-provider-equinix/internal/config"
)

func GetResourceMeta(
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) *config.Config {
	meta, ok := req.ProviderData.(*config.Config)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *http.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)
		return nil
	}

	return meta
}

// NewBaseResource returns a new instance of the BaseResource
// struct for cleaner initialization.
func NewBaseResource(cfg BaseResourceConfig) BaseResource {
	return BaseResource{
		Config: cfg,
	}
}

// BaseResourceConfig contains all configurable base resource fields.
type BaseResourceConfig struct {
	Name   string
	IDAttr string

	// Optional
	Schema *schema.Schema
}

// BaseResource contains various re-usable fields and methods
// intended for use in resource implementations by composition.
type BaseResource struct {
	Config BaseResourceConfig
	Meta   *config.Config
}

func (r *BaseResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.Meta = GetResourceMeta(req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *BaseResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = r.Config.Name
}

func (r *BaseResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	if r.Config.Schema == nil {
		resp.Diagnostics.AddError(
			"Missing Schema",
			"Base resource was not provided a schema. "+
				"Please provide a Schema config attribute or implement, the Schema(...) function.",
		)
		return
	}

	resp.Schema = *r.Config.Schema
}

// ImportState should be overridden for resources with
// complex read logic (e.g. parent ID).
func (r *BaseResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Enforce defaults
	idAttr := r.Config.IDAttr
	if idAttr == "" {
		idAttr = "id"
	}

	attrPath := path.Root(idAttr)

	if attrPath.Equal(path.Empty()) {
		resp.Diagnostics.AddError(
			"Resource Import Passthrough Missing Attribute Path",
			"This is always an error in the provider. Please report the following to the provider developer:\n\n"+
				"Resource ImportState path must be set to a valid attribute path.",
		)
		return
	}

	// Handle type conversion
	var err error
	var idValue any

	idValue = req.ID

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert ID attribute",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, attrPath, idValue)...)
}