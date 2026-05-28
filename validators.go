package goflag

import (
	"cmp"
	"fmt"
	"slices"
)

// Choices returns a validator function that checks if the provided value is one of the allowed choices.
func Choices[T comparable](choices []T) func(v any) error {
	return func(v any) error {
		concreteType, ok := v.(T)
		if !ok {
			return fmt.Errorf("invalid generic type for %v", v)
		}
		if slices.Contains(choices, concreteType) {
			return nil
		}
		return fmt.Errorf("expected value to be one of: %v", choices)
	}
}

// MinStringLen returns a validator function that checks if the provided string value has a minimum length.
func MinStringLen(length int) func(v any) error {
	return func(v any) error {
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("MinStringLen must be used only with strings")
		}

		if len(s) < length {
			return fmt.Errorf("string must be at least %d characters long", length)
		}
		return nil
	}
}

// MaxStringLen returns a validator function that checks if the provided string value has a maximum length.
func MaxStringLen(length int) func(v any) error {
	return func(v any) error {
		s, ok := v.(string)
		if !ok {
			return fmt.Errorf("MaxStringLen must be used only with strings")
		}

		if len(s) > length {
			return fmt.Errorf("string must be at most %d characters long", length)
		}
		return nil
	}
}

// Max returns a validator function that checks if the provided value is less than or equal to the specified maximum value.
func Max[T cmp.Ordered](maxValue T) func(v any) error {
	return func(v any) error {
		value := v.(T)
		if value > maxValue {
			return fmt.Errorf("value %v is greater than maximum value: %v", v, maxValue)
		}
		return nil
	}
}

// Min returns a validator function that checks if the provided value is greater than or equal to the specified minimum value.
func Min[T cmp.Ordered](minValue T) func(v any) error {
	return func(v any) error {
		value := v.(T)
		if value < minValue {
			return fmt.Errorf("value %v is less than minimum value: %v", v, minValue)
		}
		return nil
	}
}

// Range returns a validator function that checks if the provided value is within the specified range (inclusive).
func Range[T cmp.Ordered](minValue, maxValue T) func(v any) error {
	return func(v any) error {
		value := v.(T)
		if value < minValue {
			return fmt.Errorf("value %v is less than minimum value: %v", v, minValue)
		}
		if value > maxValue {
			return fmt.Errorf("value %v is greater than maximum value: %v", v, maxValue)
		}
		return nil
	}
}

// NotEmpty returns a validator function that checks if the provided string value is not empty.
func NotEmpty(v any) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("NotEmpty validator can only be used with strings")
	}
	if s == "" {
		return fmt.Errorf("value cannot be empty")
	}
	return nil
}
