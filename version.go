package forge

import "runtime/debug"

var (
	// Version is the major.minor.patch version of forge.
	// Used to build a semantic version.
	// Meant to be be overridden at build time.
	Version = "0.0.0"
	// Prerelease is the prelease of forge e.g. "alpha".
	// Used to build a semantic version.
	// Meant to be overridden at build time.
	Prerelease = ""
)

// Semver returns the semantic version of forge as built from
// Version, Prerelease and debug build info.
func Semver() string {
	version := Version

	if Prerelease != "" {
		version += "-" + Prerelease
	}

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
			version += "+" + revision[:i]
		}

		if modified {
			version += "*"
		}
	}

	return version
}
