package sanity

import (
	"errors"
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
	pubStruct := publicStruct{}
	err := FieldsInitiated(pubStruct)
	if err == nil {
		t.Error("expected error since nothing is set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}

	pubStruct.Fun1 = func() {}
	pubStruct.Fun2 = func(truth bool) bool { return !truth }
	err = FieldsInitiated(pubStruct)
	if err == nil {
		t.Error("expected error since simple types not set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
	err = FieldsInitiated(pubStruct,
		Except("Int"),
		Except("Hello"),
		Except("Time"),
	)
	if err != nil {
		t.Errorf("expected nil error since simple types excepted. err: %v", err)
	}

	pubStruct.Int = 123
	pubStruct.Hello = "Hi"
	pubStruct.Time = time.Second
	err = FieldsInitiated(pubStruct)
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
	privStruct := privateStruct{}
	privStruct.fun = func() {}
	privStruct.i = 123
	privStruct.Str = "abc"
	err := FieldsInitiated(privStruct)
	if err == nil {
		t.Error("expected error since some fields are not public.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
}

type superStruct struct {
	Int       int
	SubStruct publicStruct
}
type superStructPriv struct {
	Int       int
	SubStruct privateStruct
}

func TestSubStruct(t *testing.T) {
	topStruct := superStruct{}
	err := FieldsInitiated(topStruct)
	if err == nil {
		t.Error("expected error since nothing is set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}

	topStruct.Int = 1
	err = FieldsInitiated(topStruct)
	if err == nil {
		t.Error("expected error since nothing is set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}

	pubStruct := publicStruct{}
	pubStruct.Fun1 = func() {}
	pubStruct.Fun2 = func(truth bool) bool { return !truth }
	topStruct.SubStruct = pubStruct
	err = FieldsInitiated(topStruct)
	if err == nil {
		t.Error("expected error since simple types not set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
	err = FieldsInitiated(topStruct,
		Except("Int"),
		Except("Hello"),
		Except("Time"),
	)
	if err != nil {
		t.Errorf("expected nil error since simple types excepted. err: %v", err)
	}

	pubStruct.Int = 123
	pubStruct.Hello = "Hi"
	pubStruct.Time = time.Second
	topStruct.SubStruct = pubStruct
	err = FieldsInitiated(topStruct)
	if err != nil {
		t.Errorf("expected nil error. err: %v", err)
	}

	topStruct2 := superStructPriv{} // Not set and not public
	err = FieldsInitiated(topStruct2)
	if !errors.As(err, &NotPublicError{}) {
		t.Errorf("expected NotPublicError error, got err: %v", err)
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
}
