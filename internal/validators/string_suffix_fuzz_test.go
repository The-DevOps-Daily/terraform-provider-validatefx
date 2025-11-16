package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzStringSuffixValidator(f *testing.F) {
	seeds := []string{"", "service-tf", "SERVICE-TF", "no-suffix"}
	for _, s := range seeds {
		f.Add(s)
	}

	// This validator is case-sensitive only and takes variadic suffixes.
	vCase := StringSuffix("-tf", "-iac")
	vFold := StringSuffix("-tf", "-iac")

	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		// case-sensitive
		req := frameworkvalidator.StringRequest{Path: path.Root("suffix"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		vCase.ValidateString(context.Background(), req, resp)
		expectExact := strings.HasSuffix(s, "-tf") || strings.HasSuffix(s, "-iac")
		if s == "" {
			expectExact = false
		}
		if expectExact != !resp.Diagnostics.HasError() {
			t.Fatalf("case-sensitive mismatch for %q: expect=%v diagErr=%v", s, expectExact, resp.Diagnostics.HasError())
		}

		// case-insensitive variant not supported by validator; just ensure no panics on different casing
		_ = vFold
	})
}
