package command

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type stringToStringOrBoolValue struct {
	value   *map[string]any
	changed bool
}

func newStringToStringOrBool(val map[string]any, p *map[string]any) *stringToStringOrBoolValue {
	ssv := new(stringToStringOrBoolValue)
	ssv.value = p
	*ssv.value = val
	return ssv
}

func (s *stringToStringOrBoolValue) Set(val string) error {
	var ss []string
	n := strings.Count(val, "=")
	switch n {
	case 0:
		return fmt.Errorf("%s must be formatted as key=value", val)
	case 1:
		ss = append(ss, val)
	default:
		r := csv.NewReader(strings.NewReader(val))
		var err error
		ss, err = r.Read()
		if err != nil {
			return err
		}
	}

	out := make(map[string]any, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		if b, err := strconv.ParseBool(kv[1]); err == nil {
			out[kv[0]] = b
		} else {
			out[kv[0]] = strings.Trim(kv[1], `"'`)
		}
	}
	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}
	s.changed = true
	return nil
}

func (s *stringToStringOrBoolValue) Type() string {
	return "stringToStringOrBool"
}

func (s *stringToStringOrBoolValue) String() string {
	records := make([]string, 0, len(*s.value)>>1)
	for k, v := range *s.value {
		records = append(records, fmt.Sprint(k, "=", v))
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	if err := w.Write(records); err != nil {
		panic(err)
	}
	w.Flush()
	return "[" + strings.TrimSpace(buf.String()) + "]"
}
