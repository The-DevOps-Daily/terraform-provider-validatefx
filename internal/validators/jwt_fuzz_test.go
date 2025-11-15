package validators

import (
	"context"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func FuzzJWTValidator(f *testing.F) {
	// minimal valid-like components
	validSeg := base64.RawURLEncoding.EncodeToString([]byte("{}"))
	good := strings.Join([]string{validSeg, validSeg, validSeg}, ".")
	seeds := []string{"", good, "a.b", "a..b", "a.b.c.d", "not-a-jwt", "a.b.c"}
	for _, s := range seeds {
		f.Add(s)
	}

	v := JWT()
	f.Fuzz(func(t *testing.T, s string) {
		t.Parallel()
		req := frameworkvalidator.StringRequest{Path: path.Root("jwt"), ConfigValue: types.StringValue(s)}
		resp := &frameworkvalidator.StringResponse{}
		v.ValidateString(context.Background(), req, resp)
		if strings.TrimSpace(s) == "" {
			if resp.Diagnostics.HasError() {
				t.Fatalf("empty should not error")
			}
			return
		}

		parts := strings.Split(s, ".")
		expect := len(parts) == 3
		if expect {
			for _, p := range parts {
				if p == "" {
					expect = false
					break
				}
				if _, err := base64.RawURLEncoding.DecodeString(p); err != nil {
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
