package githubactions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/tarutil"
	"github.com/google/go-github/v50/github"
)

// ActionYAMLFilenames holds the possible names of
// a GitHub Action metadata file.
var ActionYAMLFilenames = []string{"action.yml", "action.yaml"}

// DownloadAction takes a Uses reference and returns the corresponding GitHub Action Metadata,
// a tarball of the GitHub Action repository.
func DownloadAction(ctx context.Context, u *Uses) (*Metadata, io.ReadCloser, error) {
	var (
		_        = forge.LoggerFrom(ctx)
		client   *github.Client
		metadata *Metadata
	)

	if u.IsLocal() {
		return nil, nil, fmt.Errorf("clone local action: %s", u.Path)
	}

	if token := os.Getenv(EnvVarToken); token != "" {
		client = github.NewTokenClient(ctx, token)
	} else {
		client = github.NewClient(http.DefaultClient)
	}

	client.BaseURL = GetGitHubAPIURL().JoinPath("/")

	// Get the sha in parallel for speed.
	// Used later to know what directory of the action's tarball
	// the repository's contents actually reside in.
	shaC := make(chan string, 1)
	go func() {
		defer close(shaC)

		if ref, _, err := client.Git.GetRef(ctx, u.GetOwner(), u.GetRepository(), "tags/"+u.Version); err == nil {
			shaC <- ref.GetObject().GetSHA()
		} else {
			if ref, _, err := client.Git.GetRef(ctx, u.GetOwner(), u.GetRepository(), "heads/"+u.Version); err == nil {
				shaC <- ref.GetObject().GetSHA()
			} else {
				shaC <- u.Version
			}
		}
	}()

	for _, filename := range ActionYAMLFilenames {
		rc, _, err := client.Repositories.DownloadContents(ctx, u.GetOwner(), u.GetRepository(), u.GetActionPath()+"/"+filename, &github.RepositoryContentGetOptions{
			Ref: u.Version,
		})
		if err != nil {
			return nil, nil, err
		}
		defer rc.Close()

		metadata, err = NewMetadataFromReader(rc)
		if err != nil {
			return nil, nil, err
		}

		if metadata != nil {
			break
		}
	}

	if metadata == nil {
		return nil, nil, ErrNotAnAction
	}

	link, _, err := client.Repositories.GetArchiveLink(
		ctx,
		u.GetOwner(), u.GetRepository(),
		github.Tarball,
		&github.RepositoryContentGetOptions{Ref: u.Version},
		true,
	)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	sha := <-shaC

	if matched, err := regexp.MatchString("[0-9a-f]{40}", sha); err != nil {
		return nil, nil, err
	} else if !matched {
		return nil, nil, fmt.Errorf("unable to get action sha")
	}

	return metadata, tarutil.Subdir(res.Body, u.GetOwner()+"-"+u.GetRepository()+"-"+sha[0:7]+"/", tarutil.IsGzipped), nil
}
