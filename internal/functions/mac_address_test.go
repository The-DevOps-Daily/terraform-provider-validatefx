package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestMACAddressFunction(t *testing.T) {
	t.Parallel()

	fn := NewMACAddressFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{
			name:       "valid colon separated",
			value:      types.StringValue("00:1A:2B:3C:4D:5E"),
			expectTrue: true,
		},
		{
			name:       "valid dash separated",
			value:      types.StringValue("00-1A-2B-3C-4D-5E"),
			expectTrue: true,
		},
		{
			name:       "valid compact",
			value:      types.StringValue("001A2B3C4D5E"),
			expectTrue: true,
		},
		{
			name:        "invalid mac",
			value:       types.StringValue("AA:BB:CC:DD:EE"),
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
