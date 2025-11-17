package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func FuzzListSubsetValidator(f *testing.F) {
	v := NewListSubset([]string{"a", "b", "c"})
	for _, s := range []string{"a", "a,b", "d", ""} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, csv string) {
		vals := splitCSV(csv)
		list := basetypes.NewListValueMust(basetypes.StringType{}, stringSliceToAttr(vals))
		req := frameworkvalidator.ListRequest{Path: path.Root("value"), ConfigValue: list}
		resp := &frameworkvalidator.ListResponse{}
		v.ValidateList(context.Background(), req, resp)
	})
}

func stringSliceToAttr(vs []string) []attr.Value {
	out := make([]attr.Value, 0, len(vs))
	for _, s := range vs {
		out = append(out, basetypes.NewStringValue(s))
	}
	return out
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	var out []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
