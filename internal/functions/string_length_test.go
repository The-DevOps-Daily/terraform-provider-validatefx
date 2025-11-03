package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringLengthFunction(t *testing.T) {
	t.Parallel()

	fn := NewStringLengthFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "within bounds",
			args: []attr.Value{
				types.StringValue("hello"),
				types.Int64Value(3),
				types.Int64Value(10),
			},
			expectTrue: true,
		},
		{
			name: "too short",
			args: []attr.Value{
				types.StringValue("hi"),
				types.Int64Value(3),
				types.Int64Value(10),
			},
			expectError: true,
		},
		{
			name: "too long",
			args: []attr.Value{
				types.StringValue("this is a really long string"),
				types.Int64Value(3),
				types.Int64Value(10),
			},
			expectError: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				types.Int64Null(),
				types.Int64Null(),
			},
			expectUnknown: true,
		},
	}

	for _, tc := range cases {
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
