package validators

import (
	"context"
	"testing"
)

// TestDescriptionMethods ensures Description and MarkdownDescription methods are covered.
// These methods are required by the validator interface but often return simple strings.
func TestDescriptionMethods(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		validator interface {
			Description(context.Context) string
			MarkdownDescription(context.Context) string
		}
	}{
		{"ARN", ARN()},
		{"Base32", Base32Validator()},
		{"Base64", Base64Validator()},
		{"CIDR", CIDR()},
		{"CreditCard", CreditCard()},
		{"Domain", Domain()},
		{"Email", Email()},
		{"FQDN", FQDN()},
		{"Hex", Hex()},
		{"Hostname", Hostname()},
		{"IP", IP()},
		{"JSON", JSON()},
		{"JWT", JWT()},
		{"MACAddress", MACAddress()},
		{"Phone", Phone()},
		{"PrivateIP", PrivateIP()},
		{"PublicIP", PublicIP()},
		{"SemVer", SemVer()},
		{"SSHPublicKey", SSHPublicKeyValidator()},
		{"StringContains", StringContains([]string{"test"}, true)},
		{"StringLength", NewStringLengthValidator(intPtr(1), intPtr(10))},
		{"StringPrefix", StringPrefix([]string{"test"}, true)},
		{"StringSuffix", StringSuffix("test")},
		{"Subnet", Subnet()},
		{"URI", URI()},
		{"URL", URL()},
		{"Username", DefaultUsernameValidator()},
		{"UUID", UUID()},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			desc := tc.validator.Description(ctx)
			if desc == "" {
				t.Errorf("%s: Description() returned empty string", tc.name)
			}

			mdDesc := tc.validator.MarkdownDescription(ctx)
			if mdDesc == "" {
				t.Errorf("%s: MarkdownDescription() returned empty string", tc.name)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
