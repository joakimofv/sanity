package sanity

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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
	notPublics := ""
	notSets := ""
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath != "" {
			notPublics += "\t" + t.Field(i).Name + "\n"
		} else if v.Field(i).IsZero() {
			isExcepted := false
			for _, e := range exceptedFieldNames {
				if t.Field(i).Name == e {
					isExcepted = true
					break
				}
			}
			if !isExcepted {
				notSets += "\t" + t.Field(i).Name + "\n"
			}
		} else if t.Field(i).Type.Kind() == reflect.Struct {
			// Recursive call for the sub-struct.
			if err := FieldsInitiated(v.Field(i).Interface(), ee...); err != nil {
				notSets += "\t" + t.Field(i).Name + ": " + increaseIndents(err.Error()) + "\n"
			}
		}
	}
	errStr := ""
	if notPublics != "" {
		errStr += fmt.Sprintf("Some fields not public on %v.%v:\n%s", t.PkgPath(), t.Name(), notPublics)
	}
	if notSets != "" {
		errStr += fmt.Sprintf("Some fields not set on %v.%v:\n%s", t.PkgPath(), t.Name(), notSets)
	}
	if errStr != "" {
		return errors.New(errStr)
	}
	return nil
}

func increaseIndents(s string) string {
	builder := new(strings.Builder)
	lastWasIndent := false
	for _, r := range s {
		if r == '\t' {
			lastWasIndent = true
		} else if lastWasIndent {
			builder.WriteRune('\t')
			lastWasIndent = false
		}
		builder.WriteRune(r)
	}
	return builder.String()
}
