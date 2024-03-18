package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestImmutableStringSet(t *testing.T) {
	testCases := []struct {
		Old, New, Expected int64 
		ExpectError bool
	}{
		{
			Old:         0,
			New:         1234,
			Expected:    1234,
			ExpectError: false,
		},
		{
			Old:         1234,
			New:         4321,
			Expected:    0,
			ExpectError: true,
		},
	}

	testPlanModifier := ImmutableInt64()

	for i, testCase := range testCases {
		stateValue := types.Int64Value(testCase.Old)
		planValue := types.Int64Value(testCase.New)
		expectedValue := types.Int64Null() 
		if testCase.Expected != 0 {
			expectedValue = types.Int64Value(testCase.Expected)
		}

		req := planmodifier.Int64Request{
			StateValue: stateValue,
			PlanValue:  planValue,
			Path: path.Root("test"),
		}

		var resp planmodifier.Int64Response

		testPlanModifier.PlanModifyInt64(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			if testCase.ExpectError == false {
				t.Fatalf("%d: got error modifying plan: %v", i, resp.Diagnostics.Errors())
			}
		}

		if !resp.PlanValue.Equal(expectedValue) {
			t.Fatalf("%d: output plan value does not equal expected. Want %d plan value, got %d", i, expectedValue, resp.PlanValue.ValueInt64())
		}
	}
}