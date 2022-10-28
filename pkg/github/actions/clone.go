package actions

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/frantjc/forge"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	ActionYAMLFilenames = []string{"action.yml", "action.yaml"}
)

func CloneUses(ctx context.Context, u *Uses, opts *CloneOpts) (*Metadata, error) {
	_ = forge.LoggerFrom(ctx)

	if u.IsLocal() {
		return nil, fmt.Errorf("cloning local action: %s", u.Path)
	}

	if opts == nil {
		opts = &CloneOpts{}
	}

	cloneURL := opts.GitHubURL
	if cloneURL == nil {
		cloneURL = DefaultURL
	}
	cloneURL.Path = u.Repository()

	if opts.Path == "" {
		opts.Path = "."
	}

	clopts := &git.CloneOptions{
		URL:               cloneURL.String(),
		ReferenceName:     plumbing.NewTagReferenceName(u.Version),
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		InsecureSkipTLS:   opts.Insecure,
		Tags:              git.TagFollowing,
	}
	repo, err := git.PlainCloneContext(ctx, opts.Path, false, clopts)
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		repo, err = git.PlainOpen(opts.Path)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		clopts.ReferenceName = plumbing.NewBranchReferenceName(u.Version)
		repo, err = git.PlainCloneContext(ctx, opts.Path, false, clopts)
		if err != nil {
			return nil, err
		}
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	for _, filename := range ActionYAMLFilenames {
		if f, err := commit.File(filepath.Join(strings.TrimPrefix(u.Path, u.Repository()), filename)); err == nil {
			if m, err := f.Reader(); err == nil {
				return NewMetadataFromReader(m)
			}
		}
	}

	return nil, err
}
