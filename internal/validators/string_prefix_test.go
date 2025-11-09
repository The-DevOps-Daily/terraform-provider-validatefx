package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringPrefixValidatorValid(t *testing.T) {
	t.Parallel()

	validator := StringPrefix([]string{"tf-", "iac-"}, false)
	req := frameworkvalidator.StringRequest{
		Path:        path.Root("prefix"),
		ConfigValue: types.StringValue("tf-project"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics for valid prefix: %v", resp.Diagnostics)
	}
}

func TestStringPrefixValidatorInvalid(t *testing.T) {
	t.Parallel()

	validator := StringPrefix([]string{"tf-", "iac-"}, false)
	req := frameworkvalidator.StringRequest{
		Path:        path.Root("prefix"),
		ConfigValue: types.StringValue("dev-project"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Fatalf("expected diagnostics for invalid prefix")
	}

	d := resp.Diagnostics[0]
	if d.Severity() != diag.SeverityError {
		t.Fatalf("expected severity error, got %s", d.Severity())
	}

	if d.Summary() != "Invalid Prefix" {
		t.Fatalf("unexpected summary: %s", d.Summary())
	}
}

func TestStringPrefixValidatorIgnoreCase(t *testing.T) {
	t.Parallel()

	validator := StringPrefix([]string{"tf-"}, true)
	req := frameworkvalidator.StringRequest{
		Path:        path.Root("prefix"),
		ConfigValue: types.StringValue("TF-environment"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics when ignoring case, got: %v", resp.Diagnostics)
	}
}

func TestStringPrefixValidatorHandlesUnknownAndNull(t *testing.T) {
	t.Parallel()

	validator := StringPrefix([]string{"tf-"}, false)
	requests := map[string]frameworkvalidator.StringRequest{
		"unknown": {
			Path:        path.Root("prefix"),
			ConfigValue: types.StringUnknown(),
		},
		"null": {
			Path:        path.Root("prefix"),
			ConfigValue: types.StringNull(),
		},
	}

	for name, req := range requests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &frameworkvalidator.StringResponse{}
			validator.ValidateString(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("expected no diagnostics for %s input", name)
			}
		})
	}
}

func TestStringPrefixValidatorNoPrefixes(t *testing.T) {
	t.Parallel()

	validator := StringPrefix(nil, false)
	req := frameworkvalidator.StringRequest{
		Path:        path.Root("prefix"),
		ConfigValue: types.StringValue("any"),
	}
	resp := &frameworkvalidator.StringResponse{}

	validator.ValidateString(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no diagnostics when no prefixes configured")
	}
}
