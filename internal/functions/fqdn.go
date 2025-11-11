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

type fqdnFunction struct{}

var _ function.Function = (*fqdnFunction)(nil)

func NewFQDNFunction() function.Function { return &fqdnFunction{} }

func (fqdnFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "fqdn"
}

func (fqdnFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validate that a string is a fully qualified domain name (FQDN).",
		MarkdownDescription: "Returns true when the input string is a valid FQDN.",
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

func (fqdnFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input types.String
	if err := req.Arguments.GetArgument(ctx, 0, &input); err != nil {
		resp.Error = err
		return
	}

	if input.IsNull() || input.IsUnknown() {
		resp.Result = function.NewResultData(types.BoolUnknown())
		return
	}

	v := validators.FQDN()
	r := frameworkvalidator.StringResponse{}
	v.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: input}, &r)
	if r.Diagnostics.HasError() {
		resp.Error = function.FuncErrorFromDiags(ctx, r.Diagnostics)
		return
	}
	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
