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

type semverRangeFunction struct{}

var _ function.Function = (*semverRangeFunction)(nil)

func NewSemVerRangeFunction() function.Function { return &semverRangeFunction{} }

func (semverRangeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
    resp.Name = "semver_range"
}

func (semverRangeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
    resp.Definition = function.Definition{
        Summary:             "Validate that a string is a valid semantic version range expression.",
        MarkdownDescription: "Returns true when the input string represents a valid SemVer range, e.g., '>=1.0.0,<2.0.0'.",
        Return:              function.BoolReturn{},
        Parameters: []function.Parameter{
            function.StringParameter{
                Name:               "value",
                AllowNullValue:     true,
                AllowUnknownValues: true,
                Description:        "Version range expression to validate.",
            },
        },
    }
}

func (semverRangeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
    var input types.String
    if err := req.Arguments.GetArgument(ctx, 0, &input); err != nil {
        resp.Error = err
        return
    }

    if input.IsNull() || input.IsUnknown() {
        resp.Result = function.NewResultData(types.BoolUnknown())
        return
    }

    v := validators.SemVerRange()
    r := frameworkvalidator.StringResponse{}
    v.ValidateString(ctx, frameworkvalidator.StringRequest{Path: path.Root("value"), ConfigValue: input}, &r)
    if r.Diagnostics.HasError() {
        resp.Error = function.FuncErrorFromDiags(ctx, r.Diagnostics)
        return
    }
    resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}

