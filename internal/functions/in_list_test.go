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

func TestInListFunction_Metadata(t *testing.T) {
	fn := NewInListFunction()
	resp := &function.MetadataResponse{}
	fn.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "in_list" {
		t.Errorf("expected name 'in_list', got %q", resp.Name)
	}
}

func TestInListFunction_Definition(t *testing.T) {
	fn := NewInListFunction()
	resp := &function.DefinitionResponse{}
	fn.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if resp.Definition.Return == nil {
		t.Fatal("expected return definition, got nil")
	}
	if len(resp.Definition.Parameters) != 4 {
		t.Errorf("expected 4 parameters, got %d", len(resp.Definition.Parameters))
	}
}

func TestInListFunction_CustomMessage(t *testing.T) {
	t.Parallel()

	fn := NewInListFunction()
	ctx := context.Background()

	allowed := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("alpha"),
		basetypes.NewStringValue("beta"),
	})

	// Test with custom message on failure
	args := function.NewArgumentsData([]attr.Value{
		types.StringValue("gamma"),
		allowed,
		types.BoolValue(false),
		types.StringValue("Custom error message"),
	})

	resp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, resp)

	if resp.Error == nil {
		t.Fatalf("expected error with custom message")
	}
}

func TestInListFunction_NullAllowed(t *testing.T) {
	t.Parallel()

	fn := NewInListFunction()
	ctx := context.Background()

	args := function.NewArgumentsData([]attr.Value{
		types.StringValue("alpha"),
		basetypes.NewListNull(basetypes.StringType{}),
		types.BoolValue(false),
		types.StringNull(),
	})

	resp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, resp)

	if resp.Error == nil {
		t.Fatalf("expected error for null allowed list")
	}
}

func TestInListFunction_UnknownMessage(t *testing.T) {
	t.Parallel()

	fn := NewInListFunction()
	ctx := context.Background()

	allowed := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("alpha"),
	})

	args := function.NewArgumentsData([]attr.Value{
		types.StringValue("alpha"),
		allowed,
		types.BoolValue(false),
		types.StringUnknown(),
	})

	resp := &function.RunResponse{}
	fn.Run(ctx, function.RunRequest{Arguments: args}, resp)

	if resp.Error != nil {
		t.Fatalf("unexpected error: %s", resp.Error)
	}

	boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
	if !ok || !boolVal.IsUnknown() {
		t.Fatalf("expected unknown result when message is unknown")
	}
}

func TestSelectInListValidator_WithMessage(t *testing.T) {
	allowed := []string{"alpha", "beta"}
	message := types.StringValue("Custom validation message")

	validator := selectInListValidator(allowed, false, message)
	if validator == nil {
		t.Fatal("expected validator, got nil")
	}
}

func TestSelectInListValidator_WithoutMessage(t *testing.T) {
	allowed := []string{"alpha", "beta"}

	// Test with null message
	validator := selectInListValidator(allowed, false, types.StringNull())
	if validator == nil {
		t.Fatal("expected validator, got nil")
	}

	// Test with empty string message
	validator = selectInListValidator(allowed, true, types.StringValue(""))
	if validator == nil {
		t.Fatal("expected validator, got nil")
	}

	// Test with unknown message
	validator = selectInListValidator(allowed, false, types.StringUnknown())
	if validator == nil {
		t.Fatal("expected validator, got nil")
	}
}

func TestStringArgument_Error(t *testing.T) {
	ctx := context.Background()
	resp := &function.RunResponse{}

	// Test with invalid arguments (empty)
	req := function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{})}
	val, state, ok := stringArgument(ctx, req, resp, 0)

	if ok {
		t.Fatal("expected error when accessing invalid argument index")
	}
	if resp.Error == nil {
		t.Fatal("expected error to be set")
	}
	if state != valueKnown {
		t.Fatal("expected valueKnown state")
	}
	if !val.IsNull() {
		t.Fatal("expected null string value")
	}
}

func TestIgnoreCaseFlag_Error(t *testing.T) {
	ctx := context.Background()
	resp := &function.RunResponse{}

	// Test with invalid arguments
	req := function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{})}
	val, state, ok := ignoreCaseFlag(ctx, req, resp, 0)

	if ok {
		t.Fatal("expected error when accessing invalid argument index")
	}
	if resp.Error == nil {
		t.Fatal("expected error to be set")
	}
	if val {
		t.Fatal("expected false value")
	}
	if state != valueKnown {
		t.Fatal("expected valueKnown state")
	}
}

func TestBoolFromOptional(t *testing.T) {
	tests := []struct {
		name        string
		input       types.Bool
		expectVal   bool
		expectState valueState
	}{
		{
			name:        "true value",
			input:       types.BoolValue(true),
			expectVal:   true,
			expectState: valueKnown,
		},
		{
			name:        "false value",
			input:       types.BoolValue(false),
			expectVal:   false,
			expectState: valueKnown,
		},
		{
			name:        "null value",
			input:       types.BoolNull(),
			expectVal:   false,
			expectState: valueKnown,
		},
		{
			name:        "unknown value",
			input:       types.BoolUnknown(),
			expectVal:   false,
			expectState: valueUnknown,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			val, state := boolFromOptional(tt.input)
			if val != tt.expectVal {
				t.Errorf("expected value %v, got %v", tt.expectVal, val)
			}
			if state != tt.expectState {
				t.Errorf("expected state %v, got %v", tt.expectState, state)
			}
		})
	}
}

func TestMessageArgument(t *testing.T) {
	ctx := context.Background()

	t.Run("valid message", func(t *testing.T) {
		resp := &function.RunResponse{}
		args := function.NewArgumentsData([]attr.Value{types.StringValue("test message")})
		req := function.RunRequest{Arguments: args}

		msg, state, ok := messageArgument(ctx, req, resp, 0)
		if !ok {
			t.Fatal("expected success")
		}
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}
		if state != valueKnown {
			t.Fatal("expected valueKnown state")
		}
		if msg.ValueString() != "test message" {
			t.Errorf("expected 'test message', got %q", msg.ValueString())
		}
	})

	t.Run("unknown message", func(t *testing.T) {
		resp := &function.RunResponse{}
		args := function.NewArgumentsData([]attr.Value{types.StringUnknown()})
		req := function.RunRequest{Arguments: args}

		msg, state, ok := messageArgument(ctx, req, resp, 0)
		if !ok {
			t.Fatal("expected success")
		}
		if resp.Error != nil {
			t.Fatalf("unexpected error: %s", resp.Error)
		}
		if state != valueUnknown {
			t.Fatal("expected valueUnknown state")
		}
		if !msg.IsUnknown() {
			t.Fatal("expected unknown message")
		}
	})

	t.Run("error case", func(t *testing.T) {
		resp := &function.RunResponse{}
		// Empty arguments - will cause error
		req := function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{})}

		_, state, ok := messageArgument(ctx, req, resp, 0)
		if ok {
			t.Fatal("expected error")
		}
		if resp.Error == nil {
			t.Fatal("expected error to be set")
		}
		if state != valueKnown {
			t.Fatal("expected valueKnown state")
		}
	})
}

func TestPrepareAllowedValues_WithNullUnknownElements(t *testing.T) {
	ctx := context.Background()

	// List with mix of known, null, and unknown elements
	list := basetypes.NewListValueMust(basetypes.StringType{}, []attr.Value{
		basetypes.NewStringValue("alpha"),
		basetypes.NewStringNull(),
		basetypes.NewStringValue("beta"),
		basetypes.NewStringUnknown(),
		basetypes.NewStringValue("gamma"),
	})

	values, state, err := prepareAllowedValues(ctx, list)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if state != valueKnown {
		t.Fatal("expected valueKnown state")
	}
	// Should only contain non-null, non-unknown values
	if len(values) != 3 {
		t.Errorf("expected 3 values, got %d", len(values))
	}
	expected := []string{"alpha", "beta", "gamma"}
	for i, exp := range expected {
		if values[i] != exp {
			t.Errorf("expected values[%d] = %q, got %q", i, exp, values[i])
		}
	}
}
