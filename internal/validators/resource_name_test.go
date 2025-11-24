package validators

import (
	"context"
	"testing"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestResourceName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		value       string
		expectError bool
	}{
		// Valid resource names
		{
			name:        "simple lowercase",
			value:       "myresource",
			expectError: false,
		},
		{
			name:        "with underscores",
			value:       "my_resource_name",
			expectError: false,
		},
		{
			name:        "with hyphens",
			value:       "my-resource-name",
			expectError: false,
		},
		{
			name:        "with digits",
			value:       "resource123",
			expectError: false,
		},
		{
			name:        "starts with underscore",
			value:       "_private_resource",
			expectError: false,
		},
		{
			name:        "mixed valid characters",
			value:       "aws_s3_bucket_2024",
			expectError: false,
		},
		{
			name:        "single letter",
			value:       "a",
			expectError: false,
		},
		{
			name:        "single underscore",
			value:       "_",
			expectError: false,
		},

		// Invalid resource names
		{
			name:        "empty string",
			value:       "",
			expectError: true,
		},
		{
			name:        "uppercase letters",
			value:       "MyResource",
			expectError: true,
		},
		{
			name:        "starts with digit",
			value:       "1resource",
			expectError: true,
		},
		{
			name:        "starts with hyphen",
			value:       "-resource",
			expectError: true,
		},
		{
			name:        "contains spaces",
			value:       "my resource",
			expectError: true,
		},
		{
			name:        "contains dots",
			value:       "my.resource",
			expectError: true,
		},
		{
			name:        "contains special characters",
			value:       "my@resource",
			expectError: true,
		},
		{
			name:        "all uppercase",
			value:       "MYRESOURCE",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator := ResourceName()
			resp := &frameworkvalidator.StringResponse{}
			validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
				ConfigValue: types.StringValue(tc.value),
			}, resp)

			if tc.expectError && !resp.Diagnostics.HasError() {
				t.Fatalf("expected error for %q, got none", tc.value)
			}

			if !tc.expectError && resp.Diagnostics.HasError() {
				t.Fatalf("expected no error for %q, got: %v", tc.value, resp.Diagnostics)
			}
		})
	}
}

func TestResourceNameNullAndUnknown(t *testing.T) {
	t.Parallel()

	validator := ResourceName()

	// Test null value
	resp := &frameworkvalidator.StringResponse{}
	validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
		ConfigValue: types.StringNull(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatal("expected no error for null value")
	}

	// Test unknown value
	resp = &frameworkvalidator.StringResponse{}
	validator.ValidateString(context.Background(), frameworkvalidator.StringRequest{
		ConfigValue: types.StringUnknown(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatal("expected no error for unknown value")
	}
}

func TestResourceNameDescription(t *testing.T) {
	t.Parallel()

	validator := ResourceName()
	desc := validator.Description(context.Background())

	if desc == "" {
		t.Fatal("expected non-empty description")
	}

	markdownDesc := validator.MarkdownDescription(context.Background())
	if markdownDesc == "" {
		t.Fatal("expected non-empty markdown description")
	}
}
