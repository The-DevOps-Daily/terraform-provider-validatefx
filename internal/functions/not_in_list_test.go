package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNotInListFunction(t *testing.T) {
	t.Parallel()

	fn := NewNotInListFunction()
	ctx := context.Background()

	list := func(values ...string) attr.Value {
		attrs := make([]attr.Value, 0, len(values))
		for _, v := range values {
			attrs = append(attrs, types.StringValue(v))
		}
		return types.ListValueMust(types.StringType, attrs)
	}

	cases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
	}{
		{
			name:        "disallowed present",
			args:        []attr.Value{types.StringValue("beta"), list("alpha", "beta"), types.BoolValue(false)},
			expectError: true,
		},
		{
			name:        "allowed value (not present)",
			args:        []attr.Value{types.StringValue("delta"), list("alpha", "beta"), types.BoolValue(false)},
			expectError: false,
		},
		{
			name:          "unknown disallowed list",
			args:          []attr.Value{types.StringValue("delta"), types.ListValueMust(types.StringType, []attr.Value{types.StringUnknown()}), types.BoolNull()},
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

			if tc.expectUnknown {
				if !resp.Result.Value().IsUnknown() {
					t.Fatalf("expected unknown result")
				}
				return
			}
		})
	}
}
