package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	dnsSubdomainFmt  = regexp.MustCompile(`^([a-z0-9]([-a-z0-9]*[a-z0-9])?)(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)
	qualifiedNameFmt = regexp.MustCompile(`^[A-Za-z0-9]([A-Za-z0-9_.-]*[A-Za-z0-9])?$`)
)

// These are default values, you can adjust them as you want
const (
	MaxLabelValueLength      = 63
	MaxLabelNameLength       = 63
	MaxPrefixLength          = 253
	MaxAnnotationValueLength = 262144
)

func validateKey(key string) error {
	if key == "" {
		return errors.New("Key not valid")
	}

	parts := strings.Split(key, "/")
	if len(parts) > 2 {
		return errors.New("Key must contain at most one '/'")
	}

	var prefix, name string

	if len(parts) == 2 {
		prefix = parts[0]
		name = parts[1]

		if len(prefix) > MaxPrefixLength {
			return fmt.Errorf("Prefix exceeds %d characters", MaxPrefixLength)
		}

		if !dnsSubdomainFmt.MatchString(prefix) {
			return errors.New("DNS subdomain not valid")
		}
	} else {
		name = parts[0]
	}

	if len(name) == 0 {
		return errors.New("No name is specified")
	}
	if len(name) > MaxLabelNameLength {
		return fmt.Errorf("Name exceeds %d characters", MaxLabelNameLength)
	}
	if !qualifiedNameFmt.MatchString(name) {
		return fmt.Errorf("Name must match regex: %s", qualifiedNameFmt.String())
	}

	return nil
}

func ValidateLabelValue(value string) error {
	if len(value) > MaxLabelValueLength {
		return fmt.Errorf("Label value exceeds %d characters", MaxLabelValueLength)
	}

	if value == "" {
		return nil
	}

	if !qualifiedNameFmt.MatchString(value) {
		return errors.New("Label value must match the name regex")
	}

	return nil
}

func ValidateAnnotationValue(value string) error {
	if len(value) > MaxAnnotationValueLength {
		return fmt.Errorf("Annotation value exceeds %d bytes", MaxAnnotationValueLength)
	}
	return nil
}
