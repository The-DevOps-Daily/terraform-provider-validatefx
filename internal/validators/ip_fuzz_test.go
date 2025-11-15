package validators

import (
    "context"
    "net"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/path"
    frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// FuzzIPValidator checks robustness and aligns with net.ParseIP for basic validity.
func FuzzIPValidator(f *testing.F) {
    for _, s := range []string{"", "127.0.0.1", "::1", "256.0.0.1", "2001:db8::1", "bad ip"} { f.Add(s) }
    v := IP()
    f.Fuzz(func(t *testing.T, s string) {
        t.Parallel()
        req := frameworkvalidator.StringRequest{ Path: path.Root("ip"), ConfigValue: types.StringValue(s) }
        resp := &frameworkvalidator.StringResponse{}
        v.ValidateString(context.Background(), req, resp)
        if s == "" { if resp.Diagnostics.HasError() { t.Fatalf("empty should not error") }; return }
        expect := net.ParseIP(s) != nil
        if expect != !resp.Diagnostics.HasError() {
            t.Fatalf("mismatch for %q: parse-ok=%v diagErr=%v", s, expect, resp.Diagnostics.HasError())
        }
    })
}

