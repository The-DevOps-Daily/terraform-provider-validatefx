package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringPrefixFunction(t *testing.T) {
	t.Parallel()

	fn := NewStringPrefixFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "valid prefix",
			args: []attr.Value{
				types.StringValue("tf-project"),
				types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tf-"), types.StringValue("iac-")}),
				types.BoolValue(false),
			},
			expectTrue: true,
		},
		{
			name: "invalid prefix",
			args: []attr.Value{
				types.StringValue("prod-service"),
				types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tf-"), types.StringValue("iac-")}),
				types.BoolValue(false),
			},
			expectError: true,
		},
		{
			name: "case insensitive",
			args: []attr.Value{
				types.StringValue("TF-project"),
				types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tf-")}),
				types.BoolValue(true),
			},
			expectTrue: true,
		},
		{
			name: "null value",
			args: []attr.Value{
				types.StringNull(),
				types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tf-")}),
				types.BoolValue(false),
			},
			expectUnknown: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tf-")}),
				types.BoolValue(false),
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
