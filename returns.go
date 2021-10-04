package sanity

import (
	"fmt"
	"reflect"
)

// Returns checks that the return variables from a function are sane, meaning:
// * Last var should be an error.
// * If the error is non-nil then all other vars must be unset.
func Returns(vars ...interface{}) error {
	if len(vars) < 1 {
		return nil
	}
	lastVar := vars[len(vars)-1]
	errorInterface := reflect.TypeOf((*error)(nil)).Elem()
	t := reflect.TypeOf(lastVar)
	if t == nil {
		// Can't tell if it's an error or not, so assume it's a nil error and let it pass.
		return nil
	}
	if !t.Implements(errorInterface) {
		return fmt.Errorf("Last variable must be an error type, but is '%T'.", lastVar)
	}

	if !reflect.ValueOf(lastVar).IsZero() {
		// err is set, so all else must be unset.
		sets := ""
		for i := 0; i < len(vars)-1; i++ {
			v := reflect.ValueOf(vars[i])
			if !v.IsValid() {
				// An untyped nil.
				continue
			}
			if !v.IsZero() {
				name := ""
				for v.Kind() == reflect.Ptr {
					name += "*"
					v = v.Elem()
				}
				name += v.Type().Name()
				sets += fmt.Sprintf("[%d] %v: %v\n", i, name, v.Interface())
			}
		}
		if sets != "" {
			return fmt.Errorf("Some vars set when err is also set:\n%s", addIndent(sets))
		}
	}

	return nil
}
