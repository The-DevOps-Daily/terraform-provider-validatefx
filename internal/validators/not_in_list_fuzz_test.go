package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzNotInListValidator fuzzes the not-in-list validator. It ensures that
// values present in the disallowed set fail, and values outside pass.
func FuzzNotInListValidator(f *testing.F) {
	seeds := []string{"", "alpha", " BETA ", "γamma", "delta", "ALPHA", "💡idea", " space ", "tab\tchar"}
	for _, s := range seeds {
		f.Add(s)
	}

	disallowed := []string{"alpha", "beta", "γamma", "💡idea"}
	vCase := NewNotInListValidator(disallowed, false)
	vNoCase := NewNotInListValidator(disallowed, true)

	dis := map[string]struct{}{"alpha": {}, "beta": {}, "γamma": {}, "💡idea": {}}

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		// exact (case sensitive)
		req := frameworkvalidator.StringRequest{Path: path.Root("not_in_list"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		vCase.ValidateString(context.Background(), req, resp)

		_, isDis := dis[s]
		expectError := isDis // should error if in disallowed
		if expectError != resp.Diagnostics.HasError() {
			t.Fatalf("case-sensitive mismatch for %q: expectError=%v hasErr=%v", s, expectError, resp.Diagnostics.HasError())
		}

		// case-insensitive
		req2 := frameworkvalidator.StringRequest{Path: path.Root("not_in_list"), ConfigValue: types.StringValue(s)}
		resp2 := &frameworkvalidator.StringResponse{}
		vNoCase.ValidateString(context.Background(), req2, resp2)

		lower := strings.ToLower(s)
		_, isDisFold := dis[lower]
		expectError2 := isDisFold
		if expectError2 != resp2.Diagnostics.HasError() {
			t.Fatalf("case-insensitive mismatch for %q: expectError=%v hasErr=%v", s, expectError2, resp2.Diagnostics.HasError())
		}
	})
}
