package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMatchesRegexFunction(t *testing.T) {
	t.Parallel()

	fn := NewMatchesRegexFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "match",
			args: []attr.Value{
				types.StringValue("user_123"),
				types.StringValue("^[a-z0-9_]+$"),
			},
			expectTrue: true,
		},
		{
			name: "mismatch",
			args: []attr.Value{
				types.StringValue("Invalid-User"),
				types.StringValue("^[a-z0-9_]+$"),
			},
			expectError: true,
		},
		{
			name: "invalid pattern",
			args: []attr.Value{
				types.StringValue("value"),
				types.StringValue("("),
			},
			expectError: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				types.StringValue(".*"),
			},
			expectUnknown: true,
		},
		{
			name: "unknown pattern",
			args: []attr.Value{
				types.StringValue("test"),
				types.StringUnknown(),
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
