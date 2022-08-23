package actions

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func Clone(ctx context.Context, u *Uses, opts *CloneOpts) (*Metadata, error) {
	if opts == nil {
		opts = &CloneOpts{}
	}

	cloneURL := opts.GitHubURL
	if cloneURL == nil {
		cloneURL = DefaultURL
	}
	cloneURL.Path = u.FullRepository()

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

	var f *object.File
	f, err = commit.File(filepath.Join(u.Path, "action.yml"))
	if errors.Is(err, object.ErrFileNotFound) {
		f, err = commit.File(filepath.Join(u.Path, "action.yaml"))
		if err != nil {
			return nil, ErrNotAnAction
		}
	} else if err != nil {
		return nil, err
	}

	m, err := f.Reader()
	if err != nil {
		return nil, err
	}

	return NewMetadataFromReader(m)
}
