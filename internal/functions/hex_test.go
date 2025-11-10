package functions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestHexFunction(t *testing.T) {
	t.Parallel()

	fn := NewHexFunction()
	ctx := context.Background()

	cases := []struct {
		name        string
		arg         types.String
		wantErr     bool
		wantUnknown bool
	}{
		{"ok lowercase", types.StringValue("deadbeef"), false, false},
		{"ok uppercase", types.StringValue("DEADBEEF"), false, false},
		{"invalid char", types.StringValue("xyz"), true, false},
		{"unknown", types.StringUnknown(), false, true},
		{"null", types.StringNull(), false, true},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			resp := &function.RunResponse{}
			args := function.NewArgumentsData([]attr.Value{tc.arg})
			fn.Run(ctx, function.RunRequest{Arguments: args}, resp)

			if tc.wantErr {
				if resp.Error == nil {
					t.Fatalf("expected error")
				}
				return
			}

			if resp.Error != nil {
				t.Fatalf("unexpected error: %s", resp.Error)
			}

			if tc.wantUnknown && !resp.Result.Value().IsUnknown() {
				t.Fatalf("expected unknown result")
			}
		})
	}
}
