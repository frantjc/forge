package bin

import "regexp"

var r = regexp.MustCompile(`^\s*#!.+\n`)

func HasShebang(script string) bool {
	return r.MatchString(script)
}
