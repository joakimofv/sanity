package sanity

import (
	"errors"
	"testing"
	"time"
)

type badNameStruct struct {
	Int int
}
type badTypeConfig map[string]int

func TestBadArgument(t *testing.T) {
	badStruct := badNameStruct{Int: 4}
	err := FieldsInitiated(badStruct)
	if err == nil {
		t.Error("expected error since struct name doesn't end with Config.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}

	badType := make(badTypeConfig)
	badType["Int"] = 4
	err = FieldsInitiated(badType)
	if err == nil {
		t.Error("expected error since arg is not a struct.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
}

type publicStructConfig struct {
	Fun1  func()
	Fun2  func(bool) bool
	Int   int
	Hello string
	Time  time.Duration
}

func TestFieldsInitiated(t *testing.T) {
	pubCfg := publicStructConfig{}
	err := FieldsInitiated(pubCfg)
	if err == nil {
		t.Error("expected error since nothing is set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}

	pubCfg.Fun1 = func() {}
	pubCfg.Fun2 = func(truth bool) bool { return !truth }
	err = FieldsInitiated(pubCfg)
	if err == nil {
		t.Error("expected error since simple types not set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
	err = FieldsInitiated(pubCfg,
		Except("Int"),
		Except("Hello"),
		Except("Time"),
	)
	if err != nil {
		t.Errorf("expected nil error since simple types excepted. err: %v", err)
	}

	pubCfg.Int = 123
	pubCfg.Hello = "Hi"
	pubCfg.Time = time.Second
	err = FieldsInitiated(pubCfg)
	if err != nil {
		t.Errorf("expected nil error. err: %v", err)
	}
}

type privateStructConfig struct {
	fun func()
	i   int
	Str string
}

func TestFieldsPublic(t *testing.T) {
	privCfg := privateStructConfig{}
	privCfg.fun = func() {}
	privCfg.i = 123
	privCfg.Str = "abc"
	err := FieldsInitiated(privCfg)
	if err == nil {
		t.Error("expected error since some fields are not public.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
}

type superStructConfig struct {
	Int       int
	SubStruct publicStructConfig
}
type superStructPrivConfig struct {
	Int       int
	SubStruct privateStructConfig
}

func TestSubStruct(t *testing.T) {
	topStruct := superStructConfig{}
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

	pubCfg := publicStructConfig{}
	pubCfg.Fun1 = func() {}
	pubCfg.Fun2 = func(truth bool) bool { return !truth }
	topStruct.SubStruct = pubCfg
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

	pubCfg.Int = 123
	pubCfg.Hello = "Hi"
	pubCfg.Time = time.Second
	topStruct.SubStruct = pubCfg
	err = FieldsInitiated(topStruct)
	if err != nil {
		t.Errorf("expected nil error. err: %v", err)
	}

	topStruct2 := superStructPrivConfig{} // Not set and not public
	err = FieldsInitiated(topStruct2)
	if !errors.As(err, &NotPublicError{}) {
		t.Errorf("expected NotPublicError error, got err: %v", err)
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
}

func TestSpecificFieldsInitiated(t *testing.T) {
	pubCfg := publicStructConfig{}
	err := SpecificFieldsInitiated(pubCfg)
	if err != nil {
		t.Errorf("expected nil error. err: %v", err)
	}

	pubCfg.Fun1 = func() {}
	pubCfg.Fun2 = func(truth bool) bool { return !truth }
	err = SpecificFieldsInitiated(pubCfg, "Fun1", "Fun2", "Int", "Hello", "Time")
	if err == nil {
		t.Error("expected error since simple types not set.")
	} else {
		t.Logf("got error as expected:\n%v", err)
	}
	err = SpecificFieldsInitiated(pubCfg, "Fun1", "Fun2")
	if err != nil {
		t.Errorf("expected nil error since simple types excepted. err: %v", err)
	}

	pubCfg.Int = 123
	pubCfg.Hello = "Hi"
	pubCfg.Time = time.Second
	err = SpecificFieldsInitiated(pubCfg, "Fun1", "Fun2", "Int", "Hello", "Time")
	if err != nil {
		t.Errorf("expected nil error. err: %v", err)
	}
}
