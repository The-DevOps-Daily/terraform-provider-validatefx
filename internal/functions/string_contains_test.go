package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestStringContainsFunction(t *testing.T) {
	t.Parallel()

	fn := NewStringContainsFunction()
	ctx := context.Background()

	testCases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "valid substring",
			args: []attr.Value{
				types.StringValue("hello world"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{
					basetypes.NewStringValue("hello"),
					basetypes.NewStringValue("world"),
				}),
				types.BoolValue(false),
			},
			expectTrue: true,
		},
		{
			name: "valid substring ignore case",
			args: []attr.Value{
				types.StringValue("Hello World"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{
					basetypes.NewStringValue("hello"),
				}),
				types.BoolValue(true),
			},
			expectTrue: true,
		},
		{
			name: "invalid substring",
			args: []attr.Value{
				types.StringValue("hello"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{
					basetypes.NewStringValue("world"),
				}),
				types.BoolValue(false),
			},
			expectError: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{basetypes.NewStringValue("hello")}),
				types.BoolValue(false),
			},
			expectUnknown: true,
		},
		{
			name: "unknown substrings",
			args: []attr.Value{
				types.StringValue("value"),
				basetypes.NewListUnknown(basetypes.StringType{}),
				types.BoolValue(false),
			},
			expectUnknown: true,
		},
		{
			name: "empty substrings",
			args: []attr.Value{
				types.StringValue("value"),
				types.ListValueMust(basetypes.StringType{}, []attr.Value{}),
				types.BoolValue(false),
			},
			expectUnknown: true,
		},
		{
			name: "null substrings",
			args: []attr.Value{
				types.StringValue("value"),
				types.ListNull(basetypes.StringType{}),
				types.BoolNull(),
			},
			expectError: true,
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

			result, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if tc.expectUnknown {
				if !result.IsUnknown() {
					t.Fatalf("expected unknown result")
				}
				return
			}

			if result.IsUnknown() {
				t.Fatalf("did not expect unknown result")
			}

			if result.ValueBool() != tc.expectTrue {
				t.Fatalf("expected %t, got %t", tc.expectTrue, result.ValueBool())
			}
		})
	}
}
