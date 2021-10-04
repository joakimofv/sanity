package sanity

import (
	"errors"
	"testing"
)

func TestReturns(t *testing.T) {
	var a int
	var b MyInterface
	var c *string
	var err error
	for name, tc := range map[string]struct {
		errNonLast bool
		errSet     bool
		errTyped   bool
		abcSet     bool
		typedPtr   bool
		expectErr  bool
	}{
		"unset":           {},
		"errNonLast":      {errNonLast: true, expectErr: true},
		"errSet":          {errSet: true},
		"errTyped":        {errSet: true},
		"abcSet":          {abcSet: true},
		"errSet,abcSet":   {errSet: true, abcSet: true, expectErr: true},
		"errTyped,abcSet": {errTyped: true, abcSet: true, expectErr: true},
		"errSet,typedPtr": {errSet: true, typedPtr: true},
	} {
		t.Run(name, func(t *testing.T) {

			if tc.errSet {
				err = errors.New("test error")
			} else if tc.errTyped {
				err = &myError{}
			} else {
				err = nil
			}
			if tc.abcSet {
				a = 1
				mt := MyType("stuff")
				b = &mt
				str := "hej"
				c = &str
			} else {
				a = 0
				if tc.typedPtr {
					b = (*MyType)(nil)
				} else {
					b = nil
				}
				c = nil
			}

			vars := []interface{}{a, b, c, err}
			if tc.errNonLast {
				vars = []interface{}{a, b, err, c}
			}

			if !tc.expectErr {
				if retErr := Returns(vars...); retErr != nil {
					t.Error(retErr)
				}
			} else {
				if retErr := Returns(vars...); retErr == nil {
					t.Error("Expected error, got nil.")
				} else {
					t.Logf("got error as expected:\n%v", retErr)
				}
			}
		})
	}
}

type myError struct {
	err string
}

func (err *myError) Error() string {
	return err.err
}

type MyInterface interface {
	doStuff()
}

type MyType string

func (mt *MyType) doStuff() {
	panic(string(*mt))
}
