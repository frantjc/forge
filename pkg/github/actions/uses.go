package actions

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Uses struct {
	Owner      string
	Repository string
	Path       string
	Version    string
}

func (u *Uses) String() string {
	s := u.FullRepository()
	if u.Path != "" {
		s = fmt.Sprintf("%s/%s", s, u.Path)
	}
	if u.Version != "" {
		s = fmt.Sprintf("%s@%s", s, u.Version)
	}
	return s
}

func (u *Uses) GoString() string {
	return fmt.Sprintf("&Uses{%s}", u)
}

// TODO regexp.
func Parse(uses string) (*Uses, error) {
	r := &Uses{}

	spl1 := strings.Split(uses, "@")
	switch len(spl1) {
	case 2:
		r.Version = spl1[1]
	case 1:
	default:
		return r, fmt.Errorf("unable to parse uses: '%s'", uses)
	}

	spl2 := strings.Split(spl1[0], "/")
	switch len(spl2) {
	case 0, 1:
		return r, fmt.Errorf("unable to parse uses: '%s'", uses)
	default:
		r.Owner = spl2[0]
		r.Repository = spl2[1]
		if len(spl2) > 2 {
			r.Path = filepath.Join(spl2[2:]...)
		}
	}

	return r, nil
}

func (u *Uses) FullRepository() string {
	return fmt.Sprintf("%s/%s", u.Owner, u.Repository)
}

func (u *Uses) MarshalJSON() ([]byte, error) {
	return []byte("\"" + u.String() + "\""), nil
}
