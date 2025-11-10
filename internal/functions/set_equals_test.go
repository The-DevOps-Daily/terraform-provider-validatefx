package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestSetEqualsFunction(t *testing.T) {
	t.Parallel()

	fn := NewSetEqualsFunction()
	ctx := context.Background()

	list := func(values ...string) attr.Value {
		attrs := make([]attr.Value, 0, len(values))
		for _, v := range values {
			attrs = append(attrs, types.StringValue(v))
		}
		return types.ListValueMust(types.StringType, attrs)
	}

	makeArgs := func(values, expected attr.Value) []attr.Value {
		return []attr.Value{values, expected}
	}

	testCases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "matching sets",
			args:       makeArgs(list("a", "b", "c"), list("c", "b", "a")),
			expectTrue: true,
		},
		{
			name:        "mismatched sets",
			args:        makeArgs(list("a", "b"), list("a", "c")),
			expectError: true,
		},
		{
			name:        "missing value list",
			args:        makeArgs(types.ListNull(types.StringType), list("a")),
			expectError: true,
		},
		{
			name: "unknown element",
			args: []attr.Value{
				types.ListValueMust(types.StringType, []attr.Value{types.StringUnknown()}),
				list("a"),
			},
			expectUnknown: true,
		},
		{
			name: "unknown expected list",
			args: []attr.Value{
				list("a"),
				types.ListValueMust(types.StringType, []attr.Value{types.StringUnknown()}),
			},
			expectUnknown: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData(tc.args)}, resp)

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
