package validators

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"time"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = CreditCardExpiry()

// CreditCardExpiry returns a schema.String validator which validates credit card expiry dates.
// It supports both MM/YY and MM/YYYY formats and ensures the date is not in the past.
func CreditCardExpiry() frameworkvalidator.String {
	return creditCardExpiryValidator{}
}

type creditCardExpiryValidator struct{}

func (creditCardExpiryValidator) Description(_ context.Context) string {
	return "value must be a valid credit card expiry date in MM/YY or MM/YYYY format and not in the past"
}

func (v creditCardExpiryValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (creditCardExpiryValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	if err := validateCreditCardExpiry(value); err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Credit Card Expiry Date",
			err.Error(),
		)
	}
}

// validateCreditCardExpiry validates a credit card expiry date in MM/YY or MM/YYYY format
func validateCreditCardExpiry(expiry string) error {
	// Pattern for MM/YY or MM/YYYY
	pattern := regexp.MustCompile(`^(0[1-9]|1[0-2])/(\d{2}|\d{4})$`)
	matches := pattern.FindStringSubmatch(expiry)

	if matches == nil {
		return fmt.Errorf("value %q is not in valid format (expected MM/YY or MM/YYYY)", expiry)
	}

	month, err := strconv.Atoi(matches[1])
	if err != nil {
		return fmt.Errorf("invalid month in %q", expiry)
	}

	// Month validation already handled by regex (01-12)
	if month < 1 || month > 12 {
		return fmt.Errorf("month must be between 01 and 12, got %02d", month)
	}

	yearStr := matches[2]
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return fmt.Errorf("invalid year in %q", expiry)
	}

	// Convert 2-digit year to 4-digit year
	if len(yearStr) == 2 {
		currentYear := time.Now().Year()
		century := (currentYear / 100) * 100
		year = century + year

		// Handle century rollover
		// If the 2-digit year is less than current year's last 2 digits minus some threshold,
		// assume it's the next century
		if year < currentYear-10 {
			year += 100
		}
	}

	// Create expiry date as the last day of the expiry month
	expiryDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, 0).Add(-time.Second)

	// Compare with current date
	now := time.Now().UTC()
	if expiryDate.Before(now) {
		return fmt.Errorf("credit card expiry date %q is in the past", expiry)
	}

	return nil
}
