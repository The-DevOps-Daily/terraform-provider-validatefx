package validators

import (
	"context"
	"net/url"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzURLValidator ensures the URL() validator remains robust for arbitrary inputs
// and aligns with url.ParseRequestURI expectations for basic validity (scheme+host present).
func FuzzURLValidator(f *testing.F) {
	seeds := []string{
		"", "http://example.com", "https://example.com/path?x=1#frag", "ftp://example.com", "://bad", "http:/bad", "https://", "http://",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	v := URL()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("url"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)

		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error: %v", resp.Diagnostics)
			}
			return
		}

		parsed, err := url.ParseRequestURI(s)
		basicValid := err == nil && parsed != nil && parsed.Scheme != "" && parsed.Host != ""

		// Our validator only permits http/https schemes.
		if basicValid && parsed.Scheme != "http" && parsed.Scheme != "https" {
			basicValid = false
		}

		if basicValid != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch: basicValid=%v diagnostics.HasError=%v value=%q", basicValid, resp.Diagnostics.HasError(), s)
		}
	})
}
