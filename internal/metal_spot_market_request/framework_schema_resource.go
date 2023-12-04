package metal_spot_market_request

import (
    "context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
    "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
    "github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
)

func metalSpotMarketRequestResourceSchema(ctx context.Context) *schema.Schema {
    return &schema.Schema{
        Attributes: map[string]schema.Attribute{
            "timeouts": timeouts.Attributes(ctx, timeouts.Opts{
                Delete: true,
                Create: true,
            }),
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "The unique identifier of the reserved IP block",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "devices_min": schema.Int64Attribute{
                Required:    true,
                Description: "Minimum number devices to be created",
                PlanModifiers: []planmodifier.Int64{
                    int64planmodifier.RequiresReplace(),
                },
            },
            "devices_max": schema.Int64Attribute{
                Required:    true,
                Description: "Maximum number devices to be created",
                PlanModifiers: []planmodifier.Int64{
                    int64planmodifier.RequiresReplace(),
                },
            },
            "max_bid_price": schema.Float64Attribute{
                Required:    true,
                Description: "Maximum price user is willing to pay per hour per device",
                PlanModifiers: []planmodifier.Float64{
                    float64planmodifier.RequiresReplace(),
                },
                //TODO (ocobles)
                // DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// 	oldF, err := strconv.ParseFloat(old, 64)
				// 	if err != nil {
				// 		return false
				// 	}
				// 	newF, err := strconv.ParseFloat(new, 64)
				// 	if err != nil {
				// 		return false
				// 	}
				// 	// suppress diff if the difference between existing and new bid price
				// 	// is less than 2%
				// 	diffThreshold := .02
				// 	priceDiff := oldF / newF

				// 	if diffThreshold < priceDiff {
				// 		return true
				// 	}
				// 	return false
				// },
            },
            "facilities": schema.ListAttribute{
                Optional:    true,
                Computed:    true,
                ElementType: types.StringType,
                Description: "Facility IDs where devices should be created",
                DeprecationMessage: "Use metro instead of facility.  For more information, read the migration guide: https://registry.terraform.io/providers/equinix/equinix/latest/docs/guides/migration_guide_facilities_to_metros_devices",
                PlanModifiers: []planmodifier.List{
                    listplanmodifier.RequiresReplace(),
                },
                Validators: []validator.List{
                    listvalidator.ExactlyOneOf(path.Expressions{
                        path.MatchRoot("metro"),
                    }...),
                },
                //TODO (ocobles)
                // DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// 	oldData, newData := d.GetChange("facilities")

				// 	// If this function is called and oldData or newData is nil,
				// 	// then the attribute changed
				// 	if oldData == nil || newData == nil {
				// 		return false
				// 	}

				// 	oldArray := oldData.([]interface{})
				// 	newArray := newData.([]interface{})

				// 	// If the number of items in the list is different,
				// 	// then the attribute changed
				// 	if len(oldArray) != len(newArray) {
				// 		return false
				// 	}

				// 	// Convert data to string arrays
				// 	oldFacilities := make([]string, len(oldArray))
				// 	newFacilities := make([]string, len(newArray))
				// 	for i, oldFacility := range oldArray {
				// 		oldFacilities[i] = fmt.Sprint(oldFacility)
				// 	}
				// 	for j, newFacility := range newArray {
				// 		newFacilities[j] = fmt.Sprint(newFacility)
				// 	}
				// 	// Sort the old and new arrays so that we don't show a diff
				// 	// if the facilities are the same but the order is different
				// 	sort.Strings(oldFacilities)
				// 	sort.Strings(newFacilities)
				// 	return reflect.DeepEqual(oldFacilities, newFacilities)
            },
            "metro": schema.StringAttribute{
                Optional:    true,
                Description: "Metro where devices should be created",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
                //TODO (ocobles)
                // StateFunc:     toLower,
            },
            "project_id": schema.StringAttribute{
                Required:    true,
                Description: "Project ID",
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(),
                },
            },
            "wait_for_devices": schema.BoolAttribute{
                Optional:    true,
                Description: "On resource creation - wait until all desired devices are active, on resource destruction - wait until devices are removed",
                PlanModifiers: []planmodifier.Bool{
                    boolplanmodifier.RequiresReplace(),
                },
            },
            "instance_parameters": schema.SingleNestedAttribute{
                Required:    true,
                Description: "Parameters for devices provisioned from this request. You can find the parameter description from the [equinix_metal_device doc](device.md)",
                Attributes:  instanceParametersSchema, // Referencing the defined schema
                PlanModifiers: []planmodifier.Object{
                    objectplanmodifier.RequiresReplace(),
                },
            },
        },
    }
}

var instanceParametersSchema = map[string]schema.Attribute{
    "billing_cycle": schema.StringAttribute{
        Required:    true,
        Description: "Billing cycle for the instance",
    },
    "plan": schema.StringAttribute{
        Required:    true,
        Description: "The plan or size of the instance",
    },
    "operating_system": schema.StringAttribute{
        Required:    true,
        Description: "The operating system of the instance",
    },
    "hostname": schema.StringAttribute{
        Required:    true,
        Description: "The hostname of the instance",
    },
    "termintation_time": schema.StringAttribute{
        Computed:    true,
        Description: "The termination time of the instance",
        DeprecationMessage: "Use instance_parameters.termination_time instead",
    },
    "termination_time": schema.StringAttribute{
        Computed:    true,
        Description: "The termination time of the instance",
    },
    "always_pxe": schema.BoolAttribute{
        Optional: true,
        Default:  booldefault.StaticBool(false),
    },
    "description": schema.StringAttribute{
        Optional: true,
    },
    "features": schema.ListAttribute{
        Optional:    true,
        ElementType:  types.StringType,
    },
    "locked": schema.StringAttribute{
        Optional: true,
    },
    "project_ssh_keys": schema.ListAttribute{
        Optional:    true,
        ElementType:  types.StringType,
    },
    "user_ssh_keys": schema.ListAttribute{
        Optional:    true,
        ElementType:  types.StringType,
    },
    "userdata": schema.StringAttribute{
        Optional: true,
    },
    "customdata": schema.StringAttribute{
        Optional: true,
    },
    "ipxe_script_url": schema.StringAttribute{
        Optional: true,
    },
    "tags": schema.ListAttribute{
        Optional:    true,
        ElementType:  types.StringType,
    },
}

