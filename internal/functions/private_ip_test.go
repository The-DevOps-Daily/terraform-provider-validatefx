package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestPrivateIPFunction(t *testing.T) {
	t.Parallel()

	fn := NewPrivateIPFunction()
	ctx := context.Background()

	cases := []struct {
		name       string
		value      attr.Value
		expectErr  bool
		expectTrue bool
		expectUnk  bool
	}{
		{"ipv4-rfc1918-10", types.StringValue("10.0.0.1"), false, true, false},
		{"ipv4-rfc1918-172", types.StringValue("172.16.10.5"), false, true, false},
		{"ipv4-rfc1918-192", types.StringValue("192.168.1.10"), false, true, false},
		{"ipv6-ula-fd", types.StringValue("fd00::1"), false, true, false},
		{"public-ip", types.StringValue("8.8.8.8"), true, false, false},
		{"unknown", types.StringUnknown(), false, false, true},
		{"null", types.StringNull(), false, false, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			var resp function.RunResponse
			req := function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}
			fn.Run(ctx, req, &resp)

			if tc.expectErr {
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
