package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzDomainValidator(f *testing.F) {
	seeds := []string{"", "example.com", "-bad.example", "good-label.example", "xn--bcher-kva.example", "too..dots"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := Domain()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("domain"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}

		expect := isValidDomain(s)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: expect=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
