package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestNonEmptyListFunction(t *testing.T) {
	t.Parallel()

	fn := NewNonEmptyListFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		values        attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		// Valid scenarios - non-empty lists
		{
			name:       "single element",
			values:     listValue([]string{"a"}),
			expectTrue: true,
		},
		{
			name:       "multiple elements",
			values:     listValue([]string{"a", "b", "c"}),
			expectTrue: true,
		},
		{
			name:       "many elements",
			values:     listValue([]string{"1", "2", "3", "4", "5"}),
			expectTrue: true,
		},

		// Invalid scenarios - empty list
		{
			name:        "empty list",
			values:      listValue([]string{}),
			expectError: true,
		},

		// Null and unknown
		{
			name:        "null list",
			values:      types.ListNull(basetypes.StringType{}),
			expectError: true,
		},
		{
			name:          "unknown list",
			values:        types.ListUnknown(basetypes.StringType{}),
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.values})}, resp)

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
