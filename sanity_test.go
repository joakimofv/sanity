package sanity

import (
	"testing"
	"time"
)

type publicStruct struct {
	Fun1  func()
	Fun2  func(bool) bool
	Int   int
	Hello string
	Time  time.Duration
}

func TestFieldsInitiated(t *testing.T) {
	myStruct := publicStruct{}

	err := FieldsInitiated(myStruct)
	if err == nil {
		t.Error("expected error since nothing is set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}

	myStruct.Fun1 = func() {}
	myStruct.Fun2 = func(truth bool) bool { return !truth }

	err = FieldsInitiated(myStruct)
	if err == nil {
		t.Error("expected error since simple types not set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
	err = FieldsInitiated(myStruct,
		Except("Int"),
		Except("Hello"),
		Except("Time"),
	)
	if err != nil {
		t.Errorf("expected nil error since simple types excepted. err: %v", err)
	}

	myStruct.Int = 123
	myStruct.Hello = "Hi"
	myStruct.Time = time.Second

	err = FieldsInitiated(myStruct)
	if err != nil {
		t.Errorf("expected nil error. err: %v", err)
	}
}

type privateStruct struct {
	fun func()
	i   int
	Str string
}

func TestFieldsPublic(t *testing.T) {
	myStruct := privateStruct{}
	myStruct.fun = func() {}
	myStruct.i = 123
	myStruct.Str = "abc"

	err := FieldsInitiated(myStruct)
	if err == nil {
		t.Error("expected error since some fields are not public.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
}
