package functions

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// ProviderFunctionFactories returns all Terraform function constructors exposed by the provider.
func ProviderFunctionFactories() []func() function.Function {
	return []func() function.Function{
		NewAssertFunction,
		NewEmailFunction,
		NewUUIDFunction,
		NewUUIDv4OnlyFunction,
		NewBase64Function,
		NewBase32Function,
		NewCreditCardFunction,
		NewDomainFunction,
		NewHostnameFunction,
		NewDateTimeFunction,
		NewJSONFunction,
		NewSemVerFunction,
		NewSemVerRangeFunction,
		NewHexFunction,
		NewIntegerFunction,
		NewSSHPublicKeyFunction,
		NewIPFunction,
		NewMatchesRegexFunction,
		NewStringLengthFunction,
		NewPhoneFunction,
		NewMACAddressFunction,
		NewMIMETypeFunction,
		NewSlugFunction,
		NewPositiveNumberFunction,
		NewNonNegativeNumberFunction,
		NewSizeBetweenFunction,
		NewMapKeysMatchFunction,
		NewMutuallyExclusiveFunction,
		NewNonEmptyListFunction,
		NewResourceNameFunction,
		NewURLFunction,
		NewBetweenFunction,
		NewInListFunction,
		NewNotInListFunction,
		NewStringContainsFunction,
		NewStringSuffixFunction,
		NewStringPrefixFunction,
		NewUsernameFunction,
		NewSetEqualsFunction,
		NewAllValidFunction,
		NewAnyValidFunction,
		NewExactlyOneValidFunction,
		NewVersionFunction,
		NewCIDRFunction,
		NewPasswordStrengthFunction,
		NewFQDNFunction,
		NewJWTFunction,
		NewCIDROverlapFunction,
		NewPortRangeFunction,
		NewPrivateIPFunction,
		NewURIFunction,
		NewListSubsetFunction,
		NewPortNumberFunction,
		NewSubnetFunction,
		NewPublicIPFunction,
		NewARNFunction,
		NewAWSRegionFunction,
		NewGCPRegionFunction,
		NewGCPZoneFunction,
		NewAzureLocationFunction,
		NewIPRangeSizeFunction,
		NewListLengthBetweenFunction,
	}
}

// FunctionDoc captures high level documentation details for a Terraform function.
type FunctionDoc struct {
	Name        string
	Summary     string
	Description string
}

// AvailableFunctionDocs returns documentation metadata for every exported Terraform function.
func AvailableFunctionDocs(ctx context.Context) ([]FunctionDoc, error) {
	factories := ProviderFunctionFactories()

	docs := make([]FunctionDoc, 0, len(factories))

	for _, factory := range factories {
		fn := factory()

		metaResp := &function.MetadataResponse{}
		fn.Metadata(ctx, function.MetadataRequest{}, metaResp)

		defResp := &function.DefinitionResponse{}
		fn.Definition(ctx, function.DefinitionRequest{}, defResp)

		summary := strings.TrimSpace(defResp.Definition.Summary)
		description := strings.TrimSpace(defResp.Definition.MarkdownDescription)
		if summary == "" {
			summary = description
		}

		docs = append(docs, FunctionDoc{
			Name:        metaResp.Name,
			Summary:     summary,
			Description: description,
		})
	}

	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Name < docs[j].Name
	})

	return docs, nil
}
