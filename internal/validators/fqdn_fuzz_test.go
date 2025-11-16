package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzFQDNValidator(f *testing.F) {
	seeds := []string{"", "app.example.com", "xn--bcher-kva.example", "bad..label", "singlelabel"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := FQDN()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("fqdn"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if strings.TrimSpace(s) == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}

		// basic oracle using internal helpers
		parts := strings.Split(s, ".")
		expect := len(parts) >= 2 && len(s) <= 253
		if expect {
			for _, label := range parts {
				if label == "" || !(fqdnLabelASCII.MatchString(label) || fqdnLabelPuny.MatchString(label)) {
					expect = false
					break
				}
			}
		}
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch for %q: expect=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
		}
	})
}
