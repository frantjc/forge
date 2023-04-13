package envconv

import (
	"fmt"
	"strings"
)

// ArrToMap takes an environment array of the form
//
//	["KEY1=val1", "KEY2=val2"]
//
// and returns a corresponding map of the form
//
//	{
//		"KEY1": "val1",
//		"KEY2": "val2"
//	}
func ArrToMap(a []string) map[string]string {
	m := map[string]string{}
	for _, s := range a {
		kv := strings.Split(s, "=")
		if len(kv) >= 2 && kv[0] != "" {
			m[kv[0]] = kv[1]
		}
	}
	return m
}

// ToMap is a convenience function wrapping ArrToMap.
// It takes an environment array of the form
//
//	["KEY1=val1", "KEY2=val2"]
//
// and returns a corresponding map of the form
//
//	{
//		"KEY1": "val1",
//		"KEY2": "val2"
//	}
func ToMap(ss ...string) map[string]string {
	return ArrToMap(ss)
}

// MapToArr takes an map of the form
//
//	{
//		"KEY1": "val1",
//		"KEY2": "val2"
//	}
//
// and returns a corresponding array of the form
//
//	["KEY1=val1", "KEY2=val2"].
func MapToArr(m map[string]string) []string {
	a := []string{}
	for k, v := range m {
		if k != "" {
			a = append(a, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return a
}
