package unixtable

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"
)

var (
	ErrTypeNotSupported = errors.New("type not supported")
)

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{tabwriter.NewWriter(w, DefaultMinWidth, DefaultTabWidth, DefaultPadding, ' ', tabwriter.DiscardEmptyColumns)}
}

type Encoder struct {
	*tabwriter.Writer
}

func (e *Encoder) Encode(a any) error {
	if marshaler, ok := a.(Marshaler); ok {
		b, err := marshaler.MarshalUnixTable()
		if err != nil {
			return err
		}

		_, err = e.Write(b)
		return err
	}

	v := reflect.ValueOf(a)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	k := v.Kind()
	switch k {
	case reflect.Struct:
		var (
			numFields = v.NumField()
			keys      = make([]string, numFields)
			vals      = make([]string, numFields)
		)

		for i := 0; i < numFields; i++ {
			key := v.Type().Field(i).Name
			if tag := v.Type().Field(i).Tag.Get(Tag); tag != "" {
				key = tag
			}

			keys[i] = key
			vals[i] = fmt.Sprint(v.Field(i))
		}

		if _, err := fmt.Fprintln(e, strings.Join(keys, "\t")); err != nil {
			return err
		}

		if _, err := fmt.Fprintln(e, strings.Join(vals, "\t")); err != nil {
			return err
		}
	default:
		return ErrTypeNotSupported
	}

	return e.Flush()
}
