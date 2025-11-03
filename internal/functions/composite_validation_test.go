package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestAllValidFunctionErrorsOnInvalidList(t *testing.T) {
	t.Parallel()

	fn := NewAllValidFunction()

	resp := &function.RunResponse{}
	fn.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{
			types.ListValueMust(
				types.Int64Type,
				[]attr.Value{types.Int64Value(1)},
			),
		}),
	}, resp)

	if resp.Error == nil {
		t.Fatalf("expected error for non-boolean list elements")
	}
}

func TestCompositeFunctionsNilList(t *testing.T) {
	t.Parallel()

	args := []attr.Value{types.ListNull(basetypes.BoolType{})}

	for _, tc := range []struct {
		name string
		fn   function.Function
	}{
		{"all_valid", NewAllValidFunction()},
		{"any_valid", NewAnyValidFunction()},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			tc.fn.Run(context.Background(), function.RunRequest{Arguments: function.NewArgumentsData(args)}, resp)

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if !boolVal.IsUnknown() {
				t.Fatalf("expected unknown result for nil list")
			}
		})
	}
}

func TestSummarizeBoolValues(t *testing.T) {
	t.Parallel()

	eval := summarizeBoolValues([]basetypes.BoolValue{
		basetypes.NewBoolValue(true),
		basetypes.NewBoolUnknown(),
		basetypes.NewBoolValue(false),
	})

	if !eval.anyTrue || !eval.anyFalse || !eval.anyUnknown {
		t.Fatalf("expected summary to record true, false, and unknown")
	}
}
