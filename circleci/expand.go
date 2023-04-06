package circleci

import "strings"

type ExpandFunc func(string) string

func (e ExpandFunc) ExpandString(s string) string {
	return ExpandString(s, e)
}

func (e ExpandFunc) Expand(b []byte) []byte {
	return Expand(b, e)
}

func ExpandString(s string, mapping ExpandFunc) string {
	return string(Expand([]byte(s), mapping))
}

func Expand(b []byte, mapping ExpandFunc) (p []byte) {
	i := 0
	for j := 0; j < len(b); j++ {
		if b[j] == '<' {
			if p == nil {
				p = make([]byte, 0, 2*len(b))
			}
			p = append(p, b[i:j]...)
			name, w := getParameterName(b[j:])
			switch {
			case name == "" && w > 0:
				// encountered invalid syntax; eat the characters
			case name == "":
				// valid syntax, but << >> contained no name
				p = append(p, b[j])
			default:
				p = append(p, mapping(name)...)
			}
			j += w
			i = j + 1
		}
	}
	if p == nil {
		return b
	} else if i >= len(b) {
		i = len(b)
	}
	return append(p, b[i:]...)
}

func getParameterName(b []byte) (s string, w int) {
	if len(b) > 3 && b[0] == '<' && b[1] == '<' {
		i := 2
		//nolint:revive
		for ; i+1 < len(b) && b[i] != '>'; i++ {
		}
		if b[i] == '>' && i+1 < len(b) && b[i+1] != '>' {
			return "", 0 // bad syntax
		}
		return strings.TrimSpace(string(b[2:i])), i + 2
	}

	return "", 0
}
