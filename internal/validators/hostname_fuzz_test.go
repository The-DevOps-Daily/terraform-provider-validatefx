package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzHostnameValidator ensures the hostname validator is robust and broadly
// consistent with RFC1123 constraints enforced by isRFC1123Hostname.
func FuzzHostnameValidator(f *testing.F) {
	seeds := []string{"", "example", "example.com", "-bad", "toolonglabeltoolonglabeltoolonglabeltoolonglabeltoolonglabeltool", "good-host", "bad..dots", "xn--bcher-kva.example"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := Hostname()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("hostname"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}

		expect := isRFC1123Hostname(s)
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: expectValid=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
