package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestInListFunction(t *testing.T) {
	t.Parallel()

	fn := NewInListFunction()
	ctx := context.Background()

	run := func(value attr.Value, allowed attr.Value, ignore attr.Value) *function.RunResponse {
		resp := &function.RunResponse{}
		// Pass null for the optional custom message parameter by default.
		args := function.NewArgumentsData([]attr.Value{value, allowed, ignore, types.StringNull()})
		fn.Run(ctx, function.RunRequest{Arguments: args}, resp)
		return resp
	}

	allowed := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("alpha"),
		basetypes.NewStringValue("beta"),
	})

	t.Run("valid value", func(t *testing.T) {
		t.Parallel()

		resp := run(types.StringValue("alpha"), allowed, types.BoolValue(false))
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !boolVal.ValueBool() {
			t.Fatalf("expected true result")
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		t.Parallel()

		resp := run(types.StringValue("gamma"), allowed, types.BoolValue(false))
		if resp.Error == nil {
			t.Fatalf("expected error for invalid value")
		}
	})

	t.Run("ignore case", func(t *testing.T) {
		t.Parallel()

		resp := run(types.StringValue("BETA"), allowed, types.BoolValue(true))
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !boolVal.ValueBool() {
			t.Fatalf("expected true result with ignore case")
		}
	})

	t.Run("unknown allowed", func(t *testing.T) {
		t.Parallel()

		resp := run(types.StringValue("alpha"), basetypes.NewListUnknown(basetypes.StringType{}), types.BoolValue(false))
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !boolVal.IsUnknown() {
			t.Fatalf("expected unknown result for unknown allowed list")
		}
	})

	t.Run("empty allowed", func(t *testing.T) {
		t.Parallel()

		resp := run(types.StringValue("alpha"), basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{}), types.BoolValue(false))
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !boolVal.IsUnknown() {
			t.Fatalf("expected unknown result when allowed values empty")
		}
	})

	t.Run("invalid allowed element type", func(t *testing.T) {
		t.Parallel()

		invalidAllowed := basetypes.NewListValueMust(basetypes.Int64Type{}, []attr.Value{basetypes.NewInt64Value(1)})
		resp := run(types.StringValue("alpha"), invalidAllowed, types.BoolNull())
		if resp.Error == nil {
			t.Fatalf("expected error for non-string allowed values")
		}
	})

	t.Run("unknown ignore case", func(t *testing.T) {
		t.Parallel()

		resp := run(types.StringValue("alpha"), allowed, types.BoolUnknown())
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !boolVal.IsUnknown() {
			t.Fatalf("expected unknown result when ignore_case unknown")
		}
	})
}
