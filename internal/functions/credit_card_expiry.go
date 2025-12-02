package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewCreditCardExpiryFunction exposes the credit card expiry validator as a Terraform function.
func NewCreditCardExpiryFunction() function.Function {
	return newStringValidationFunction(
		"credit_card_expiry",
		"Validate that a string is a valid credit card expiry date in MM/YY or MM/YYYY format and not in the past.",
		"Returns true when the input is a valid expiry date format with a valid month (01-12) and the date is not in the past.",
		validators.CreditCardExpiry(),
	)
}
