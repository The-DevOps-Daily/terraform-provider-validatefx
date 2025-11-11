package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

// NewPasswordStrengthFunction returns a Terraform function
func NewPasswordStrengthFunction() function.Function {
	return function.New(&passwordStrengthFunction{})
}

type passwordStrengthFunction struct{}

func (f *passwordStrengthFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "password_strength"
}

func (f *passwordStrengthFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Checks if a password meets strength requirements",
		MarkdownDescription: "Validates that a password has minimum length and contains uppercase, lowercase, number, and special character.",
		Arguments: []function.Argument{
			{
				Name:     "password",
				Type:     function.TypeString,
				Required: true,
			},
		},
		Returns: []function.Return{
			{
				Type: function.TypeBool,
			},
		},
	}
}

func (f *passwordStrengthFunction) Call(ctx context.Context, req function.CallRequest, resp *function.CallResponse) {
	password, ok := req.Args[0].(string)
	if !ok {
		resp.Diagnostics.AddError("Invalid Argument", "Expected a string password")
		return
	}

	err := validators.PasswordStrength(password)
	if err != nil {
		resp.Responses = []interface{}{false}
		resp.Diagnostics.AddError("Password Validation Failed", err.Error())
		return
	}

	resp.Responses = []interface{}{true}
}
