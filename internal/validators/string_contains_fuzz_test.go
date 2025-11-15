package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzStringContainsValidator(f *testing.F) {
	seeds := []string{"", "Hello Terraform", "validatefx rocks", "no-match"}
	for _, s := range seeds {
		f.Add(s)
	}

	vCase := StringContains([]string{"Terraform", "ValidateFX"}, false)
	vFold := StringContains([]string{"Terraform", "ValidateFX"}, true)

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		// case-sensitive
		req := frameworkvalidator.StringRequest{Path: path.Root("contains"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		vCase.ValidateString(context.Background(), req, resp)
		expectExact := strings.Contains(s, "Terraform") || strings.Contains(s, "ValidateFX")
		if s == "" {
			expectExact = false
		}
		if expectExact != !resp.Diagnostics.HasError() {
			t.Fatalf("case-sensitive mismatch for %q: expect=%v diagErr=%v", s, expectExact, resp.Diagnostics.HasError())
		}

		// case-insensitive
		req2 := frameworkvalidator.StringRequest{Path: path.Root("contains"), ConfigValue: types.StringValue(s)}
		resp2 := &frameworkvalidator.StringResponse{}
		vFold.ValidateString(context.Background(), req2, resp2)
		lower := strings.ToLower(s)
		expectFold := strings.Contains(lower, "terraform") || strings.Contains(lower, "validatefx")
		if s == "" {
			expectFold = false
		}
		if expectFold != !resp2.Diagnostics.HasError() {
			t.Fatalf("case-insensitive mismatch for %q: expect=%v diagErr=%v", s, expectFold, resp2.Diagnostics.HasError())
		}
	})
}
