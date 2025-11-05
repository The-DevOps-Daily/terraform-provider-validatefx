package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func newArguments(values ...basetypes.BoolValue) function.ArgumentsData {
	list := make([]attr.Value, len(values))
	for i, v := range values {
		list[i] = v
	}

	args := []attr.Value{basetypes.NewListValueMust(basetypes.BoolType{}, list)}
	return function.NewArgumentsData(args)
}

func runComposite(t *testing.T, fn function.Function, checks ...basetypes.BoolValue) basetypes.BoolValue {
	t.Helper()

	runResp := &function.RunResponse{}
	fn.Run(context.Background(), function.RunRequest{Arguments: newArguments(checks...)}, runResp)

	if runResp.Error != nil {
		t.Fatalf("unexpected function error: %s", runResp.Error)
	}

	value := runResp.Result.Value()

	boolVal, ok := value.(basetypes.BoolValue)
	if !ok {
		t.Fatalf("unexpected result value type: %T", value)
	}

	return boolVal
}

func TestAllValidFunction(t *testing.T) {
	t.Parallel()

	fn := NewAllValidFunction()

	t.Run("all true", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(true), basetypes.NewBoolValue(true))
		if !res.ValueBool() {
			t.Fatalf("expected true, got %v", res)
		}
	})

	t.Run("contains false", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(true), basetypes.NewBoolValue(false))
		if res.ValueBool() {
			t.Fatalf("expected false, got %v", res)
		}
	})

	t.Run("unknown propagates", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(true), basetypes.NewBoolUnknown())
		if !res.IsUnknown() {
			t.Fatalf("expected unknown result")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn)
		if !res.ValueBool() {
			t.Fatalf("expected true for empty list")
		}
	})
}

func TestAnyValidFunction(t *testing.T) {
	t.Parallel()

	fn := NewAnyValidFunction()

	t.Run("any true", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(false), basetypes.NewBoolValue(true))
		if !res.ValueBool() {
			t.Fatalf("expected true result")
		}
	})

	t.Run("all false", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(false), basetypes.NewBoolValue(false))
		if res.ValueBool() {
			t.Fatalf("expected false result")
		}
	})

	t.Run("unknown propagates", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(false), basetypes.NewBoolUnknown())
		if !res.IsUnknown() {
			t.Fatalf("expected unknown result")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn)
		if res.ValueBool() {
			t.Fatalf("expected false for empty list")
		}
	})
}

func TestExactlyOneValidFunction(t *testing.T) {
	t.Parallel()

	fn := NewExactlyOneValidFunction()

	t.Run("one true", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(false), basetypes.NewBoolValue(true), basetypes.NewBoolValue(false))
		if !res.ValueBool() {
			t.Fatalf("expected true when exactly one value is true")
		}
	})

	t.Run("multiple true", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(true), basetypes.NewBoolValue(true))
		if res.ValueBool() {
			t.Fatalf("expected false when more than one value is true")
		}
	})

	t.Run("no true", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(false), basetypes.NewBoolValue(false))
		if res.ValueBool() {
			t.Fatalf("expected false when no values are true")
		}
	})

	t.Run("unknown only", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolUnknown(), basetypes.NewBoolValue(false))
		if !res.IsUnknown() {
			t.Fatalf("expected unknown result when outcomes depend on unknown values")
		}
	})

	t.Run("unknown with single true", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn, basetypes.NewBoolValue(true), basetypes.NewBoolUnknown(), basetypes.NewBoolValue(false))
		if res.IsUnknown() {
			t.Fatalf("expected deterministic result when exactly one value is true")
		}

		if res.ValueBool() {
			t.Fatalf("expected false when exactly one true value and unknown present")
		}
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()

		res := runComposite(t, fn)
		if res.ValueBool() {
			t.Fatalf("expected false for empty list")
		}
	})
}
