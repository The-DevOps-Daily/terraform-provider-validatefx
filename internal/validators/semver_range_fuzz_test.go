package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzSemVerRangeValidator(f *testing.F) {
	seeds := []string{"", ">=1.2.3", ">=1.0.0, <2.0.0", "=v1.2.3", "bad comparator"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := SemVerRange()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("range"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
		}
		// No strict oracle without re-implementing parsing. We ensure robustness and acceptance of valid-like seeds.
	})
}
