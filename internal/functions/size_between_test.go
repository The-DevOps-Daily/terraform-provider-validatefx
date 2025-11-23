package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestSizeBetweenFunction(t *testing.T) {
	t.Parallel()

	fn := NewSizeBetweenFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		value         attr.Value
		min           attr.Value
		max           attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid values within range
		{
			name:       "value in range",
			value:      types.StringValue("5"),
			min:        types.StringValue("1"),
			max:        types.StringValue("10"),
			expectTrue: true,
		},
		{
			name:       "value at min",
			value:      types.StringValue("1"),
			min:        types.StringValue("1"),
			max:        types.StringValue("10"),
			expectTrue: true,
		},
		{
			name:       "value at max",
			value:      types.StringValue("10"),
			min:        types.StringValue("1"),
			max:        types.StringValue("10"),
			expectTrue: true,
		},
		{
			name:       "decimal in range",
			value:      types.StringValue("0.5"),
			min:        types.StringValue("0"),
			max:        types.StringValue("1"),
			expectTrue: true,
		},

		// Invalid values
		{
			name:        "value below min",
			value:       types.StringValue("0"),
			min:         types.StringValue("1"),
			max:         types.StringValue("10"),
			expectError: true,
		},
		{
			name:        "value above max",
			value:       types.StringValue("11"),
			min:         types.StringValue("1"),
			max:         types.StringValue("10"),
			expectError: true,
		},
		{
			name:        "non-numeric value",
			value:       types.StringValue("abc"),
			min:         types.StringValue("1"),
			max:         types.StringValue("10"),
			expectError: true,
		},

		// Null and unknown
		{
			name:          "null input",
			value:         types.StringNull(),
			min:           types.StringValue("1"),
			max:           types.StringValue("10"),
			expectUnknown: true,
		},
		{
			name:          "unknown input",
			value:         types.StringUnknown(),
			min:           types.StringValue("1"),
			max:           types.StringValue("10"),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value, tc.min, tc.max})}, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			boolVal, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if tc.expectUnknown {
				if !boolVal.IsUnknown() {
					t.Fatalf("expected unknown result")
				}
				return
			}

			if boolVal.IsUnknown() {
				t.Fatalf("did not expect unknown result")
			}

			if boolVal.ValueBool() != tc.expectTrue {
				t.Fatalf("expected %t, got %t", tc.expectTrue, boolVal.ValueBool())
			}
		})
	}
}
