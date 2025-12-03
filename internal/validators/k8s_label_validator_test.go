package validators

import (
	"strings"
	"testing"
)

func TestValidateLabelKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid simple key",
			key:     "app",
			wantErr: false,
		},
		{
			name:    "valid key with prefix",
			key:     "example.com/app",
			wantErr: false,
		},
		{
			name:    "valid key with subdomain prefix",
			key:     "kubernetes.io/app",
			wantErr: false,
		},
		{
			name:    "valid key with dashes",
			key:     "app-name",
			wantErr: false,
		},
		{
			name:    "valid key with underscores",
			key:     "app_name",
			wantErr: false,
		},
		{
			name:    "valid key with dots",
			key:     "app.name",
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     "",
			wantErr: true,
		},
		{
			name:    "too many slashes",
			key:     "a/b/c",
			wantErr: true,
		},
		{
			name:    "prefix too long",
			key:     strings.Repeat("a", 254) + "/name",
			wantErr: true,
		},
		{
			name:    "name too long",
			key:     strings.Repeat("a", 64),
			wantErr: true,
		},
		{
			name:    "invalid prefix format uppercase",
			key:     "Example.com/app",
			wantErr: true,
		},
		{
			name:    "invalid prefix with special chars",
			key:     "example.com$/app",
			wantErr: true,
		},
		{
			name:    "empty name with prefix",
			key:     "example.com/",
			wantErr: true,
		},
		{
			name:    "prefix starts with dash",
			key:     "-example.com/app",
			wantErr: true,
		},
		{
			name:    "prefix ends with dash",
			key:     "example.com-/app",
			wantErr: true,
		},
		{
			name:    "name with special chars",
			key:     "app@name",
			wantErr: true,
		},
		{
			name:    "name starts with dash",
			key:     "-appname",
			wantErr: true,
		},
		{
			name:    "name ends with dash",
			key:     "appname-",
			wantErr: true,
		},
		{
			name:    "max length name",
			key:     strings.Repeat("a", 63),
			wantErr: false,
		},
		{
			name:    "max length prefix",
			key:     strings.Repeat("a", 253) + "/app",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLabelKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLabelKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLabelValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid simple value",
			value:   "production",
			wantErr: false,
		},
		{
			name:    "valid value with dash",
			value:   "prod-env",
			wantErr: false,
		},
		{
			name:    "valid value with underscore",
			value:   "prod_env",
			wantErr: false,
		},
		{
			name:    "valid value with dot",
			value:   "v1.0",
			wantErr: false,
		},
		{
			name:    "empty value is valid",
			value:   "",
			wantErr: false,
		},
		{
			name:    "value with numbers",
			value:   "app123",
			wantErr: false,
		},
		{
			name:    "value too long",
			value:   strings.Repeat("a", 64),
			wantErr: true,
		},
		{
			name:    "max length value",
			value:   strings.Repeat("a", 63),
			wantErr: false,
		},
		{
			name:    "value with special chars",
			value:   "prod@env",
			wantErr: true,
		},
		{
			name:    "value starts with dash",
			value:   "-production",
			wantErr: true,
		},
		{
			name:    "value ends with dash",
			value:   "production-",
			wantErr: true,
		},
		{
			name:    "value with uppercase",
			value:   "Production",
			wantErr: false,
		},
		{
			name:    "value with mixed case",
			value:   "ProdEnv",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLabelValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLabelValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAnnotationValue(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "empty annotation",
			value:   "",
			wantErr: false,
		},
		{
			name:    "simple annotation",
			value:   "This is a simple annotation",
			wantErr: false,
		},
		{
			name:    "annotation with special chars",
			value:   "annotation@example.com: value!",
			wantErr: false,
		},
		{
			name:    "annotation with newlines",
			value:   "line1\nline2\nline3",
			wantErr: false,
		},
		{
			name:    "large annotation",
			value:   strings.Repeat("a", 100000),
			wantErr: false,
		},
		{
			name:    "max size annotation",
			value:   strings.Repeat("a", 262144),
			wantErr: false,
		},
		{
			name:    "annotation too large",
			value:   strings.Repeat("a", 262145),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAnnotationValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAnnotationValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
