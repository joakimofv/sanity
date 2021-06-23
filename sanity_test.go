package sanity

import (
	"testing"
	"time"
)

func TestFieldsInitiated(t *testing.T) {
	myStruct := struct {
		Fun1  func()
		Fun2  func(bool) bool
		Int   int
		Hello string
		Time  time.Duration
	}{}

	err := FieldsInitiated(myStruct)
	if err == nil {
		t.Error("expected error since nothing is set.")
	}

	myStruct.Fun1 = func() {}
	myStruct.Fun2 = func(truth bool) bool { return !truth }

	err = FieldsInitiated(myStruct)
	if err == nil {
		t.Error("expected error since simple types not set.")
	}
	err = FieldsInitiated(myStruct,
		"Int",
		"Hello",
		"Time",
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

func TestFieldsPublic(t *testing.T) {
	myStruct := struct {
		fun  func()
		i   int
		Str	string
	}{}
	myStruct.fun = func() {}
	myStruct.i = 123
	myStruct.Str = "abc"

	err := FieldsInitiated(myStruct)
	if err == nil {
		t.Error("expected error since some fields are not public.")
	}
}
