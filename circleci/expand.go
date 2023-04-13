package circleci

import "strings"

// ExpandFunc takes a variable name and returns the variable value.
type ExpandFunc func(string) string

// ExpandString is a convenience method for Expanding strings. See Expand.
func ExpandString(s string, expand ExpandFunc) string {
	return string(Expand([]byte(s), expand))
}

// Expand takes bytes and an ExpandFunc. Whenever it encounters a variable
// in the bytes, signified by e.g. "<< example >>", it calls ExpandFunc
// with the variable name and replaces the variable with the result.
// It returns the fully expanded bytes.
func Expand(b []byte, expand ExpandFunc) (p []byte) {
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
				p = append(p, expand(name)...)
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
