package actions2container

import (
	"strings"

	"github.com/frantjc/forge/pkg/github/actions"
)

func UsesToVolumeName(uses *actions.Uses) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(uses.String(), "/", "-"), "@", "-",
	)
}
