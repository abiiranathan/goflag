package goflag

import (
	"cmp"
	"fmt"
)

func Choices[T comparable](choices []T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		concreteType, ok := v.(T)
		if !ok {
			return false, fmt.Sprintf("Invalid generic type for %v", v)
		}
		for _, choice := range choices {
			if choice == concreteType {
				return true, ""
			}
		}
		return false, fmt.Sprintf("%v is not a valid choice. Expected on of %v", v, choices)
	}
}

func MinStringLen(length int) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		s, ok := v.(string)
		if !ok {
			return false, "MinStringLen must be used only with strings"
		}

		return len(s) >= length, ""
	}
}

func MaxStringLen(length int) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		s, ok := v.(string)
		if !ok {
			return false, "MinStringLen must be used only with strings"
		}

		return len(s) <= length, ""
	}
}

func Max[T cmp.Ordered](maxValue T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		value := v.(T)
		return value <= maxValue, fmt.Sprintf("value %v is greater than maximum value: %v", v, maxValue)
	}
}

func Min[T cmp.Ordered](minValue T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		value := v.(T)
		return value >= minValue, fmt.Sprintf("value %v is less than minimum value: %v", v, minValue)
	}
}

func Range[T cmp.Ordered](minValue, maxValue T) func(v any) (bool, string) {
	return func(v any) (bool, string) {
		value := v.(T)
		return value >= minValue && value <= maxValue, fmt.Sprintf("value %v is not in range [%v, %v]", v, minValue, maxValue)
	}
}
