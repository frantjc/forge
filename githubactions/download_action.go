package githubactions

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/frantjc/forge"
	"github.com/frantjc/forge/internal/tarutil"
	"github.com/google/go-github/v50/github"
)

// ActionYAMLFilenames holds the possible names of
// a GitHub Action metadata file.
var ActionYAMLFilenames = []string{"action.yml", "action.yaml"}

// DownloadAction takes a Uses reference and returns the corresponding GitHub Action Metadata,
// a tarball of the GitHub Action repository and a common path prefix shared by all of the
// tar headers for files from the repository.
// TODO use github.com/frantjc/forge/internal/tarutil.StripPrefix
// instead of returning the prefix to be stripped.
func DownloadAction(ctx context.Context, u *Uses) (*Metadata, io.ReadCloser, error) {
	var (
		_        = forge.LoggerFrom(ctx)
		client   *github.Client
		metadata *Metadata
		sha      = u.Version
	)

	if u.IsLocal() {
		return nil, nil, fmt.Errorf("clone local action: %s", u.Path)
	}

	if token := os.Getenv(EnvVarToken); token != "" {
		client = github.NewTokenClient(ctx, os.Getenv(EnvVarToken))
	} else {
		client = github.NewClient(http.DefaultClient)
	}

	if ref, _, err := client.Git.GetRef(ctx, u.GetOwner(), u.GetRepository(), "tags/"+u.Version); err == nil {
		sha = ref.GetObject().GetSHA()
	} else {
		if ref, _, err := client.Git.GetRef(ctx, u.GetOwner(), u.GetRepository(), "heads/"+u.Version); err == nil {
			sha = ref.GetObject().GetSHA()
		}
	}

	if len(sha) < 7 {
		panic("unable to get action sha")
	}

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

	return metadata, tarutil.StripPrefix(res.Body, u.GetOwner()+"-"+u.GetRepository()+"-"+sha[0:7]+"/", tarutil.WithGzip), nil
}
