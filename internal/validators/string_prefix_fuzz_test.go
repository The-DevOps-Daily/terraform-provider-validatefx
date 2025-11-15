package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzStringPrefixValidator(f *testing.F) {
	seeds := []string{"", "tf-backend", "IAC-service", "misc", " tf- leading"}
	for _, s := range seeds {
		f.Add(s)
	}

	vCase := StringPrefix([]string{"tf-", "iac-"}, false)
	vFold := StringPrefix([]string{"tf-", "iac-"}, true)

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		// case-sensitive
		req := frameworkvalidator.StringRequest{Path: path.Root("prefix"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		vCase.ValidateString(context.Background(), req, resp)

		expectExact := strings.HasPrefix(s, "tf-") || strings.HasPrefix(s, "iac-")
		if expectExact != !resp.Diagnostics.HasError() {
			t.Fatalf("case-sensitive mismatch for %q: expect=%v diagErr=%v", s, expectExact, resp.Diagnostics.HasError())
		}

		// case-insensitive
		req2 := frameworkvalidator.StringRequest{Path: path.Root("prefix"), ConfigValue: types.StringValue(s)}
		resp2 := &frameworkvalidator.StringResponse{}
		vFold.ValidateString(context.Background(), req2, resp2)

		lower := strings.ToLower(s)
		expectFold := strings.HasPrefix(lower, "tf-") || strings.HasPrefix(lower, "iac-")
		if expectFold != !resp2.Diagnostics.HasError() {
			t.Fatalf("case-insensitive mismatch for %q: expect=%v diagErr=%v", s, expectFold, resp2.Diagnostics.HasError())
		}
	})
}
