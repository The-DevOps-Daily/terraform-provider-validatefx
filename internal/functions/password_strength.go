package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	frameworkvalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
)

type passwordStrengthFunction struct{}

var _ function.Function = (*passwordStrengthFunction)(nil)

func NewPasswordStrengthFunction() function.Function { return &passwordStrengthFunction{} }

func (passwordStrengthFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "password_strength"
}

func (passwordStrengthFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Checks if a password meets strength requirements",
		MarkdownDescription: "Validates that a password has minimum length and contains uppercase, lowercase, number, and special character.",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "password",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "Password string to validate.",
			},
		},
	}
}

func (passwordStrengthFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input types.String
	if err := req.Arguments.GetArgument(ctx, 0, &input); err != nil {
		resp.Error = err
		return
	}

	if input.IsNull() || input.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	v := validators.PasswordStrengthValidator()
	r := frameworkvalidator.StringResponse{}
	v.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("password"), ConfigValue: input}, &r)
	if r.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, r.Diagnostics)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
