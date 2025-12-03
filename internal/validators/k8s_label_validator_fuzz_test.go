package validators

import (
	"testing"
)

func FuzzValidateLabelKey(f *testing.F) {
	// Seed corpus with valid and invalid inputs
	f.Add("app")
	f.Add("kubernetes.io/name")
	f.Add("example.com/app-name")
	f.Add("app_name")
	f.Add("")
	f.Add("a/b/c")
	f.Add("app@name")

	f.Fuzz(func(t *testing.T, key string) {
		// Just ensure it doesn't panic
		_ = ValidateLabelKey(key)
	})
}

func FuzzValidateLabelValue(f *testing.F) {
	// Seed corpus
	f.Add("production")
	f.Add("")
	f.Add("prod-env")
	f.Add("v1.0")
	f.Add("app123")
	f.Add("Production")

	f.Fuzz(func(t *testing.T, value string) {
		// Just ensure it doesn't panic
		_ = ValidateLabelValue(value)
	})
}

func FuzzValidateAnnotationValue(f *testing.F) {
	// Seed corpus
	f.Add("")
	f.Add("This is an annotation")
	f.Add("annotation@example.com: value!")
	f.Add("line1\nline2")

	f.Fuzz(func(t *testing.T, value string) {
		// Just ensure it doesn't panic
		_ = ValidateAnnotationValue(value)
	})
}
