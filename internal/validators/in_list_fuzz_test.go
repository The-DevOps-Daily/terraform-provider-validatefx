package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzInListValidator fuzzes the in-list validator with varied seeds
// including Unicode, whitespace, and mixed case. It checks that values
// present in the allowed set pass, and values outside fail.
func FuzzInListValidator(f *testing.F) {
	seeds := []string{"", "alpha", " BETA ", "γamma", "delta", "ALPHA", "💡idea", " space ", "tab\tchar"}
	for _, s := range seeds {
		f.Add(s)
	}

	allowed := []string{"alpha", "beta", "γamma", "💡idea"}
	vCase := NewInListValidator(allowed, false)
	vNoCase := NewInListValidator(allowed, true)

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		// exact (case sensitive)
		req := frameworkvalidator.StringRequest{Path: path.Root("in_list"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		vCase.ValidateString(context.Background(), req, resp)

		// normalize like validator for comparison purposes
		exactAllowed := map[string]struct{}{"alpha": {}, "beta": {}, "γamma": {}, "💡idea": {}}
		_, expectExact := exactAllowed[s]
		if expectExact != !resp.Diagnostics.HasError() {
			t.Fatalf("case-sensitive mismatch for %q: expect=%v hasErr=%v", s, expectExact, resp.Diagnostics.HasError())
		}

		// case-insensitive
		req2 := frameworkvalidator.StringRequest{Path: path.Root("in_list"), ConfigValue: types.StringValue(s)}
		resp2 := &frameworkvalidator.StringResponse{}
		vNoCase.ValidateString(context.Background(), req2, resp2)

		// expect true if lowercased version is in set
		lower := strings.ToLower(s)
		_, expectFold := exactAllowed[lower]
		if expectFold != !resp2.Diagnostics.HasError() {
			t.Fatalf("case-insensitive mismatch for %q: expect=%v hasErr=%v", s, expectFold, resp2.Diagnostics.HasError())
		}
	})
}
