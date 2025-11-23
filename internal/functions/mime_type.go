package functions

import (
	"github.com/hashicorp/terraform-plugin-framework/function"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewMIMETypeFunction exposes the MIME type validator as a Terraform function.
func NewMIMETypeFunction() function.Function {
	return newStringValidationFunction(
		"mime_type",
		"Validate that a string is a valid MIME type.",
		"Returns true when the input is a valid MIME type (e.g. application/json, text/html).",
		validators.MIMEType(),
	)
}
