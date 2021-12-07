[![Go Reference](https://pkg.go.dev/badge/github.com/joakimofv/sanity.svg)](https://pkg.go.dev/github.com/joakimofv/sanity)

# sanity

Opinionated sanity checking of structs. Intended for structs named "Config" that are given to state machines before start up.
Having any field that is not actively set could result in bugs, especially if it is a function (since calling a nil function causes panic).

Useful during development as a reminder to set all values.

## Intended Use

Suppose that your main function looks like this:

```go
cfg := mypkg.Config{
	MyInt: 123,
	// You are supposed to set all values here, or at least before Start().
}

machine := mypkg.New(cfg)
if err := machine.Start(); err != nil {
	return err
}
```

Then do the sanity check inside mypkg:

```go
import (
	"github.com/joakimofv/sanity"
)

type Config struct {
	//...
}

type Machine struct {
	Cfg Config
	//...
}

func New(cfg Config) *Machine {
	machine := &Machine{
		Cfg: cfg,
		//...
	}
	//...
}

func (m *Machine) Start() error {
	if err := sanity.FieldsInitiated(m.Cfg); err != nil {
		return err
	}
	//...
}
```

## Excepting Fields

For some fields it might be natural to use the zero-value, then except them from the check:

```go
if err := sanity.FieldsInitiated(m.Cfg,
	sanity.Except("MyInt"),
	sanity.Except("MyFloat"),
	sanity.Except("MyBool")); err != nil {
	return err
}
```

## Specific Fields

To check on only specifically specified fields, use the SpecificFieldsInitiated function like this:

```go
if err := sanity.SpecificFieldsInitiated(m.Cfg,	"MyInt", "MyFloat"); err != nil {
	return err
}
```
