package unnecessary

import (
	"testing"
)

func TestValueToString(t *testing.T) {
	out1 := ValueToString(134)
	if out1 != "134" {
		t.Errorf("the operation to string for int is incorrect: %s", out1)
	}

	var someint = 432
	out2 := ValueToString(&someint)
	if out2 != "432" {
		t.Errorf("the operation to string for reference is incorrect: %s", out2)
	}

	out3 := ValueToString("opa")
	if out3 != "opa" {
		t.Errorf("the operation to string for string is incorrect: %s", out3)
	}
}

func TestValueToSlice(t *testing.T) {
	in1 := []int{3, 7}
	out1, err := ValueToSlice(&in1)
	if err != nil {
		t.Error(err)
	}
	if len(out1) != len(in1) {
		t.Fatalf("invalid len")
	}
	for i := 0; i < len(in1); i++ {
		if out1[i] != in1[i] {
			t.Errorf("element %d is invalid: %v != %v", i, out1[i], in1[i])
		}
	}

	//
	out2, err := ValueToSlice("12")
	if err == nil {
		t.Error("int converted to slice")
	}
	if out2 != nil {
		t.Error("int converted to slice with value")
	}
}
