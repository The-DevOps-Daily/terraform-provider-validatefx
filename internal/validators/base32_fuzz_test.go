package validators

import (
	"context"
	"encoding/base32"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzBase32Validator(f *testing.F) {
	seeds := []string{"", "MZXW6===", "NBSWY3DP", "INVALID*BASE32", "GEZDGNBV"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := Base32Validator()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("b32"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if s == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}
		_, err := base32.StdEncoding.DecodeString(s)
		expect := err == nil
		if expect != !resp.Diagnostics.HasError() {
			t.Fatalf("mismatch: decode-ok=%v diagErr=%v for %q", expect, resp.Diagnostics.HasError(), s)
		}
	})
}
