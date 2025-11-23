package validators

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ frameworkvalidator.String = CreditCard()

// CreditCard returns a schema.String validator which validates credit card numbers using the Luhn algorithm.
func CreditCard() frameworkvalidator.String {
	return creditCardValidator{}
}

type creditCardValidator struct{}

func (creditCardValidator) Description(_ context.Context) string {
	return "value must be a valid credit card number (Luhn algorithm)"
}

func (creditCardValidator) MarkdownDescription(_ context.Context) string {
	return "value must be a valid credit card number"
}

func (creditCardValidator) ValidateString(_ context.Context, req frameworkvalidator.StringRequest, resp *frameworkvalidator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value == "" {
		return
	}

	if !isValidCreditCard(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Credit Card Number",
			fmt.Sprintf("Value %q is not a valid credit card number according to the Luhn algorithm", value),
		)
	}
}

// isValidCreditCard validates a credit card number using the Luhn algorithm
func isValidCreditCard(cardNumber string) bool {
	// Remove spaces and hyphens
	cardNumber = strings.ReplaceAll(cardNumber, " ", "")
	cardNumber = strings.ReplaceAll(cardNumber, "-", "")

	// Check if all characters are digits
	for _, char := range cardNumber {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Credit card numbers should be between 13-19 digits
	length := len(cardNumber)
	if length < 13 || length > 19 {
		return false
	}

	// Reject all zeros (invalid in practice)
	if strings.Trim(cardNumber, "0") == "" {
		return false
	}

	return luhnCheck(cardNumber)
}

// luhnCheck implements the Luhn algorithm for credit card validation
func luhnCheck(cardNumber string) bool {
	var sum int
	alternate := false

	// Process digits from right to left
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(cardNumber[i]))
		if err != nil {
			return false
		}

		if alternate {
			digit *= 2
			if digit > 9 {
				digit = digit%10 + digit/10
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}
