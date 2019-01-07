package validator

import (
	"reflect"
	"testing"
)

func TestOrderValPSByPledge_Success(t *testing.T) {
	first := CreateValidatorWithPledgeAmountForTests(1)
	second := CreateValidatorWithPledgeAmountForTests(2)
	third := CreateValidatorWithPledgeAmountForTests(23)
	fourth := CreateValidatorWithPledgeAmountForTests(23)
	fifth := CreateValidatorWithPledgeAmountForTests(45)

	input := []Validator{
		first,
		fifth,
		third,
		fourth,
		second,
	}

	expected := []Validator{
		first,
		second,
		third,
		fourth,
		fifth,
	}

	out, _ := orderValPSByPledge(input, 0, len(input))
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}

}

func TestOrderValPSByPledge_withNonZeroIndex_Success(t *testing.T) {
	first := CreateValidatorWithPledgeAmountForTests(1)
	second := CreateValidatorWithPledgeAmountForTests(2)
	third := CreateValidatorWithPledgeAmountForTests(23)
	fourth := CreateValidatorWithPledgeAmountForTests(23)
	fifth := CreateValidatorWithPledgeAmountForTests(45)

	input := []Validator{
		first,
		fifth,
		third,
		fourth,
		second,
	}

	expected := []Validator{
		third,
		fourth,
		fifth,
	}

	out, _ := orderValPSByPledge(input, 2, len(input))
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}
}

func TestOrderValPSByPledge_withTooBigIndex_Success(t *testing.T) {
	first := CreateValidatorWithPledgeAmountForTests(1)
	second := CreateValidatorWithPledgeAmountForTests(2)
	third := CreateValidatorWithPledgeAmountForTests(23)
	fourth := CreateValidatorWithPledgeAmountForTests(23)
	fifth := CreateValidatorWithPledgeAmountForTests(45)

	input := []Validator{
		first,
		fifth,
		third,
		fourth,
		second,
	}

	expected := []Validator{}
	out, _ := orderValPSByPledge(input, 20, len(input))
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}
}

func TestOrderValPSByPledge_withTooBigAmount_Success(t *testing.T) {
	first := CreateValidatorWithPledgeAmountForTests(1)
	second := CreateValidatorWithPledgeAmountForTests(2)
	third := CreateValidatorWithPledgeAmountForTests(23)
	fourth := CreateValidatorWithPledgeAmountForTests(23)
	fifth := CreateValidatorWithPledgeAmountForTests(45)

	input := []Validator{
		first,
		fifth,
		third,
		fourth,
		second,
	}

	expected := []Validator{
		first,
		second,
		third,
		fourth,
		fifth,
	}

	out, _ := orderValPSByPledge(input, 0, len(input)*30)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}
}

func TestOrderValPSByPledge_withMinusIndex_Success(t *testing.T) {
	first := CreateValidatorWithPledgeAmountForTests(1)
	second := CreateValidatorWithPledgeAmountForTests(2)
	third := CreateValidatorWithPledgeAmountForTests(23)
	fourth := CreateValidatorWithPledgeAmountForTests(23)
	fifth := CreateValidatorWithPledgeAmountForTests(45)

	input := []Validator{
		first,
		fifth,
		third,
		fourth,
		second,
	}

	expected := []Validator{
		first,
		second,
		third,
		fourth,
		fifth,
	}

	out, _ := orderValPSByPledge(input, -20, len(input))
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}
}

func TestOrderValPSByPledge_withMinusAmount_Success(t *testing.T) {
	first := CreateValidatorWithPledgeAmountForTests(1)
	second := CreateValidatorWithPledgeAmountForTests(2)
	third := CreateValidatorWithPledgeAmountForTests(23)
	fourth := CreateValidatorWithPledgeAmountForTests(23)
	fifth := CreateValidatorWithPledgeAmountForTests(45)

	input := []Validator{
		first,
		fifth,
		third,
		fourth,
		second,
	}

	expected := []Validator{}
	out, _ := orderValPSByPledge(input, 0, -20)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}
}

func TestOrderValPSByPledge_withZeroElements_Success(t *testing.T) {
	input := []Validator{}
	expected := []Validator{}
	out, _ := orderValPSByPledge(input, 0, 50)
	if !reflect.DeepEqual(out, expected) {
		t.Errorf("the output is invalid")
		return
	}
}
