package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestNonNegativeNumberFunction(t *testing.T) {
	t.Parallel()

	fn := NewNonNegativeNumberFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid non-negative numbers
		{
			name:       "zero",
			value:      types.StringValue("0"),
			expectTrue: true,
		},
		{
			name:       "positive integer",
			value:      types.StringValue("42"),
			expectTrue: true,
		},
		{
			name:       "positive decimal",
			value:      types.StringValue("3.14"),
			expectTrue: true,
		},
		{
			name:       "small positive",
			value:      types.StringValue("0.001"),
			expectTrue: true,
		},
		{
			name:       "zero decimal",
			value:      types.StringValue("0.0"),
			expectTrue: true,
		},

		// Invalid values
		{
			name:        "negative integer",
			value:       types.StringValue("-1"),
			expectError: true,
		},
		{
			name:        "negative decimal",
			value:       types.StringValue("-3.14"),
			expectError: true,
		},
		{
			name:        "non-numeric",
			value:       types.StringValue("abc"),
			expectError: true,
		},

		// Null and unknown
		{
			name:          "null input",
			value:         types.StringNull(),
			expectUnknown: true,
		},
		{
			name:          "unknown input",
			value:         types.StringUnknown(),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}, resp)

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
