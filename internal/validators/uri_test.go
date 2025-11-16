package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestURIValidator(t *testing.T) {
	t.Parallel()
	v := URI()

	run := func(s string) *frameworkvalidator.StringResponse {
		req := frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		return resp
	}

	valids := []string{
		"http://example.com",
		"https://example.com/path?x=1#frag",
		"ftp://ftp.example.com/resource",
		"ssh://user@host",
		"postgres://user:pass@localhost:5432/db?sslmode=disable",
		"mysql://user@127.0.0.1/db",
		"amqp://guest:guest@localhost:5672/",
		"urn:isbn:0451450523", // non-hierarchical allowed
	}
	for _, s := range valids {
		if resp := run(s); resp.Diagnostics.HasError() {
			t.Fatalf("expected valid URI %q, got %v", s, resp.Diagnostics)
		}
	}

	invalids := []string{
		"example.com",  // missing scheme
		"http:///path", // missing host
		"ssh://",       // missing host
		"://host",      // bad scheme
		"http://",      // missing host
		"",             // empty
		"   ",          // spaces
	}
	for _, s := range invalids {
		if resp := run(s); !resp.Diagnostics.HasError() {
			t.Fatalf("expected invalid URI %q", s)
		}
	}
}
