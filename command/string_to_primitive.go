package command

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
)

type stringToPrimitive struct {
	Value *map[string]any

	changed bool
}

func newStringToPrimitive(val map[string]any, p *map[string]any) *stringToPrimitive {
	ssv := new(stringToPrimitive)
	ssv.Value = p
	*ssv.Value = val
	return ssv
}

// Set implements pflag.Value.
func (s *stringToPrimitive) Set(val string) error {
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
		if i, err := strconv.Atoi(kv[1]); err == nil {
			out[kv[0]] = i
		} else if b, err := strconv.ParseBool(kv[1]); err == nil {
			out[kv[0]] = b
		} else {
			out[kv[0]] = strings.Trim(kv[1], `"'`)
		}
	}
	if !s.changed {
		*s.Value = out
	} else {
		for k, v := range out {
			(*s.Value)[k] = v
		}
	}
	s.changed = true
	return nil
}

// Type implements pflag.Value.
func (s *stringToPrimitive) Type() string {
	return "stringToPrimitive"
}

// String implements pflag.Value.
func (s *stringToPrimitive) String() string {
	records := make([]string, 0, len(*s.Value)>>1)
	for k, v := range *s.Value {
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
