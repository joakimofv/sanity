package sanity

import (
	"fmt"
	"reflect"
)

// Except can optionally be used to wrap extra parameters to FieldsInitiated,
// for semantic clarity in that those field names are "excepted" from the check.
func Except(s string) string {
	return s
}

// FieldsInitiated is meant for checking that Configs are not neglected to be filled in (otherwise easy to forget).
func FieldsInitiated(i interface{}, ee ...string) error {
	var exceptedFieldNames []string
	for _, e := range ee {
		exceptedFieldNames = append(exceptedFieldNames, e)
	}
	problems := ""
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		if field.IsZero() {
			isExcepted := false
			for _, e := range exceptedFieldNames {
				if t.Field(i).Name == e {
					isExcepted = true
					break
				}
			}
			if !isExcepted {
				problems += t.Field(i).Name + "\n"
			}
		}
	}
	if problems != "" {
		return fmt.Errorf("Some fields not set on %v.%v:\n%s", t.PkgPath(), t.Name(), problems)
	}
	return nil
}
