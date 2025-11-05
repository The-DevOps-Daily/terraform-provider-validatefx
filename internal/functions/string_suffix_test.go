package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringSuffixFunction(t *testing.T) {
	t.Parallel()

	fn := NewStringSuffixFunction()
	ctx := context.Background()

	tests := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "valid suffix",
			args: []attr.Value{
				types.StringValue("config.yaml"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{
					basetypes.NewStringValue(".yaml"),
					basetypes.NewStringValue(".json"),
				}),
			},
			expectTrue: true,
		},
		{
			name: "invalid suffix",
			args: []attr.Value{
				types.StringValue("config.yml"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{
					basetypes.NewStringValue(".yaml"),
				}),
			},
			expectError: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue(".log")}),
			},
			expectUnknown: true,
		},
		{
			name: "unknown suffix list",
			args: []attr.Value{
				types.StringValue("config.yaml"),
				basetypes.NewListUnknown(basetypes.StringType{}),
			},
			expectUnknown: true,
		},
		{
			name: "empty suffix list",
			args: []attr.Value{
				types.StringValue("config.yaml"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{}),
			},
			expectUnknown: true,
		},
		{
			name: "non-string suffix entry",
			args: []attr.Value{
				types.StringValue("config.yaml"),
				types.ListValueMust(basetypes.Int64Type{}, []attr.Value{basetypes.NewInt64Value(1)}),
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
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
