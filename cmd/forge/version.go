package main

import (
	"runtime/debug"
	"strings"
)

var (
	version = ""
)

// SemVer returns the semantic version of `sindri` as
// built from ldflags and debug build info.
func SemVer() string {
	semver := version

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

			if !strings.Contains(semver, revision[:i]) {
				semver += "+" + revision[:i]
			}
		}

		if modified {
			semver += "*"
		}
	}

	return semver
}
