package forge

import (
	"runtime/debug"
	"strings"
)

// Semver is the semantic version of Forge.
// Meant to be be overridden at build time,
// but kept up-to-date sometimes to best
// support `go install`.
var Semver = "0.13.5"

// GetSemver returns the semantic version of Forge as built from
// Semver and debug build info.
func GetSemver() string {
	version := Semver

	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		var (
			revision string
			modified bool
		)
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.modified":
				modified = setting.Value == "true"
			}
		}

		if revision != "" {
			i := len(revision)
			if i > 7 {
				i = 7
			}

			if !strings.Contains(version, revision[:i]) {
				version += "+" + revision[:i]
			}
		}

		if modified {
			version += "*"
		}
	}

	return version
}
