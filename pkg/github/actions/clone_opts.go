package actions

import "net/url"

type CloneOpts struct {
	Path      string
	Insecure  bool
	GitHubURL *url.URL
}
