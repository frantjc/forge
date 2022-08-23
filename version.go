package forge

import "runtime/debug"

var (
	Version    = "0.0.0"
	Prerelease = ""
)

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
