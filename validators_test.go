package goflag

import (
	"testing"
)

func TestChoices(t *testing.T) {
	choicesValidator := Choices([]int{1, 2, 3})
	err := choicesValidator(2)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = choicesValidator(4)
	if err == nil {
		t.Errorf("choice validator should have returned an error for value 4")
	}
}

func TestMinStringLen(t *testing.T) {
	minLenValidator := MinStringLen(5)
	err := minLenValidator("hello")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = minLenValidator("hi")
	if err == nil {
		t.Errorf("min string len validator should have returned an error for string with length less than min")
	}
}

func TestMaxStringLen(t *testing.T) {
	maxLenValidator := MaxStringLen(5)
	err := maxLenValidator("hello")
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = maxLenValidator("hello world")
	if err == nil {
		t.Errorf("max string len validator should have returned an error for string with length greater than max")
	}
}

func TestMax(t *testing.T) {
	maxValidator := Max(10)
	err := maxValidator(9)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = maxValidator(11)
	if err == nil {
		t.Errorf("max validator should have returned an error for value greater than max")
	}
}

func TestMin(t *testing.T) {
	minValidator := Min(5)
	err := minValidator(6)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = minValidator(4)
	if err == nil {
		t.Errorf("min validator should have returned an error for value less than min")
	}
}

func TestRange(t *testing.T) {
	rangeValidator := Range(5, 10)
	err := rangeValidator(7)
	if err != nil {
		t.Errorf("%s", err.Error())
	}

	err = rangeValidator(4)
	if err == nil {
		t.Errorf("range validator should have returned an error for value less than min")
	}

	err = rangeValidator(11)
	if err == nil {
		t.Errorf("range validator should have returned an error for value greater than max")
	}
}
