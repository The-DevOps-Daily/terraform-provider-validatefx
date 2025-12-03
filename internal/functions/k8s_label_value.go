package functions

import (
	"context"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ function.Function = &K8sLabelValueFunction{}

// K8sLabelValueFunction validates Kubernetes label values.
type K8sLabelValueFunction struct{}

func NewK8sLabelValueFunction() function.Function {
	return &K8sLabelValueFunction{}
}

func (f *K8sLabelValueFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "k8s_label_value"
}

func (f *K8sLabelValueFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validates Kubernetes label value format",
		Description:         "Validates that a string is a valid Kubernetes label value (63 characters max, alphanumeric with dots, dashes, and underscores).",
		MarkdownDescription: "Validates that a string is a valid Kubernetes label value (63 characters max, alphanumeric with dots, dashes, and underscores).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "value",
				Description: "The label value to validate",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f *K8sLabelValueFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &value))
	if resp.Error != nil {
		return
	}

	if err := validators.ValidateLabelValue(value); err != nil {
		resp.Error = function.NewArgumentFuncError(0, "Invalid Kubernetes Label Value: "+err.Error())
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
