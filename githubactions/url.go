package githubactions

import (
	"net/url"
	"os"
	"path"
)

var (
	DefaultURL        *url.URL
	DefaultAPIURL     *url.URL
	DefaultGraphQLURL *url.URL
)

func init() {
	var err error
	if DefaultURL, err = url.Parse("https://github.com/"); err != nil {
		panic("github.com/frantjc/forge/githubactions.DefaultURL is not a valid URL")
	}

	if DefaultAPIURL, err = APIURLFromBaseURL(DefaultURL); err != nil {
		panic("github.com/frantjc/forge/githubactionsDefaultAPIURL is not a valid URL")
	}

	DefaultGraphQLURL, err = GraphQLURLFromBaseURL(DefaultURL)
	if err != nil {
		panic("github.com/frantjc/forge/githubactions.DefaultGraphQLURL is not a valid URL")
	}
}

func APIURLFromBaseURL(base *url.URL) (*url.URL, error) {
	api, err := url.Parse(base.String())
	if err != nil {
		return nil, err
	}

	if api.Hostname() == "github.com" {
		api.Host = "api.github.com"
	} else {
		api.Path = path.Join(api.Path, "/api/v3")
	}
	return api, nil
}

func GraphQLURLFromBaseURL(base *url.URL) (*url.URL, error) {
	graphql, err := APIURLFromBaseURL(base)
	if err != nil {
		return nil, err
	}
	graphql.Path = path.Join(graphql.Path, "graphql")
	return graphql, nil
}

func GetGitHubURL() *url.URL {
	envVar := os.Getenv(EnvVarServerURL)
	if u, err := url.Parse(envVar); err == nil && envVar != "" {
		return u
	}

	return DefaultURL
}

func GetGitHubServerURL() *url.URL {
	return GetGitHubURL()
}

func GetGitHubAPIURL() *url.URL {
	if u, err := url.Parse(os.Getenv(EnvVarAPIURL)); err == nil {
		return u
	}

	return DefaultAPIURL
}

func GetGraphQLURL() *url.URL {
	if u, err := url.Parse(os.Getenv(EnvVarGraphQLURL)); err == nil {
		return u
	}

	return DefaultGraphQLURL
}
