package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestBetweenFunction(t *testing.T) {
	t.Parallel()

	fn := NewBetweenFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		args          []attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name: "within range",
			args: []attr.Value{
				types.StringValue("7.5"),
				types.StringValue("5"),
				types.StringValue("10"),
			},
			expectTrue: true,
		},
		{
			name: "too low",
			args: []attr.Value{
				types.StringValue("2"),
				types.StringValue("5"),
				types.StringValue("10"),
			},
			expectError: true,
		},
		{
			name: "too high",
			args: []attr.Value{
				types.StringValue("11"),
				types.StringValue("5"),
				types.StringValue("10"),
			},
			expectError: true,
		},
		{
			name: "invalid bounds",
			args: []attr.Value{
				types.StringValue("7"),
				types.StringValue("10"),
				types.StringValue("5"),
			},
			expectError: true,
		},
		{
			name: "unknown value",
			args: []attr.Value{
				types.StringUnknown(),
				types.StringValue("5"),
				types.StringValue("10"),
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
