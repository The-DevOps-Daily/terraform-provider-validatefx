package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestJWTFunction(t *testing.T) {
	t.Parallel()

	fn := NewJWTFunction()
	ctx := context.Background()

	valid := types.StringValue("eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

	cases := []struct {
		name          string
		arg           attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{name: "valid", arg: valid, expectTrue: true},
		{name: "missing parts", arg: types.StringValue("abc.def"), expectError: true},
		{name: "empty segment", arg: types.StringValue("abc..def"), expectError: true},
		{name: "bad base64", arg: types.StringValue("abc.def.!@#"), expectError: true},
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
