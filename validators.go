package goflag

import (
	"cmp"
	"fmt"
	"slices"
)

// Choices returns a validator function that checks if the provided value is one of the allowed choices.
func Choices[T comparable](choices []T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		concreteType, ok := v.(T)
		if !ok {
			return false, fmt.Sprintf("Invalid generic type for %v", v)
		}
		if slices.Contains(choices, concreteType) {
			return true, ""
		}
		return false, fmt.Sprintf("Expected value to be one of: %v", choices)
	}
}

// MinStringLen returns a validator function that checks if the provided string value has a minimum length.
func MinStringLen(length int) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		s, ok := v.(string)
		if !ok {
			return false, "MinStringLen must be used only with strings"
		}

		return len(s) >= length, ""
	}
}

// MaxStringLen returns a validator function that checks if the provided string value has a maximum length.
func MaxStringLen(length int) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		s, ok := v.(string)
		if !ok {
			return false, "MaxStringLen must be used only with strings"
		}

		return len(s) <= length, ""
	}
}

// Max returns a validator function that checks if the provided value is less than or equal to the specified maximum value.
func Max[T cmp.Ordered](maxValue T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		value := v.(T)
		return value <= maxValue, fmt.Sprintf("value %v is greater than maximum value: %v", v, maxValue)
	}
}

// Min returns a validator function that checks if the provided value is greater than or equal to the specified minimum value.
func Min[T cmp.Ordered](minValue T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		value := v.(T)
		return value >= minValue, fmt.Sprintf("value %v is less than minimum value: %v", v, minValue)
	}
}

// Range returns a validator function that checks if the provided value is within the specified range (inclusive).
func Range[T cmp.Ordered](minValue, maxValue T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		value := v.(T)
		return value >= minValue && value <= maxValue, fmt.Sprintf("value %v is not in range [%v, %v]", v, minValue, maxValue)
	}
}

// NotEmpty returns a validator function that checks if the provided string value is not empty.
func NotEmpty(v any) (bool, string) {
	s, ok := v.(string)
	if !ok {
		return false, "NotEmpty validator can only be used with strings"
	}
	return s != "", "value cannot be empty"
}
