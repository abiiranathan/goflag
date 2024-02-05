package goflag

import (
	"testing"
)

func TestChoices(t *testing.T) {
	choicesValidator := Choices([]int{1, 2, 3})
	valid, _ := choicesValidator(2)
	if !valid {
		t.Errorf("Choices validator failed for value in choices")
	}

	valid, _ = choicesValidator(4)
	if valid {
		t.Errorf("Choices validator passed for value not in choices")
	}
}

func TestMinStringLen(t *testing.T) {
	minLenValidator := MinStringLen(5)
	valid, _ := minLenValidator("hello")
	if !valid {
		t.Errorf("MinStringLen validator failed for string with length equal to min")
	}

	valid, _ = minLenValidator("hi")
	if valid {
		t.Errorf("MinStringLen validator passed for string with length less than min")
	}
}

func TestMaxStringLen(t *testing.T) {
	maxLenValidator := MaxStringLen(5)
	valid, _ := maxLenValidator("hello")
	if !valid {
		t.Errorf("MaxStringLen validator failed for string with length equal to max")
	}

	valid, _ = maxLenValidator("hello world")
	if valid {
		t.Errorf("MaxStringLen validator passed for string with length greater than max")
	}
}

func TestMax(t *testing.T) {
	maxValidator := Max(10)
	valid, _ := maxValidator(9)
	if !valid {
		t.Errorf("Max validator failed for value less than max")
	}

	valid, _ = maxValidator(11)
	if valid {
		t.Errorf("Max validator passed for value greater than max")
	}
}

func TestMin(t *testing.T) {
	minValidator := Min(5)
	valid, _ := minValidator(6)
	if !valid {
		t.Errorf("Min validator failed for value greater than min")
	}

	valid, _ = minValidator(4)
	if valid {
		t.Errorf("Min validator passed for value less than min")
	}
}

func TestRange(t *testing.T) {
	rangeValidator := Range(5, 10)
	valid, _ := rangeValidator(7)
	if !valid {
		t.Errorf("Range validator failed for value within range")
	}

	valid, _ = rangeValidator(4)
	if valid {
		t.Errorf("Range validator passed for value less than range")
	}

	valid, _ = rangeValidator(11)
	if valid {
		t.Errorf("Range validator passed for value greater than range")
	}
}
