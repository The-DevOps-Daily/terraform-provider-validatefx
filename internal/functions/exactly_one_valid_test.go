package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestExactlyOneValidFunctionRun(t *testing.T) {
	t.Parallel()

	fn := NewExactlyOneValidFunction()
	ctx := context.Background()

	run := func(values []attr.Value) *function.RunResponse {
		arguments := function.NewArgumentsData([]attr.Value{
			basetypes.NewListValueMust(basetypes.BoolType{}, values),
		})
		resp := &function.RunResponse{}
		fn.Run(ctx, function.RunRequest{Arguments: arguments}, resp)
		return resp
	}

	t.Run("exactly one true", func(t *testing.T) {
		t.Parallel()

		resp := run([]attr.Value{
			basetypes.NewBoolValue(false),
			basetypes.NewBoolValue(true),
			basetypes.NewBoolValue(false),
		})

		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
		if !ok || !boolVal.ValueBool() {
			t.Fatalf("expected true result, got %v", resp.Result.Value())
		}
	})

	t.Run("multiple true -> false", func(t *testing.T) {
		t.Parallel()

		resp := run([]attr.Value{
			basetypes.NewBoolValue(true),
			basetypes.NewBoolValue(true),
		})

		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, _ := resp.Result.Value().(basetypes.BoolValue)
		if boolVal.ValueBool() {
			t.Fatalf("expected false when multiple values are true")
		}
	})

	t.Run("none true -> false", func(t *testing.T) {
		t.Parallel()

		resp := run([]attr.Value{
			basetypes.NewBoolValue(false),
			basetypes.NewBoolValue(false),
		})

		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, _ := resp.Result.Value().(basetypes.BoolValue)
		if boolVal.ValueBool() {
			t.Fatalf("expected false when no values are true")
		}
	})

	t.Run("unknown result", func(t *testing.T) {
		t.Parallel()

		resp := run([]attr.Value{
			basetypes.NewBoolValue(false),
			basetypes.NewBoolUnknown(),
			basetypes.NewBoolValue(false),
		})

		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, _ := resp.Result.Value().(basetypes.BoolValue)
		if !boolVal.IsUnknown() {
			t.Fatalf("expected unknown result when outcomes depend on unknown values")
		}
	})

	t.Run("invalid element type", func(t *testing.T) {
		t.Parallel()

		arguments := function.NewArgumentsData([]attr.Value{
			basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("invalid")}),
		})

		resp := &function.RunResponse{}
		fn.Run(ctx, function.RunRequest{Arguments: arguments}, resp)

		if resp.Error == nil {
			t.Fatalf("expected error for invalid element type")
		}
	})

	t.Run("null list -> unknown", func(t *testing.T) {
		t.Parallel()

		arguments := function.NewArgumentsData([]attr.Value{basetypes.NewListNull(basetypes.BoolType{})})
		resp := &function.RunResponse{}
		fn.Run(ctx, function.RunRequest{Arguments: arguments}, resp)

		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}

		boolVal, _ := resp.Result.Value().(basetypes.BoolValue)
		if !boolVal.IsUnknown() {
			t.Fatalf("expected unknown result for null list")
		}
	})
}
