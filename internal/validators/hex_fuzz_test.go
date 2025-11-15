package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzHexValidator(f *testing.F) {
	seeds := []string{"", "deadbeef", "CAFEBABE", "12345g", "0011aa", "ZZZ"}
	for _, s := range seeds {
		f.Add(s)
	}
	v := Hex()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("hex"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if !resp.Diagnostics.HasError() {
				t.Fatalf("empty should error")
			}
			return
		}
		expect := hexRe.MatchString(s)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: hex=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
