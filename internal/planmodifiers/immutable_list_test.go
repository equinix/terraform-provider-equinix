package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestImmutableListSet(t *testing.T) {
	testCases := []struct {
		Old, New, Expected []string
		ExpectError        bool
	}{
		{
			Old:         []string{},
			New:         []string{"test"},
			Expected:    []string{"test"},
			ExpectError: false,
		},
		{
			Old:         []string{"test"},
			New:         []string{},
			Expected:    []string{},
			ExpectError: true,
		},
		{
			Old:         []string{"foo"},
			New:         []string{"bar"},
			Expected:    []string{"bar"},
			ExpectError: true,
		},
	}

	testPlanModifier := ImmutableList()

	for i, testCase := range testCases {
		stateValue, _ := types.ListValueFrom(context.Background(), types.StringType, testCase.Old)
		planValue, _ := types.ListValueFrom(context.Background(), types.StringType, testCase.New)
		expectedValue, _ := types.ListValueFrom(context.Background(), types.StringType, testCase.Expected)

		req := planmodifier.ListRequest{
			StateValue: stateValue,
			PlanValue:  planValue,
			Path:       path.Root("test"),
		}

		var resp planmodifier.ListResponse

		testPlanModifier.PlanModifyList(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			if testCase.ExpectError == false {
				t.Fatalf("%d: got error modifying plan: %v", i, resp.Diagnostics.Errors())
			}
		}

		if !resp.PlanValue.Equal(expectedValue) {
			value, _ := resp.PlanValue.ToListValue(context.Background())
			t.Fatalf("%d: output plan value does not equal expected. Want %v plan value, got %v", i, expectedValue, value)
		}
	}
}
