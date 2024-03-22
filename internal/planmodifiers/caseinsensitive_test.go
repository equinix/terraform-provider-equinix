package planmodifiers

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCaseInsensitiveSet(t *testing.T) {
	testCases := []struct {
		Old, New, Expected string
	}{
		{
			Old:      "foo",
			New:      "foo",
			Expected: "foo",
		},
		{
			Old:      "Bar",
			New:      "bar",
			Expected: "Bar",
		},
		{
			Old:      "foo",
			New:      "fOO",
			Expected: "foo",
		},
	}

	testPlanModifier := CaseInsensitiveString()

	for i, testCase := range testCases {
		stateValue := types.StringValue(testCase.Old)
		planValue := types.StringValue(testCase.New)
		expectedValue := types.StringValue(testCase.Expected)

		req := planmodifier.StringRequest{
			StateValue: stateValue,
			PlanValue:  planValue,
		}

		var resp planmodifier.StringResponse

		testPlanModifier.PlanModifyString(context.Background(), req, &resp)

		if resp.Diagnostics.HasError() {
			t.Fatalf("%d: got error modifying plan: %v", i, resp.Diagnostics.Errors())
		}

		if !resp.PlanValue.Equal(expectedValue) {
			t.Fatalf("%d: output plan value does not equal expected plan value", i)
		}
	}
}
