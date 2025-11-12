package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestFQDNFunction(t *testing.T) {
	t.Parallel()

	fn := NewFQDNFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		arg           attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{name: "valid", arg: types.StringValue("example.com"), expectTrue: true},
		{name: "valid punycode", arg: types.StringValue("xn--bcher-kva.example"), expectTrue: true},
		{name: "no dot", arg: types.StringValue("localhost"), expectError: true},
		{name: "empty label", arg: types.StringValue("example..com"), expectError: true},
		{name: "bad chars", arg: types.StringValue("exa_mple.com"), expectError: true},
		{name: "unknown", arg: types.StringUnknown(), expectUnknown: true},
		{name: "null", arg: types.StringNull(), expectUnknown: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.arg})}, resp)

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

			if !boolVal.ValueBool() {
				t.Fatalf("expected true result")
			}
		})
	}
}
