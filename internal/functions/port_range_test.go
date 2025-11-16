package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestPortRangeFunction(t *testing.T) {
	t.Parallel()

	fn := NewPortRangeFunction()
	ctx := context.Background()

	cases := []struct {
		name        string
		value       attr.Value
		expectError bool
		expectTrue  bool
		expectUnk   bool
	}{
		{"ok", types.StringValue("80-8080"), false, true, false},
		{"spaces", types.StringValue(" 0 - 65535 "), false, true, false},
		{"bad", types.StringValue("80:8080"), true, false, false},
		{"unknown", types.StringUnknown(), false, false, true},
		{"null", types.StringNull(), false, false, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			req := function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}
			fn.Run(ctx, req, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			b, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}

			if tc.expectUnk {
				if !b.IsUnknown() {
					t.Fatalf("expected unknown")
				}
				return
			}

			if b.IsUnknown() {
				t.Fatalf("did not expect unknown")
			}
			if b.ValueBool() != tc.expectTrue {
				t.Fatalf("expected %t got %t", tc.expectTrue, b.ValueBool())
			}
		})
	}
}
