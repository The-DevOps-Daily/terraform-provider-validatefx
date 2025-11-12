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

type jwtFunction struct{}

var _ function.Function = (*jwtFunction)(nil)

func NewJWTFunction() function.Function { return &jwtFunction{} }

func (jwtFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "jwt"
}

func (jwtFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string is a well-formed JSON Web Token (JWT).",
		MarkdownDescription: "Returns true when the input string is a well-formed JWT (three base64url segments).",
		Return:              function.BoolReturn{},
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:               "value",
				AllowNullValue:     true,
				AllowUnknownValues: true,
				Description:        "String value to validate.",
			},
		},
	}
}

func (jwtFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input types.String
	if err := req.Arguments.GetArgument(ctx, 0, &input); err != nil {
		resp.Error = err
		return
	}

	if input.IsNull() || input.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	v := validators.JWT()
	r := frameworkvalidator.StringResponse{}
	v.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: input}, &r)
	if r.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, r.Diagnostics)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
