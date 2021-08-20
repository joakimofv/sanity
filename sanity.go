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
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("Type passed to FieldsInitiated must be of kind 'struct', have '%v'.", t.Kind())
	}
	if !strings.HasSuffix(t.Name(), "Config") {
		return fmt.Errorf("Type passed to FieldsInitiated must be named '[...]Config', have '%v'.", t.Name())
	}
	v := reflect.ValueOf(i)
FIELD_LOOP:
	for i := 0; i < t.NumField(); i++ {
		for _, e := range exceptedFieldNames {
			if t.Field(i).Name == e {
				continue FIELD_LOOP
			}
		}
		if t.Field(i).PkgPath != "" {
			notPublics += t.Field(i).Name + "\n"
		} else if t.Field(i).Type.Kind() == reflect.Struct && strings.HasSuffix(t.Field(i).Type.Name(), "Config") {
			// Recursive call for the nested Config struct.
			if err := FieldsInitiated(v.Field(i).Interface(), ee...); err != nil {
				if errors.As(err, &NotPublicError{}) {
					notPublics += t.Field(i).Name + " -> " + err.Error() + "\n"
				} else if v.Field(i).IsZero() {
					notSets += t.Field(i).Name + "\n"
				} else {
					notSets += t.Field(i).Name + " -> " + err.Error() + "\n"
				}
			}
		} else if v.Field(i).IsZero() {
			notSets += t.Field(i).Name + "\n"
		}
	}
	if notPublics != "" {
		return NotPublicError{fmt.Sprintf("Some fields not public on %v.%v:\n%s", t.PkgPath(), t.Name(), addIndent(notPublics))}
	}
	if notSets != "" {
		return NotSetError{fmt.Sprintf("Some fields not set on %v.%v:\n%s", t.PkgPath(), t.Name(), addIndent(notSets))}
	}
	return nil
}

type NotPublicError struct {
	err string
}

func (err NotPublicError) Error() string {
	return string(err.err)
}

type NotSetError struct {
	err string
}

func (err NotSetError) Error() string {
	return string(err.err)
}

func addIndent(s string) string {
	s = strings.TrimRight(s, "\n")
	builder := new(strings.Builder)
	indent := `    `
	builder.WriteString(indent)
	for _, r := range s {
		builder.WriteRune(r)
		if r == '\n' {
			builder.WriteString(indent)
		}
	}
	return builder.String()
}
