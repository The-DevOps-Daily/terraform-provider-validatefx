package functions

import (
	"context"

	"github.com/The-DevOps-Daily/terraform-provider-validatefx/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ function.Function = &K8sAnnotationValueFunction{}

// K8sAnnotationValueFunction validates Kubernetes annotation values.
type K8sAnnotationValueFunction struct{}

func NewK8sAnnotationValueFunction() function.Function {
	return &K8sAnnotationValueFunction{}
}

func (f *K8sAnnotationValueFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "k8s_annotation_value"
}

func (f *K8sAnnotationValueFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Validates Kubernetes annotation value format",
		Description:         "Validates that a string is a valid Kubernetes annotation value (up to 256KB).",
		MarkdownDescription: "Validates that a string is a valid Kubernetes annotation value (up to 256KB).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "value",
				Description: "The annotation value to validate",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f *K8sAnnotationValueFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var value string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &value))
	if resp.Error != nil {
		return
	}

	if err := validators.ValidateAnnotationValue(value); err != nil {
		resp.Error = function.NewArgumentFuncError(0, "Invalid Kubernetes Annotation Value: "+err.Error())
		return
	}

	resp.Result = function.NewResultData(basetypes.NewBoolValue(true))
}
