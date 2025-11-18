package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestARNFunction(t *testing.T) {
	t.Parallel()
	fn := NewARNFunction()
	ctx := context.Background()

	cases := []struct {
		name          string
		value         attr.Value
		expectError   bool
		expectUnknown bool
		expectTrue    bool
	}{
		{name: "valid iam role", value: types.StringValue("arn:aws:iam::123456789012:role/Admin"), expectTrue: true},
		{name: "invalid format", value: types.StringValue("not-an-arn"), expectError: true},
		{name: "null", value: types.StringNull(), expectUnknown: true},
		{name: "unknown", value: types.StringUnknown(), expectUnknown: true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			fn.Run(ctx, function.RunRequest{Arguments: function.NewArgumentsData([]attr.Value{tc.value})}, resp)

			if tc.expectError {
				if resp.Error == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if resp.Error != nil {
				t.Fatalf("unexpected error: %v", resp.Error)
			}
			b, ok := resp.Result.Value().(basetypes.BoolValue)
			if !ok {
				t.Fatalf("unexpected result type %T", resp.Result.Value())
			}
			if tc.expectUnknown {
				if !b.IsUnknown() {
					t.Fatalf("expected unknown")
				}
				return
			}
			if b.IsUnknown() {
				t.Fatalf("did not expect unknown")
			}
			if !b.ValueBool() {
				t.Fatalf("expected true result")
			}
		})
	}
}
