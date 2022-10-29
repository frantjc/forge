package actions

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/frantjc/forge"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var (
	ActionYAMLFilenames = []string{"action.yml", "action.yaml"}
)

func CheckoutUses(ctx context.Context, u *Uses, opts *CheckoutOpts) (*Metadata, error) {
	_ = forge.LoggerFrom(ctx)

	if u.IsLocal() {
		return nil, fmt.Errorf("clone local action: %s", u.Path)
	}

	if opts == nil {
		opts = &CheckoutOpts{}
	}

	cloneURL := opts.GitHubURL
	if cloneURL == nil {
		cloneURL = DefaultURL
	}
	cloneURL.Path = u.GetRepository()

	if opts.Path == "" {
		opts.Path = "."
	}

	clopts := &git.CloneOptions{
		URL:               cloneURL.String(),
		ReferenceName:     plumbing.NewTagReferenceName(u.GetVersion()),
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
		clopts.ReferenceName = plumbing.NewBranchReferenceName(u.GetVersion())
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
		if f, err := commit.File(filepath.Join(u.GetActionPath(), filename)); err == nil {
			if m, err := f.Reader(); err == nil {
				return NewMetadataFromReader(m)
			}
		}
	}

	return nil, err
}
