package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ImmutableInt64() planmodifier.Int64 {
	return &immutableInt64PlanModifier{}
}

type immutableInt64PlanModifier struct{}

func (d *immutableInt64PlanModifier) PlanModifyInt64(ctx context.Context, request planmodifier.Int64Request, response *planmodifier.Int64Response) {
	if request.StateValue.IsNull() && request.PlanValue.IsUnknown() {
		return
	}

	oldValue := request.StateValue.ValueInt64()
	newValue := request.PlanValue.ValueInt64()

	if oldValue != 0 && newValue != oldValue {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Change not allowed",
			fmt.Sprintf(
				"Cannot modify the value of the `%s` field. Resource recreation would be required.",
				request.Path.String(),
			),
		)
		return
	}

	response.PlanValue = types.Int64Value(newValue)
}

func (d *immutableInt64PlanModifier) Description(ctx context.Context) string {
	return "Prevents modification of a int64 value if the old value is not null."
}

func (d *immutableInt64PlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
