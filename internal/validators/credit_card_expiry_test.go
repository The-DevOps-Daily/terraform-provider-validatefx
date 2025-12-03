package validators

import (
	"context"
	"fmt"
	"testing"
	"time"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestCreditCardExpiryValidator(t *testing.T) {
	ctx := context.Background()

	// Get current date for testing
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	// Future dates
	nextMonth := currentMonth + 1
	nextYear := currentYear
	if nextMonth > 12 {
		nextMonth = 1
		nextYear++
	}

	// Past dates
	pastYear := currentYear - 1
	lastYear2Digit := pastYear % 100

	testCases := []struct {
		name        string
		input       string
		expectError bool
	}{
		// Valid formats - future dates
		{
			name:        "Valid MM/YY format - next month",
			input:       fmt.Sprintf("%02d/%02d", nextMonth, nextYear%100),
			expectError: false,
		},
		{
			name:        "Valid MM/YYYY format - next year",
			input:       fmt.Sprintf("%02d/%d", currentMonth, currentYear+1),
			expectError: false,
		},
		{
			name:        "Valid - far future 2-digit year",
			input:       "12/99",
			expectError: false,
		},
		{
			name:        "Valid - far future 4-digit year",
			input:       "12/2099",
			expectError: false,
		},
		{
			name:        "Valid - January future year",
			input:       fmt.Sprintf("01/%d", currentYear+2),
			expectError: false,
		},
		{
			name:        "Valid - December future year",
			input:       fmt.Sprintf("12/%d", currentYear+1),
			expectError: false,
		},

		// Invalid formats
		{
			name:        "Invalid format - single digit month",
			input:       "1/25",
			expectError: true,
		},
		{
			name:        "Invalid format - single digit year",
			input:       "01/5",
			expectError: true,
		},
		{
			name:        "Invalid format - 3-digit year",
			input:       "01/025",
			expectError: true,
		},
		{
			name:        "Invalid format - no slash",
			input:       "0125",
			expectError: true,
		},
		{
			name:        "Invalid format - dash separator",
			input:       "01-25",
			expectError: true,
		},
		{
			name:        "Invalid format - space separator",
			input:       "01 25",
			expectError: true,
		},
		{
			name:        "Invalid format - letters",
			input:       "AB/CD",
			expectError: true,
		},

		// Invalid month values
		{
			name:        "Invalid month - 00",
			input:       "00/25",
			expectError: true,
		},
		{
			name:        "Invalid month - 13",
			input:       "13/25",
			expectError: true,
		},
		{
			name:        "Invalid month - 99",
			input:       "99/25",
			expectError: true,
		},

		// Past dates
		{
			name:        "Past date - last year 2-digit",
			input:       fmt.Sprintf("12/%02d", lastYear2Digit),
			expectError: true,
		},
		{
			name:        "Past date - last year 4-digit",
			input:       fmt.Sprintf("12/%d", pastYear),
			expectError: true,
		},
		{
			name:        "Past date - year 2000",
			input:       "12/2000",
			expectError: true,
		},
		{
			name:        "Past date - year 2020",
			input:       "01/2020",
			expectError: true,
		},

		// Edge cases
		{
			name:        "Empty string",
			input:       "",
			expectError: false, // Empty strings are allowed (handled upstream)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator := CreditCardExpiry()

			req := frameworkvalidator.StringRequest{
				ConfigValue: types.StringValue(tc.input),
			}
			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(ctx, req, resp)

			hasErrors := resp.Diagnostics.HasError()
			if hasErrors != tc.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tc.expectError, hasErrors, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestCreditCardExpiryValidator_NullAndUnknown(t *testing.T) {
	ctx := context.Background()
	validator := CreditCardExpiry()

	testCases := []struct {
		name  string
		value types.String
	}{
		{
			name:  "Null value",
			value: types.StringNull(),
		},
		{
			name:  "Unknown value",
			value: types.StringUnknown(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := frameworkvalidator.StringRequest{
				ConfigValue: tc.value,
			}
			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(ctx, req, resp)

			if resp.Diagnostics.HasError() {
				t.Errorf("Expected no errors for %s, got: %v", tc.name, resp.Diagnostics.Errors())
			}
		})
	}
}

func TestCreditCardExpiryValidator_Description(t *testing.T) {
	ctx := context.Background()
	validator := CreditCardExpiry()

	desc := validator.Description(ctx)
	if desc == "" {
		t.Error("Description should not be empty")
	}

	mdDesc := validator.MarkdownDescription(ctx)
	if mdDesc == "" {
		t.Error("MarkdownDescription should not be empty")
	}

	if desc != mdDesc {
		t.Error("Description and MarkdownDescription should match")
	}
}

// Test the edge case where we're in the expiry month
func TestCreditCardExpiryValidator_CurrentMonth(t *testing.T) {
	ctx := context.Background()
	validator := CreditCardExpiry()

	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	// Current month/year should be valid (expires at end of month)
	testCases := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "Current month - 2 digit year",
			input:       fmt.Sprintf("%02d/%02d", currentMonth, currentYear%100),
			expectError: false,
		},
		{
			name:        "Current month - 4 digit year",
			input:       fmt.Sprintf("%02d/%d", currentMonth, currentYear),
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := frameworkvalidator.StringRequest{
				ConfigValue: types.StringValue(tc.input),
			}
			resp := &frameworkvalidator.StringResponse{}

			validator.ValidateString(ctx, req, resp)

			hasErrors := resp.Diagnostics.HasError()
			if hasErrors != tc.expectError {
				t.Errorf("Expected error: %v, got error: %v. Diagnostics: %v",
					tc.expectError, hasErrors, resp.Diagnostics.Errors())
			}
		})
	}
}
