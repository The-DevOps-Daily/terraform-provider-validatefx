package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TestListSubsetValidator_ListAndSet(t *testing.T) {
	t.Parallel()

	v := NewListSubset([]string{"read", "write", "admin"})

	runList := func(values []string) *frameworkvalidator.ListResponse {
		list := basetypes.NewListValueMust(basetypes.StringType{}, toAttr(values))
		req := frameworkvalidator.ListRequest{Path: path.Root("value"), ConfigValue: list}
		resp := &frameworkvalidator.ListResponse{}
		v.ValidateList(context.Background(), req, resp)
		return resp
	}

	runSet := func(values []string) *frameworkvalidator.SetResponse {
		set := basetypes.NewSetValueMust(basetypes.StringType{}, toAttr(values))
		req := frameworkvalidator.SetRequest{Path: path.Root("value"), ConfigValue: set}
		resp := &frameworkvalidator.SetResponse{}
		v.ValidateSet(context.Background(), req, resp)
		return resp
	}

	if resp := runList([]string{"read", "write"}); resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
	}
	if resp := runSet([]string{"read", "unknown"}); !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics for disallowed element")
	}
}

func toAttr(vs []string) []attr.Value {
	out := make([]attr.Value, 0, len(vs))
	for _, s := range vs {
		out = append(out, basetypes.NewStringValue(s))
	}
	return out
}
