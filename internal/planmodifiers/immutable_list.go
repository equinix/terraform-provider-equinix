package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func ImmutableList() planmodifier.List {
	return &immutableListPlanModifier{}
}

type immutableListPlanModifier struct{}

func (d *immutableListPlanModifier) PlanModifyList(ctx context.Context, request planmodifier.ListRequest, response *planmodifier.ListResponse) {

	if request.StateValue.IsNull() && request.PlanValue.IsNull() {
		return
	}

	if request.PlanValue.IsNull() && len(request.StateValue.Elements()) > 0 {
		response.Diagnostics.AddAttributeError(
			request.Path,
			"Change not allowed",
			fmt.Sprintf(
				"Elements of the `%s` list field can not be removed. Resource recreation would be required.",
				request.Path.String(),
			),
		)
		return
	}

	response.PlanValue = request.PlanValue
}

func (d *immutableListPlanModifier) Description(ctx context.Context) string {
	return "Allows adding elements to a list if it was initially empty and permits modifications, but disallows removals, requiring resource recreation."
}

func (d *immutableListPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}
