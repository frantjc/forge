package githubactions

import "net/url"

type CheckoutOpts struct {
	Path      string
	Insecure  bool
	GitHubURL *url.URL
}
