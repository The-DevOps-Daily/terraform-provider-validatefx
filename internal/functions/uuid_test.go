package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestUUIDFunction(t *testing.T) {
	t.Parallel()

	fn := NewUUIDFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid uuid",
			value:      types.StringValue("d9428888-122b-11e1-b85c-61cd3cbb3210"),
			expectTrue: true,
		},
		{
			name:        "invalid uuid",
			value:       types.StringValue("not-a-uuid"),
			expectError: true,
		},
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

	for _, tc := range cases {
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
